# Two-Page Layout with URL Query Parameter — Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Split the web UI into a centered landing page (`/`) and a scan/results page (`/scan?q=<target>`), with `?follow=1` and `?md=1` URL parameters wired end-to-end.

**Architecture:** The Go server gains a `GET /scan` route. `indexHandler` redirects to `/scan?q=...` when `q` is present; `scanPageHandler` redirects to `/` when `q` is absent, runs an inline scan for `?md=1`, or renders `scan.html`. The scan page auto-fires via JS reading URL params on load; the landing page is a plain HTML GET form pointing at `/scan`.

**Tech Stack:** Go 1.26 (`net/http`, `html/template`, `regexp`), vanilla JS, vanilla CSS.

---

## File Map

| File | Action | What changes |
| --- | --- | --- |
| `pkg/webserver/handlers_test.go` | Create | Unit tests: routing redirects, `stripANSI` |
| `pkg/webserver/handlers.go` | Modify | Update `indexHandler`; add `scanPageHandler`, `markdownScanHandler`, `stripANSI` |
| `pkg/webserver/jobs.go` | Modify | Add `runInlineScan` |
| `pkg/webserver/server.go` | Modify | Register `GET /scan` route |
| `pkg/webserver/templates/index.html` | Modify | Landing page: centered form, `action="/scan" method="GET"` |
| `pkg/webserver/templates/scan.html` | Create | Scan/results page: pre-filled form, spinner, logs, results |
| `pkg/webserver/static/style.css` | Modify | Add `.landing-main` / `.landing-container` / `.landing-title` |
| `pkg/webserver/static/script.js` | Modify | Path-based init: landing redirect check; scan auto-start |
| `tests/e2e/webscan.spec.ts` | Modify | Update for new two-page flow |

---

## Task 1: Write failing handler tests

**Files:**

- Create: `pkg/webserver/handlers_test.go`

- [ ] **Step 1: Create the test file**

```go
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
```

- [ ] **Step 2: Run tests — expect FAIL (compilation error: `stripANSI` and `scanPageHandler` undefined)**

```bash
cd /var/home/thetillhoff/code/webscan && go test ./pkg/webserver/...
```

Expected: compilation error mentioning `stripANSI` and `scanPageHandler` undefined.

---

## Task 2: Backend — routing logic and stubs

**Files:**

- Modify: `pkg/webserver/handlers.go`
- Modify: `pkg/webserver/server.go`

- [ ] **Step 1: Add `regexp` import and `stripANSI` to `handlers.go`**

Add `"regexp"` to the import block. Add after the `getRemoteIP` function (line 187):

