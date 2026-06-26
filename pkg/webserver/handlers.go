package webserver

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type ScanRequest struct {
	Target string `json:"target"`
	Follow *bool  `json:"follow,omitempty"`
}

type ScanResponse struct {
	JobID     string    `json:"job_id,omitempty"`
	Target    string    `json:"target,omitempty"`
	Status    string    `json:"status"`
	Results   string    `json:"results,omitempty"`
	Stderr    string    `json:"stderr,omitempty"`
	Error     string    `json:"error,omitempty"`
	Timestamp time.Time `json:"timestamp"`
	Duration  string    `json:"duration,omitempty"`
}

func (s *Server) scanHandler(w http.ResponseWriter, r *http.Request) {
	remoteIP := getRemoteIP(r)
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")

	r.Body = http.MaxBytesReader(w, r.Body, s.maxRequestBodyBytes)
	defer func() { _ = r.Body.Close() }()

	var req ScanRequest
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{
			"error": fmt.Sprintf("invalid request: %v", err),
		})
		return
	}

	req.Target = strings.TrimSpace(req.Target)
	if req.Target == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "target cannot be empty"})
		return
	}
	if len(req.Target) > s.maxTargetLength {
		writeJSON(w, http.StatusBadRequest, map[string]string{
			"error": fmt.Sprintf("target too long (max %d chars)", s.maxTargetLength),
		})
		return
	}

	if s.isBlocked(r.Context(), req.Target) {
		writeJSON(w, http.StatusForbidden, map[string]string{"error": "target is not allowed"})
		return
	}

	queueLen, err := s.redis.LLen(r.Context(), s.queueKey).Result()
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "queue unavailable"})
		return
	}
	if queueLen >= int64(s.maxQueueSize) {
		writeJSON(w, http.StatusTooManyRequests, map[string]string{"error": "queue full, please retry shortly"})
		return
	}

	jobID, err := s.nextJobID(r.Context())
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "could not create job"})
		return
	}

	follow := s.follow
	if req.Follow != nil {
		follow = *req.Follow
	}

	now := time.Now().UTC().Format(time.RFC3339Nano)
	jobKey := s.jobKey(jobID)
	values := map[string]any{
		"id":         jobID,
		"target":     req.Target,
		"follow":     strconv.FormatBool(follow),
		"status":     statusQueued,
		"created_at": now,
	}
	if err := s.redis.HSet(r.Context(), jobKey, values).Err(); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "could not persist job"})
		return
	}
	if err := s.redis.Expire(r.Context(), jobKey, s.jobTTL).Err(); err != nil {
		slog.Warn("could not set job ttl", "job_id", jobID, "error", err)
	}
	if err := s.redis.LPush(r.Context(), s.queueKey, jobID).Err(); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "could not enqueue job"})
		return
	}

	slog.Info("Queued scan job", "job_id", jobID, "target", req.Target, "remote_ip", remoteIP)
	writeJSON(w, http.StatusAccepted, ScanResponse{
		JobID:     jobID,
		Target:    req.Target,
		Status:    statusQueued,
		Timestamp: time.Now().UTC(),
	})
}

func (s *Server) scanStatusHandler(w http.ResponseWriter, r *http.Request) {
	jobID := strings.TrimSpace(strings.TrimPrefix(r.URL.Path, "/api/scan/"))
	if jobID == "" || strings.Contains(jobID, "/") {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid job id"})
		return
	}

	job, err := s.loadJob(r.Context(), jobID)
	if err != nil {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "job not found"})
		return
	}

	writeJSON(w, http.StatusOK, job)
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		slog.Warn("failed to encode json response", "error", err)
	}
}

// isBlocked decides whether a scan target must be refused. It resolves
// hostnames so that a name pointing at an internal address can't slip past the
// IP checks (the core SSRF protection). Private/loopback/link-local/reserved
// addresses are blocked unless allowPrivateTargets is set (for deliberate LAN
// scanning). ctx bounds the DNS lookup.
//
// ponytail: residual DNS-rebinding TOCTOU remains — the target is re-resolved
// by the scan engine after this check. Closing that needs the engine to scan
// the exact IP resolved here; tracked in TODO.md.
func (s *Server) isBlocked(ctx context.Context, target string) bool {
	host := strings.ToLower(strings.TrimSpace(extractHost(target)))

	if ip := net.ParseIP(host); ip != nil {
		return s.ipBlocked(ip)
	}

	for _, entry := range s.blockedDomains {
		if host == entry || strings.HasSuffix(host, "."+entry) {
			return true
		}
	}

	addrs, err := net.DefaultResolver.LookupIPAddr(ctx, host)
	if err != nil {
		// Let unresolvable targets through; the scan itself will fail cleanly.
		return false
	}
	for _, addr := range addrs {
		if s.ipBlocked(addr.IP) {
			return true
		}
	}
	return false
}

