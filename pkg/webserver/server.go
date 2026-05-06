package webserver

import (
	"context"
	"embed"
	"html/template"
	"io/fs"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/thetillhoff/webscan/v3/pkg/webscan"
)

const (
	defaultScanTimeout        = 30 * time.Second
	defaultMaxRequestBytes    = int64(4096)
	defaultMaxTargetLength    = 2048
	defaultMaxConcurrentScans = 1
	defaultMaxQueueSize       = 1000
	defaultJobTTL             = 24 * time.Hour
)

var (
	//go:embed static/*
	staticAssets embed.FS

	//go:embed templates/*
	htmlTemplates embed.FS
)

type Server struct {
	writeMutex          *sync.Mutex
	dnsServer           string
	follow              bool
	requestTimeout      time.Duration
	scanOptions         webscan.ScanOptions
	disableColor        bool
	port                string
	router              http.Handler
	templates           *template.Template
	scanTimeout         time.Duration
	maxRequestBodyBytes int64
	maxTargetLength     int
	workerCount         int
	maxQueueSize        int
	jobTTL              time.Duration
	redis               *redis.Client
	queueKey            string
	jobPrefix           string
	jobIDKey            string
}

func NewServer(
	noColor bool,
	dnsServer string,
	followRedirects bool,
	requestTimeout time.Duration,
	port string,
	writeMutex *sync.Mutex,
	scanTimeout time.Duration,
	maxConcurrentScans int,
	maxRequestBodyBytes int64,
	redisAddr string,
	redisPassword string,
	redisDB int,
	maxQueueSize int,
) (*Server, error) {
	if writeMutex == nil {
		writeMutex = &sync.Mutex{}
	}
	if scanTimeout <= 0 {
		scanTimeout = defaultScanTimeout
	}
	if maxConcurrentScans <= 0 {
		maxConcurrentScans = defaultMaxConcurrentScans
	}
	if maxRequestBodyBytes <= 0 {
		maxRequestBodyBytes = defaultMaxRequestBytes
	}

	if maxQueueSize <= 0 {
		maxQueueSize = defaultMaxQueueSize
	}
	if redisAddr == "" {
		redisAddr = "127.0.0.1:6379"
	}

	tmpl := template.Must(template.ParseFS(htmlTemplates, "templates/*"))
	redisClient := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: redisPassword,
		DB:       redisDB,
	})
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := redisClient.Ping(ctx).Err(); err != nil {
		return nil, err
	}

	if requestTimeout <= 0 {
		requestTimeout = 5 * time.Second
	}

	server := &Server{
		writeMutex:          writeMutex,
		dnsServer:           dnsServer,
		follow:              followRedirects,
		requestTimeout:      requestTimeout,
		scanOptions:         webscan.WebScanOptions(),
		port:                port,
		disableColor:        noColor,
		templates:           tmpl,
		scanTimeout:         scanTimeout,
		maxRequestBodyBytes: maxRequestBodyBytes,
		maxTargetLength:     defaultMaxTargetLength,
		workerCount:         maxConcurrentScans,
		maxQueueSize:        maxQueueSize,
		jobTTL:              defaultJobTTL,
		redis:               redisClient,
		queueKey:            "webscan:jobs:queue",
		jobPrefix:           "webscan:job:",
		jobIDKey:            "webscan:jobs:next_id",
	}

	server.setupRouter()
	return server, nil
}

func (s *Server) setupRouter() {
	mux := http.NewServeMux()

	staticFiles, err := fs.Sub(staticAssets, "static")
	if err != nil {
		panic(err)
	}
	mux.Handle("GET /static/", http.StripPrefix("/static/", http.FileServer(http.FS(staticFiles))))

	mux.HandleFunc("GET /", s.indexHandler)
	mux.HandleFunc("GET /api/health", s.healthHandler)
	mux.HandleFunc("POST /api/scan", s.scanHandler)
	mux.HandleFunc("GET /api/scan/", s.scanStatusHandler)

	s.router = s.withRequestLogging(mux)
}

func (s *Server) withRequestLogging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
	})
}

func (s *Server) indexHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := s.templates.ExecuteTemplate(w, "index.html", map[string]any{
		"title": "webscan - Web Security Scanner",
	}); err != nil {
		slog.Error("failed to render index template", "error", err)
		http.Error(w, "template rendering failed", http.StatusInternalServerError)
	}
}

func (s *Server) healthHandler(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "healthy"})
}

func (s *Server) Run() error {
	slog.Info("Starting webscan web server", "port", s.port)
	ctx, cancelWorkers := context.WithCancel(context.Background())
	for i := 0; i < s.workerCount; i++ {
		go s.runWorker(ctx, i+1)
	}
	defer cancelWorkers()

	server := &http.Server{
		Addr:         ":" + s.port,
		Handler:      s.router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: s.scanTimeout + 30*time.Second,
		IdleTimeout:  120 * time.Second,
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	serverErrors := make(chan error, 1)
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			serverErrors <- err
		}
		close(serverErrors)
	}()

	select {
	case err := <-serverErrors:
		return err
	case <-stop:
		slog.Info("Shutting down web server...")
		cancelWorkers()
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		return server.Shutdown(ctx)
	}
}