```go
var ansiEscapeRe = regexp.MustCompile(`\x1b\[[0-9;?]*[ -/]*[@-~]`)

func stripANSI(s string) string {
 return ansiEscapeRe.ReplaceAllString(s, "")
}
```

- [ ] **Step 2: Move and update `indexHandler`**

`indexHandler` currently lives at the bottom of `server.go` (lines 178–187). Delete those lines from `server.go`, then add the updated version to `handlers.go` instead (it needs `strings` and `url` imports that are already present there):

```go
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
```

- [ ] **Step 3: Add `scanPageHandler` and a `markdownScanHandler` stub to `handlers.go`**

Append to `handlers.go`:

```go
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

// markdownScanHandler stub — replaced in Task 6
func (s *Server) markdownScanHandler(w http.ResponseWriter, _ *http.Request, _ string, _ bool) {
 http.Error(w, "not yet implemented", http.StatusNotImplemented)
}
```

- [ ] **Step 4: Register `GET /scan` in `server.go` `setupRouter`**

In `setupRouter` (around line 164), add after the `GET /` line:

```go
mux.HandleFunc("GET /scan", s.scanPageHandler)
```

- [ ] **Step 5: Run tests — routing redirects should pass; `TestScanPageHandler_WithQ_RendersPage` still fails (scan.html not found)**

```bash
cd /var/home/thetillhoff/code/webscan && go test ./pkg/webserver/...
```

Expected: `TestStripANSI`, `TestIndexHandler_*`, `TestScanPageHandler_NoQ_RedirectsToRoot` PASS; `TestScanPageHandler_WithQ_RendersPage` FAIL with template error.

---

## Task 3: Templates — landing page and scan page

**Files:**

- Modify: `pkg/webserver/templates/index.html`
- Create: `pkg/webserver/templates/scan.html`

- [ ] **Step 1: Replace `index.html` with the landing page**

Full file content:

```html
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <meta name="description" content="Webscan - Web Security Scanner">
    <title>webscan</title>
    <link rel="stylesheet" href="/static/style.css">
</head>
<body>
    <main class="landing-main">
        <div class="landing-container">
            <h1 class="landing-title">webscan</h1>
            <form action="/scan" method="GET" class="scan-form landing-form">
                <div class="input-group">
                    <input type="text" name="q" id="targetInput"
                           placeholder="Enter domain, IP, or URL (e.g., google.com)"
                           autocomplete="off" autofocus>
                    <button type="submit" id="scanButton">Scan</button>
                </div>
                <div class="advanced-options">
                    <label>
                        <input type="checkbox" name="follow" value="1" id="followRedirects"> Follow CNAMEs and HTTP redirects
                    </label>
                </div>
            </form>
        </div>
    </main>

    <footer>
        <a href="https://github.com/thetillhoff/webscan">github.com/thetillhoff/webscan</a>
        {{if .version}} · {{.version}}{{end}}
    </footer>

    <script src="/static/script.js"></script>
</body>
</html>
```

- [ ] **Step 2: Create `scan.html`**

Full file content:

```html
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <meta name="description" content="Webscan - Web Security Scanner">
    <title>webscan{{if .query}} — {{.query}}{{end}}</title>
    <link rel="stylesheet" href="/static/style.css">
</head>
<body>
    <main>
        <div class="scan-container">
            <form action="/scan" method="GET" id="scanForm" class="scan-form">
                <div class="input-group">
                    <input type="text" name="q" id="targetInput"
                           value="{{.query}}"
                           placeholder="Enter domain, IP, or URL (e.g., google.com)"
                           autocomplete="off" required>
                    <button type="submit" id="scanButton">Scan</button>
                </div>
                <div class="advanced-options">
                    <label>
                        <input type="checkbox" name="follow" value="1" id="followRedirects"{{if .follow}} checked{{end}}> Follow CNAMEs and HTTP redirects
                    </label>
                </div>
            </form>

            <div id="spinner" class="spinner">
                <div class="spinner-icon"></div>
                <span id="spinnerText">Scanning...</span>
            </div>

            <div id="logsSection" class="logs-section" style="display:none;">
                <button type="button" id="toggleLogs" class="toggle-logs-btn">Show logs</button>
                <pre id="logsOutput" class="logs-output" style="display:none;"></pre>
            </div>

            <div id="resultsSection" class="results-section" style="display:none;">
                <pre id="scanResults" class="scan-results"></pre>
            </div>

            <div id="errorSection" class="error-section" style="display:none;">
                <p id="errorMessage"></p>
            </div>
        </div>
    </main>

    <footer>
        <a href="https://github.com/thetillhoff/webscan">github.com/thetillhoff/webscan</a>
        {{if .version}} · {{.version}}{{end}}
    </footer>

    <script src="/static/script.js"></script>
</body>
</html>
```

- [ ] **Step 3: Run all tests — all should pass**

```bash
cd /var/home/thetillhoff/code/webscan && go test ./pkg/webserver/...
```

Expected: all tests PASS.

- [ ] **Step 4: Commit**

```bash
git add pkg/webserver/handlers_test.go pkg/webserver/handlers.go pkg/webserver/server.go \
        pkg/webserver/templates/index.html pkg/webserver/templates/scan.html
git commit -m "feat: add two-page routing with scan page handler"
```

---

## Task 4: CSS — landing page centered layout

**Files:**

- Modify: `pkg/webserver/static/style.css`

- [ ] **Step 1: Add landing page styles to the end of `style.css`**

```css
/* Landing page */
.landing-main {
    flex: 1;
    display: flex;
    align-items: center;
    justify-content: center;
    width: 100%;
    max-width: 100%;
    margin: 0;
    padding: 0 1rem;
}

.landing-container {
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 2rem;
    width: 100%;
    max-width: 600px;
}

.landing-title {
    font-size: 3rem;
    font-weight: 700;
    color: #4285f4;
    letter-spacing: -1px;
}

.landing-form {
    width: 100%;
}
```

- [ ] **Step 2: Verify it compiles**

```bash
cd /var/home/thetillhoff/code/webscan && go build ./...
```

Expected: no output (success).

- [ ] **Step 3: Commit**

```bash
git add pkg/webserver/static/style.css
git commit -m "feat: add centered landing page layout"
```

---

## Task 5: JS — two-page script

**Files:**

- Modify: `pkg/webserver/static/script.js`

- [ ] **Step 1: Replace `script.js` entirely**

```javascript
document.addEventListener('DOMContentLoaded', function () {
    const path = window.location.pathname;
    if (path === '/') {
        initLandingPage();
    } else if (path === '/scan') {
        initScanPage();
    }
});

function initLandingPage() {
    const followCheckbox = document.getElementById('followRedirects');

    if (localStorage.getItem('followRedirects') === 'true') {
        followCheckbox.checked = true;
    }
    followCheckbox.addEventListener('change', function () {
        localStorage.setItem('followRedirects', String(followCheckbox.checked));
    });

    // Redirect if someone pastes a /?q=... URL
    const q = (new URLSearchParams(window.location.search).get('q') || '').trim();
    if (q) {
        const dest = new URLSearchParams({ q });
        if (followCheckbox.checked) dest.set('follow', '1');
        window.location.replace('/scan?' + dest.toString());
    }
}

function initScanPage() {
    const params = new URLSearchParams(window.location.search);
    const q = (params.get('q') || '').trim();
    const follow = params.get('follow') === '1';

    const form = document.getElementById('scanForm');
    const input = document.getElementById('targetInput');
    const button = document.getElementById('scanButton');
    const followCheckbox = document.getElementById('followRedirects');
    const spinner = document.getElementById('spinner');
    const spinnerText = document.getElementById('spinnerText');
    const logsSection = document.getElementById('logsSection');
    const toggleLogsBtn = document.getElementById('toggleLogs');
    const logsOutput = document.getElementById('logsOutput');
    const resultsSection = document.getElementById('resultsSection');
    const scanResults = document.getElementById('scanResults');
    const errorSection = document.getElementById('errorSection');
    const errorMessage = document.getElementById('errorMessage');

    let logsExpanded = false;

    toggleLogsBtn.addEventListener('click', function () {
        logsExpanded = !logsExpanded;
        logsOutput.style.display = logsExpanded ? 'block' : 'none';
        toggleLogsBtn.textContent = logsExpanded ? 'Hide logs' : 'Show logs';
        if (logsExpanded) logsOutput.scrollTop = logsOutput.scrollHeight;
    });

    // Form submit navigates to /scan?q=... — triggers a fresh page load + scan
    form.addEventListener('submit', function (e) {
        e.preventDefault();
        const newQ = input.value.trim();
        if (!newQ) return;
        const dest = new URLSearchParams({ q: newQ });
        if (followCheckbox.checked) dest.set('follow', '1');
        window.location.href = '/scan?' + dest.toString();
    });

    // Auto-start scan from URL params on page load
    if (q) {
        runScan(q, follow, {
            button, input, spinner, spinnerText,
            logsSection, logsOutput,
            resultsSection, scanResults,
            errorSection, errorMessage,
        });
    }
}

async function runScan(target, follow, els) {
    const { button, input, spinner, spinnerText, logsSection, logsOutput,
            resultsSection, scanResults, errorSection, errorMessage } = els;

    input.disabled = true;
    button.disabled = true;
    button.textContent = 'Scanning...';

    try {
        const enqueueResp = await fetch('/api/scan', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ target, follow }),
        });
        const enqueueData = await enqueueResp.json();
        if (!enqueueResp.ok) throw new Error(enqueueData.error || enqueueResp.statusText);

        await pollScanJob(enqueueData.job_id, { spinner, spinnerText, logsSection, logsOutput, resultsSection, scanResults });
    } catch (err) {
        spinner.style.display = 'none';
        errorSection.style.display = 'block';
        errorMessage.textContent = err.message;
    } finally {
        input.disabled = false;
        button.disabled = false;
        button.textContent = 'Scan';
    }
}

async function pollScanJob(jobID, { spinner, spinnerText, logsSection, logsOutput, resultsSection, scanResults }) {
    const pollIntervalMs = 1000;
    const timeoutMs = 180000;
    const startedAt = Date.now();

    while (true) {
        if (Date.now() - startedAt > timeoutMs) {
            throw new Error('Scan timed out');
        }

        const resp = await fetch('/api/scan/' + encodeURIComponent(jobID));
        const data = await resp.json();
        if (!resp.ok) throw new Error(data.error || resp.statusText);

        const status = (data.status || '').toLowerCase();

        if (data.stderr) {
            logsSection.style.display = 'block';
            logsOutput.textContent = cleanAnsi(data.stderr);
            logsOutput.scrollTop = logsOutput.scrollHeight;
        }

        if (status === 'running') {
            spinnerText.textContent = getLastLine(data.stderr || '') || 'Scanning...';
        } else if (status === 'completed') {
            spinner.style.display = 'none';
            resultsSection.style.display = 'block';
            scanResults.textContent = data.results || '';
            return;
        } else if (status === 'failed' || status === 'timeout') {
            throw new Error(data.error || 'Scan ' + status);
        }

        await new Promise(r => setTimeout(r, pollIntervalMs));
    }
}

function cleanAnsi(text) {
    return String(text || '')
        .replace(/\x1b\[[0-9;?]*[ -/]*[@-~]/g, '')
        .replace(/\r/g, '')
        .trim();
}

function getLastLine(stderr) {
    const lines = cleanAnsi(stderr).split('\n').filter(l => l.trim());
    return lines.length ? lines[lines.length - 1] : '';
}
```

- [ ] **Step 2: Build to verify no regressions**

```bash
cd /var/home/thetillhoff/code/webscan && go build ./...
```

Expected: no output (success).

- [ ] **Step 3: Commit**

```bash
git add pkg/webserver/static/script.js
git commit -m "feat: rewrite JS for two-page scan flow"
```

---

## Task 6: Markdown mode — inline scan

**Files:**

- Modify: `pkg/webserver/jobs.go`
- Modify: `pkg/webserver/handlers.go`

- [ ] **Step 1: Add `runInlineScan` to `jobs.go`**

Append to `jobs.go` (after `newEngine`):

```go
func (s *Server) runInlineScan(ctx context.Context, target string, follow bool) (string, error) {
 outputBuffer := &synchronizedBuffer{}
 statusBuffer := &synchronizedBuffer{}

 engine, err := s.newEngine(outputBuffer, statusBuffer, follow)
 if err != nil {
  return "", fmt.Errorf("failed to initialize scan engine: %w", err)
 }

 scanCtx, cancel := context.WithTimeout(ctx, s.scanTimeout)
 defer cancel()

 done := make(chan error, 1)
 go func() {
  done <- engine.Scan(target)
 }()

 select {
 case err := <-done:
  if err != nil {
   return "", fmt.Errorf("scan failed: %w", err)
  }
  return outputBuffer.String(), nil
 case <-scanCtx.Done():
  return "", fmt.Errorf("scan timed out after %s", s.scanTimeout)
 }
}
```

- [ ] **Step 2: Replace the `markdownScanHandler` stub in `handlers.go` with the real implementation**

Replace:

```go
// markdownScanHandler stub — replaced in Task 6
func (s *Server) markdownScanHandler(w http.ResponseWriter, _ *http.Request, _ string, _ bool) {
 http.Error(w, "not yet implemented", http.StatusNotImplemented)
}
```

With:

```go
func (s *Server) markdownScanHandler(w http.ResponseWriter, r *http.Request, target string, follow bool) {
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
 fmt.Fprint(w, stripANSI(result))
}
```

- [ ] **Step 3: Run tests and build**

```bash
cd /var/home/thetillhoff/code/webscan && go test ./pkg/webserver/... && go build ./...
```

Expected: all tests PASS, no build errors.

- [ ] **Step 4: Commit**

```bash
git add pkg/webserver/jobs.go pkg/webserver/handlers.go
git commit -m "feat: add markdown mode with inline scan for ?md=1"
```

---

## Task 7: Update e2e test

**Files:**

- Modify: `tests/e2e/webscan.spec.ts`

The existing test references `#scanStatus` which does not exist in the templates. Update the test to use the new two-page flow.

- [ ] **Step 1: Replace the e2e test**

```typescript
import { expect, test } from "@playwright/test";

test("scan thetillhoff.de and show successful keyword-based output", async ({ page }) => {
  await page.goto("/scan?q=thetillhoff.de");

  await expect(page.locator("#resultsSection")).toBeVisible({ timeout: 120_000 });

  const resultsText = await page.locator("#scanResults").innerText();
  expect(resultsText).toContain("# webscan results for thetillhoff.de");
  expect(resultsText).toContain("## DNS scan results");
  expect(resultsText).not.toContain("Error:");
});
```

- [ ] **Step 2: Commit**

```bash
git add tests/e2e/webscan.spec.ts
git commit -m "test: update e2e test for two-page scan flow"
```

---

## Task 8: Final verification

- [ ] **Step 1: Full build and test**

```bash
cd /var/home/thetillhoff/code/webscan && go test ./... && go build ./...
```

Expected: all tests PASS, no build errors.

- [ ] **Step 2: Push**

```bash
git push origin main
```
