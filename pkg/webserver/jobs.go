package webserver

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log/slog"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/thetillhoff/webscan/v3/pkg/webscan"
)

const (
	statusQueued    = "queued"
	statusRunning   = "running"
	statusCompleted = "completed"
	statusFailed    = "failed"
	statusTimeout   = "timeout"

	jobStatusFlushInterval = 500 * time.Millisecond
)

type synchronizedBuffer struct {
	mu  sync.Mutex
	buf bytes.Buffer
}

func (b *synchronizedBuffer) Write(p []byte) (int, error) {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.buf.Write(p)
}

func (b *synchronizedBuffer) String() string {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.buf.String()
}

func (s *Server) jobKey(jobID string) string {
	return s.jobPrefix + jobID
}

func (s *Server) nextJobID(ctx context.Context) (string, error) {
	id, err := s.redis.Incr(ctx, s.jobIDKey).Result()
	if err != nil {
		return "", err
	}
	return strconv.FormatInt(id, 10), nil
}

func (s *Server) loadJob(ctx context.Context, jobID string) (ScanResponse, error) {
	data, err := s.redis.HGetAll(ctx, s.jobKey(jobID)).Result()
	if err != nil {
		return ScanResponse{}, err
	}
	if len(data) == 0 {
		return ScanResponse{}, redis.Nil
	}

	resp := ScanResponse{
		JobID:   jobID,
		Target:  data["target"],
		Status:  data["status"],
		Results: data["result"],
		Stderr:  data["status_output"],
		Error:   data["error"],
	}
	if v := data["duration"]; v != "" {
		resp.Duration = v
	}
	if v := data["updated_at"]; v != "" {
		if t, err := time.Parse(time.RFC3339Nano, v); err == nil {
			resp.Timestamp = t
		}
	}
	if resp.Timestamp.IsZero() {
		resp.Timestamp = time.Now().UTC()
	}

	return resp, nil
}

func (s *Server) runWorker(ctx context.Context, workerID int) {
	slog.Info("scan worker started", "worker_id", workerID)
	for {
		select {
		case <-ctx.Done():
			slog.Info("scan worker stopped", "worker_id", workerID)
			return
		default:
		}

		values, err := s.redis.BRPop(ctx, 2*time.Second, s.queueKey).Result()
		if err == redis.Nil {
			continue
		}
		if err != nil {
			if ctx.Err() != nil {
				return
			}
			slog.Warn("worker queue pop failed", "worker_id", workerID, "error", err)
			time.Sleep(500 * time.Millisecond)
			continue
		}
		if len(values) < 2 {
			continue
		}

		jobID := values[1]
		s.processJob(ctx, workerID, jobID)
	}
}

func (s *Server) processJob(ctx context.Context, workerID int, jobID string) {
	jobKey := s.jobKey(jobID)
	data, err := s.redis.HGetAll(ctx, jobKey).Result()
	if err != nil || len(data) == 0 {
		slog.Warn("worker could not load job", "worker_id", workerID, "job_id", jobID, "error", err)
		return
	}

	target := strings.TrimSpace(data["target"])
	follow, _ := strconv.ParseBool(data["follow"])

	now := time.Now().UTC().Format(time.RFC3339Nano)
	if err := s.redis.HSet(ctx, jobKey, map[string]any{
		"status":     statusRunning,
		"started_at": now,
		"updated_at": now,
	}).Err(); err != nil {
		slog.Warn("worker could not mark job running", "worker_id", workerID, "job_id", jobID, "error", err)
	}

	outputBuffer := &synchronizedBuffer{}
	statusBuffer := &synchronizedBuffer{}
	stopStatusStreaming := s.startStatusStreaming(jobKey, statusBuffer)
	defer stopStatusStreaming()

	engine, err := s.newEngine(outputBuffer, statusBuffer, follow)
	if err != nil {
		s.finishJob(ctx, jobKey, statusFailed, "", "", fmt.Sprintf("failed to initialize scan engine: %v", err), "")
		return
	}

	started := time.Now()
	done := make(chan struct{})
	var scanErr error

	go func() {
		defer close(done)
		scanErr = engine.Scan(target)
	}()

	select {
	case <-done:
		duration := time.Since(started).String()
		if scanErr != nil {
			s.finishJob(ctx, jobKey, statusFailed, "", statusBuffer.String(), fmt.Sprintf("scan failed: %v", scanErr), duration)
			return
		}
		s.finishJob(ctx, jobKey, statusCompleted, outputBuffer.String(), statusBuffer.String(), "", duration)
	case <-time.After(s.scanTimeout):
		duration := time.Since(started).String()
		s.finishJob(ctx, jobKey, statusTimeout, outputBuffer.String(), statusBuffer.String(), fmt.Sprintf("scan timed out after %s", s.scanTimeout), duration)
	}
}

func (s *Server) startStatusStreaming(jobKey string, statusBuffer *synchronizedBuffer) func() {
	done := make(chan struct{})
	stopped := make(chan struct{})

	go func() {
		defer close(stopped)
		ticker := time.NewTicker(jobStatusFlushInterval)
		defer ticker.Stop()

		lastSnapshot := ""
		for {
			select {
			case <-ticker.C:
				snapshot := statusBuffer.String()
				if snapshot == lastSnapshot {
					continue
				}
				s.persistRunningStatus(jobKey, snapshot)
				lastSnapshot = snapshot
			case <-done:
				snapshot := statusBuffer.String()
				if snapshot != lastSnapshot {
					s.persistRunningStatus(jobKey, snapshot)
				}
				return
			}
		}
	}()

	return func() {
		close(done)
		<-stopped
	}
}

func (s *Server) persistRunningStatus(jobKey, statusOutput string) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	now := time.Now().UTC().Format(time.RFC3339Nano)
	if err := s.redis.HSet(ctx, jobKey, map[string]any{
		"status_output": statusOutput,
		"updated_at":    now,
	}).Err(); err != nil {
		slog.Warn("worker could not persist running status", "job_key", jobKey, "error", err)
	}
}

func (s *Server) finishJob(ctx context.Context, jobKey, status, result, statusOutput, errMsg, duration string) {
	now := time.Now().UTC().Format(time.RFC3339Nano)
	fields := map[string]any{
		"status":        status,
		"result":        result,
		"status_output": statusOutput,
		"error":         errMsg,
		"duration":      duration,
		"finished_at":   now,
		"updated_at":    now,
	}
	if err := s.redis.HSet(ctx, jobKey, fields).Err(); err != nil {
		slog.Warn("worker could not finish job", "job_key", jobKey, "error", err)
	}
	if err := s.redis.Expire(ctx, jobKey, s.jobTTL).Err(); err != nil {
		slog.Warn("worker could not set job ttl", "job_key", jobKey, "error", err)
	}
}

func (s *Server) newEngine(stdout io.Writer, statusOut io.Writer, followRedirects bool) (*webscan.Engine, error) {
	engine, err := webscan.NewEngine(
		stdout,
		statusOut,
		s.disableColor,
		s.dnsServer,
		followRedirects,
		s.requestTimeout,
		false,
		false,
		false,
		false,
		false,
		false,
		false,
		false,
		false,
		s.writeMutex,
	)
	if err != nil {
		return nil, err
	}

	s.scanOptions.Apply(&engine)
	return &engine, nil
}
