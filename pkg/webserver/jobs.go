package webserver

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/thetillhoff/webscan/v5/pkg/webscan"
)

const (
	statusQueued    = "queued"
	statusRunning   = "running"
	statusCompleted = "completed"
	statusFailed    = "failed"
	statusTimeout   = "timeout"

	jobStatusFlushInterval = 500 * time.Millisecond

	// statusBufferMaxSize keeps only the most recent status text. Scan
	// results (a separate, unbounded buffer) are never trimmed.
	statusBufferMaxSize = 16 * 1024

	// sseWriteTimeout bounds each individual SSE write. Without it, a client
	// that stops reading (Slowloris) could block a write, and thus this
	// handler's goroutine, forever.
	sseWriteTimeout = 10 * time.Second

	// maxScanDuration is an absolute backstop independent of staleTimeout: a
	// target that keeps status trickling in just under staleTimeout would
	// otherwise occupy a worker slot forever. Generous enough that no
	// legitimate scan (a full 65535-port scan included) should ever hit it.
	maxScanDuration = 30 * time.Minute
)

type synchronizedBuffer struct {
	mu  sync.Mutex
	buf bytes.Buffer
	// maxSize caps retained bytes, keeping only the most recent ones once
	// exceeded. 0 means unbounded. A full port scan writes one status line per
	// port (up to 65535 lines); without a cap, every poll/SSE push resends the
	// whole growing history, ballooning to multi-MB payloads.
	maxSize int
	// lastWrite is when Write was last called, used to detect a stalled scan
	// (no progress for staleTimeout). Zero until the first write.
	lastWrite time.Time
}

func (b *synchronizedBuffer) Write(p []byte) (int, error) {
	b.mu.Lock()
	defer b.mu.Unlock()
	n, err := b.buf.Write(p)
	if b.maxSize > 0 && b.buf.Len() > b.maxSize {
		kept := append([]byte(nil), b.buf.Bytes()[b.buf.Len()-b.maxSize:]...)
		b.buf.Reset()
		b.buf.Write(kept)
	}
	b.lastWrite = time.Now()
	return n, err
}

func (b *synchronizedBuffer) String() string {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.buf.String()
}

func (b *synchronizedBuffer) LastWrite() time.Time {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.lastWrite
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
	fullPortScan, _ := strconv.ParseBool(data["full_port_scan"])

	now := time.Now().UTC().Format(time.RFC3339Nano)
	if err := s.redis.HSet(ctx, jobKey, map[string]any{
		"status":     statusRunning,
		"started_at": now,
		"updated_at": now,
	}).Err(); err != nil {
		slog.Warn("worker could not mark job running", "worker_id", workerID, "job_id", jobID, "error", err)
	}

	outputBuffer := &synchronizedBuffer{}
	statusBuffer := &synchronizedBuffer{maxSize: statusBufferMaxSize, lastWrite: time.Now()}
	stopStatusStreaming := s.startStatusStreaming(jobKey, statusBuffer)
	defer stopStatusStreaming()

	engine, err := s.newEngine(outputBuffer, statusBuffer, follow, fullPortScan)
	if err != nil {
		s.finishJob(ctx, jobKey, statusFailed, "", "", fmt.Sprintf("failed to initialize scan engine: %v", err), "")
		return
	}

	started := time.Now()
	// Cancel propagates to engine.Scan either when the watchdog below detects
	// no status progress for staleTimeout, or when maxScanDuration elapses
	// regardless of progress — a backstop so a target that keeps trickling
	// status just under staleTimeout can't occupy a worker slot forever.
	scanCtx, cancel := context.WithTimeout(ctx, maxScanDuration)
	defer cancel()

	done := make(chan struct{})
	var scanErr error

	go func() {
		defer close(done)
		scanErr = engine.Scan(scanCtx, target)
	}()

	staleFired := make(chan struct{})
	go s.watchForStale(scanCtx, cancel, statusBuffer, done, staleFired)

	select {
	case <-done:
		duration := time.Since(started).String()
		if scanErr != nil {
			s.finishJob(ctx, jobKey, statusFailed, "", statusBuffer.String(), fmt.Sprintf("scan failed: %v", scanErr), duration)
			return
		}
		s.finishJob(ctx, jobKey, statusCompleted, outputBuffer.String(), statusBuffer.String(), "", duration)
	case <-scanCtx.Done():
		duration := time.Since(started).String()
		errMsg := fmt.Sprintf("scan exceeded the maximum duration of %s", maxScanDuration)
		select {
		case <-staleFired:
			errMsg = fmt.Sprintf("scan produced no status update for %s, aborting", s.staleTimeout)
		default:
		}
		s.finishJob(ctx, jobKey, statusTimeout, outputBuffer.String(), statusBuffer.String(), errMsg, duration)
	}
}