func (s *Server) ipBlocked(ip net.IP) bool {
	if !s.allowPrivateTargets && isPrivateOrReserved(ip) {
		return true
	}
	for _, cidr := range s.blockedCIDRs {
		if cidr.Contains(ip) {
			return true
		}
	}
	return false
}

// cgnatRange is RFC6598 (100.64.0.0/10), not covered by net.IP.IsPrivate.
var cgnatRange = &net.IPNet{IP: net.IPv4(100, 64, 0, 0), Mask: net.CIDRMask(10, 32)}

func isPrivateOrReserved(ip net.IP) bool {
	return ip.IsLoopback() ||
		ip.IsPrivate() || // RFC1918 + RFC4193 ULA (fc00::/7)
		ip.IsLinkLocalUnicast() || // 169.254.0.0/16 + fe80::/10 (covers 169.254.169.254)
		ip.IsLinkLocalMulticast() ||
		ip.IsUnspecified() ||
		ip.IsMulticast() ||
		cgnatRange.Contains(ip)
}

func extractHost(target string) string {
	if strings.Contains(target, "://") {
		if u, err := url.Parse(target); err == nil {
			return u.Hostname()
		}
	}
	host, _, err := net.SplitHostPort(target)
	if err != nil {
		return target
	}
	return host
}

func (s *Server) indexHandler(w http.ResponseWriter, r *http.Request) {
	q := strings.TrimSpace(r.URL.Query().Get("q"))
	if q != "" {
		dest := "/scan?q=" + url.QueryEscape(q)
		if r.URL.Query().Get("follow") == "1" {
			dest += "&follow=1"
		}
		http.Redirect(w, r, dest, http.StatusFound)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := s.templates.ExecuteTemplate(w, "index.html", map[string]any{
		"title":   "webscan",
		"version": s.version,
	}); err != nil {
		slog.Error("failed to render index template", "error", err)
		http.Error(w, "template rendering failed", http.StatusInternalServerError)
	}
}

func (s *Server) scanPageHandler(w http.ResponseWriter, r *http.Request) {
	q := strings.TrimSpace(r.URL.Query().Get("q"))
	if q == "" {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	follow := r.URL.Query().Get("follow") == "1"
	if r.URL.Query().Get("md") == "1" {
		s.markdownScanHandler(w, r, q, follow)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := s.templates.ExecuteTemplate(w, "scan.html", map[string]any{
		"title":   "webscan — " + q,
		"version": s.version,
		"query":   q,
		"follow":  follow,
	}); err != nil {
		slog.Error("failed to render scan template", "error", err)
		http.Error(w, "template rendering failed", http.StatusInternalServerError)
	}
}

func (s *Server) markdownScanHandler(w http.ResponseWriter, r *http.Request, target string, follow bool) {
	if s.isBlocked(r.Context(), target) {
		http.Error(w, "target is not allowed", http.StatusForbidden)
		return
	}
	result, err := s.runInlineScan(r.Context(), target, follow)
	if err != nil {
		status := http.StatusInternalServerError
		if strings.Contains(err.Error(), "timed out") {
			status = http.StatusGatewayTimeout
		}
		http.Error(w, err.Error(), status)
		return
	}
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	if _, err := fmt.Fprint(w, stripANSI(result)); err != nil {
		slog.Debug("webserver: Error writing response", "error", err)
	}
}

var ansiEscapeRe = regexp.MustCompile(`\x1b\[[0-9;?]*[ -/]*[@-~]`)

func stripANSI(s string) string {
	return ansiEscapeRe.ReplaceAllString(s, "")
}

// getRemoteIP returns the client IP. X-Forwarded-For / X-Real-IP are only
// trusted when the direct peer is a private/loopback address (i.e. a reverse
// proxy on the same network); otherwise they are attacker-spoofable, so the
// real peer address is used.
func getRemoteIP(r *http.Request) string {
	peer := r.RemoteAddr
	if host, _, err := net.SplitHostPort(peer); err == nil {
		peer = host
	}

	peerIP := net.ParseIP(peer)
	if peerIP != nil && isPrivateOrReserved(peerIP) {
		if forwarded := r.Header.Get("X-Forwarded-For"); forwarded != "" {
			parts := strings.Split(forwarded, ",")
			return strings.TrimSpace(parts[0])
		}
		if realIP := strings.TrimSpace(r.Header.Get("X-Real-IP")); realIP != "" {
			return realIP
		}
	}
	return peer
}
