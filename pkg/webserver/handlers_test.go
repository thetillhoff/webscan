package webserver

import (
	"html/template"
	"net/http"
	"net/http/httptest"
	"testing"
)

func testServer(t *testing.T) *Server {
	t.Helper()
	tmpl, err := template.ParseFS(htmlTemplates, "templates/*")
	if err != nil {
		t.Fatalf("failed to parse templates: %v", err)
	}
	return &Server{templates: tmpl, version: "test"}
}

func TestStripANSI(t *testing.T) {
	tests := []struct {
		in   string
		want string
	}{
		{"plain text", "plain text"},
		{"\x1b[32mgreen\x1b[0m", "green"},
		{"\x1b[1;31mbold red\x1b[0m text", "bold red text"},
		{"no escapes", "no escapes"},
	}
	for _, tt := range tests {
		got := stripANSI(tt.in)
		if got != tt.want {
			t.Errorf("stripANSI(%q) = %q, want %q", tt.in, got, tt.want)
		}
	}
}

func TestIndexHandler_NoQ_RendersLandingPage(t *testing.T) {
	s := testServer(t)
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	s.indexHandler(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
}

func TestIndexHandler_WithQ_Redirects(t *testing.T) {
	s := testServer(t)
	req := httptest.NewRequest(http.MethodGet, "/?q=example.com", nil)
	rec := httptest.NewRecorder()
	s.indexHandler(rec, req)
	if rec.Code != http.StatusFound {
		t.Fatalf("expected 302, got %d", rec.Code)
	}
	if got := rec.Header().Get("Location"); got != "/scan?q=example.com" {
		t.Fatalf("expected /scan?q=example.com, got %s", got)
	}
}

func TestIndexHandler_WithQAndFollow_Redirects(t *testing.T) {
	s := testServer(t)
	req := httptest.NewRequest(http.MethodGet, "/?q=example.com&follow=1", nil)
	rec := httptest.NewRecorder()
	s.indexHandler(rec, req)
	if rec.Code != http.StatusFound {
		t.Fatalf("expected 302, got %d", rec.Code)
	}
	if got := rec.Header().Get("Location"); got != "/scan?q=example.com&follow=1" {
		t.Fatalf("expected /scan?q=example.com&follow=1, got %s", got)
	}
}

func TestScanPageHandler_NoQ_RedirectsToRoot(t *testing.T) {
	s := testServer(t)
	req := httptest.NewRequest(http.MethodGet, "/scan", nil)
	rec := httptest.NewRecorder()
	s.scanPageHandler(rec, req)
	if rec.Code != http.StatusFound {
		t.Fatalf("expected 302, got %d", rec.Code)
	}
	if got := rec.Header().Get("Location"); got != "/" {
		t.Fatalf("expected /, got %s", got)
	}
}

func TestScanPageHandler_WithQ_RendersPage(t *testing.T) {
	s := testServer(t)
	req := httptest.NewRequest(http.MethodGet, "/scan?q=example.com", nil)
	rec := httptest.NewRecorder()
	s.scanPageHandler(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d; body: %s", rec.Code, rec.Body.String())
	}
}