// watchForStale cancels scanCtx once statusBuffer has gone staleTimeout
// without a write, i.e. the scan has stopped making progress, closing
// staleFired first so the caller can tell that apart from scanCtx ending for
// another reason (maxScanDuration elapsed, parent ctx cancelled). Returns
// once done closes (the scan finished on its own) or scanCtx is done.
func (s *Server) watchForStale(scanCtx context.Context, cancel context.CancelFunc, statusBuffer *synchronizedBuffer, done <-chan struct{}, staleFired chan<- struct{}) {
	ticker := time.NewTicker(jobStatusFlushInterval)
	defer ticker.Stop()

	for {
		select {
		case <-done:
			return
		case <-scanCtx.Done():
			return
		case <-ticker.C:
			if time.Since(statusBuffer.LastWrite()) >= s.staleTimeout {
				close(staleFired)
				cancel()
				return
			}
		}
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

// scanEventsHandler streams job status over Server-Sent Events instead of
// making the client poll — replaces one HTTP round trip per poll tick with a
// single connection, pushed on the same cadence the worker persists updates.
func (s *Server) scanEventsHandler(w http.ResponseWriter, r *http.Request) {
	jobID := r.PathValue("id")
	if jobID == "" {
		http.Error(w, "invalid job id", http.StatusBadRequest)
		return
	}

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "streaming unsupported", http.StatusInternalServerError)
		return
	}

	// A scan (and thus this stream) may run longer than the server's default
	// write timeout, so it's overridden here — but with a bounded per-write
	// deadline, not cleared outright: a client that stops reading (Slowloris)
	// would otherwise be able to hold this connection open forever.
	rc := http.NewResponseController(w)

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Set("Connection", "keep-alive")
	w.WriteHeader(http.StatusOK)

	ticker := time.NewTicker(jobStatusFlushInterval)
	defer ticker.Stop()

	var lastPayload string
	for {
		job, err := s.loadJob(r.Context(), jobID)
		if err != nil {
			job = ScanResponse{JobID: jobID, Status: statusFailed, Error: "job not found"}
		}

		if payload, err := json.Marshal(job); err == nil && string(payload) != lastPayload {
			if err := rc.SetWriteDeadline(time.Now().Add(sseWriteTimeout)); err != nil {
				slog.Debug("webserver: could not set SSE write deadline", "error", err)
			}
			if _, err := fmt.Fprintf(w, "data: %s\n\n", payload); err != nil {
				return
			}
			flusher.Flush()
			lastPayload = string(payload)
		}

		switch strings.ToLower(job.Status) {
		case statusCompleted, statusFailed, statusTimeout:
			return
		}

		select {
		case <-r.Context().Done():
			return
		case <-ticker.C:
		}
	}
}

func (s *Server) newEngine(stdout io.Writer, statusOut io.Writer, followRedirects bool, fullPortScan bool) (*webscan.Engine, error) {
	engine, err := webscan.NewEngine(
		stdout,
		statusOut,
		s.disableColor,
		s.dnsServer,
		followRedirects,
		fullPortScan,
		s.requestTimeout,
		s.scanOptions,
		s.writeMutex,
	)
	if err != nil {
		return nil, err
	}

	return &engine, nil
}
