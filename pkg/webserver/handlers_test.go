package webserver

import (
	"html/template"
	"net/http"
	"net/http/httptest"
	"strings"
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

func TestIndexHandler_WithQAndFullPort_Redirects(t *testing.T) {
	s := testServer(t)
	req := httptest.NewRequest(http.MethodGet, "/?q=example.com&fullport=1", nil)
	rec := httptest.NewRecorder()
	s.indexHandler(rec, req)
	if rec.Code != http.StatusFound {
		t.Fatalf("expected 302, got %d", rec.Code)
	}
	if got := rec.Header().Get("Location"); got != "/scan?q=example.com&fullport=1" {
		t.Fatalf("expected /scan?q=example.com&fullport=1, got %s", got)
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

func TestScanPageHandler_WithFullPort_ChecksBox(t *testing.T) {
	s := testServer(t)
	req := httptest.NewRequest(http.MethodGet, "/scan?q=example.com&fullport=1", nil)
	rec := httptest.NewRecorder()
	s.scanPageHandler(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d; body: %s", rec.Code, rec.Body.String())
	}
	if !strings.Contains(rec.Body.String(), `id="fullPortScan" checked>`) {
		t.Fatalf("expected fullPortScan checkbox to be checked, body: %s", rec.Body.String())
	}
}
