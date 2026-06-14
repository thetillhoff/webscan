# Two-Page Layout with URL Query Parameter

## Overview

Split the single-page web UI into two distinct pages with clean URL semantics. Three URL parameters drive routing and behaviour: `q` (scan target), `follow` (follow redirects), and `md` (markdown output).

## URL Parameters

| Parameter | Values | Meaning |
|---|---|---|
| `q` | any non-empty string | Scan target (domain, IP, or URL) |
| `follow` | `1` | Follow CNAMEs and HTTP redirects |
| `md` | `1` | Return raw markdown output instead of the HTML page |

## Routes

| Route | Condition | Behaviour |
|---|---|---|
| `GET /` | `?q=` absent or empty | Render landing page |
| `GET /` | `?q=` present and non-empty | Redirect 302 → `/scan?q=<value>[&follow=1]` |
| `GET /scan` | `?q=` absent or empty | Redirect 302 → `/` |
| `GET /scan` | `?q=` present, `md` absent | Render scan/results HTML page, auto-start scan |
| `GET /scan` | `?q=` present, `?md=1` | Run scan, return raw output as `text/plain` |

`follow=1` is preserved on the `/` → `/scan` redirect when present.

## Landing Page (`/`)

Minimal centered layout — search bar vertically and horizontally centered in the viewport, similar to a standard search engine home page.

**Elements:**
- Product name / logo above the input
- Text input for the target (domain, IP, or URL), `name="q"`
- "Follow redirects" checkbox, `name="follow"` `value="1"`
- "Scan" submit button
- Footer with version and GitHub link (same as today)

**Behaviour:**
- Form `action="/scan"` `method="GET"` — no JavaScript required for the basic submit; the browser navigates to `/scan?q=<value>[&follow=1]`.
- If `?q=` is non-empty when the page loads, JavaScript immediately redirects to `/scan?q=<value>[&follow=1]` (handles the case where someone pastes a query-bearing `/` URL).
- `follow` checkbox initial state loaded from `localStorage` (same as today); on change, state saved back to `localStorage`.

## Scan Page (`/scan?q=<target>`)

Functionally identical to today's single page after a scan is submitted.

**Elements:**
- Search bar at the top (compact, not centered)
- Input pre-filled with the current `?q=` value
- "Follow redirects" checkbox, pre-checked if `?follow=1` is in the URL
- Spinner, live logs, results, and error sections — unchanged from today
- Footer — unchanged

**Behaviour:**
- On page load, JavaScript reads `q` and `follow` from `window.location.search` and immediately fires the scan (POST `/api/scan` with `{ target, follow }`, then polls — same flow as today).
- If the user edits the input or toggles the checkbox and clicks Scan, the page navigates to `/scan?q=<new-value>[&follow=1]`, triggering a fresh page load and a new auto-scan.
- No in-page history manipulation (`pushState`) — navigation is always a full page load so the browser back button reliably returns to `/`.
- `follow` checkbox state is no longer persisted to `localStorage` on the scan page — the URL is the source of truth.

## Markdown Mode (`/scan?q=<target>&md=1`)

Returns the scan output as `text/plain` with no HTML wrapper. Intended for CLI use (`curl`) or scripting.

**Behaviour:**
- Handler reads `q` and `follow`, runs the scan synchronously (blocking until complete or timeout), and writes the raw output directly to the response.
- No polling, no Redis job — the scan runs inline in the request goroutine.
- Response `Content-Type: text/plain; charset=utf-8`.
- HTTP status 200 on success, 400 if `q` is missing, 504 if the scan times out.
- ANSI colour codes stripped from output (same helper already used in the JS client).

## Backend Changes

- `server.go`: register `GET /scan` handler.
- `handlers.go`:
  - Update `indexHandler` — redirects to `/scan?q=...` if `q` non-empty.
  - Add `scanPageHandler` — reads `q`; if `md=1`, runs scan inline and returns plain text; otherwise renders `scan.html`.
- Two templates: `index.html` (landing) and `scan.html` (scan/results). Shared markup (head, footer) extracted into Go template blocks.

## Static Assets

`style.css` and `script.js` are shared between both pages, served from `/static/` as today.

CSS additions:
- Landing page: `.landing-container` — full-viewport flex column, `justify-content: center`, `align-items: center`.
- Scan page: existing `.scan-container` layout unchanged.

JS changes:
- On `/` load: if `?q=` non-empty, redirect to `/scan?q=...`.
- On `/scan` load: read `q` and `follow` from URL, pre-fill inputs, auto-fire scan.
- Form submit on scan page: navigate to `/scan?q=<value>[&follow=1]`.

## Out of Scope

- Any server-side rendering of scan results for the HTML page (results still fetched client-side via polling).
- Progressive enhancement / no-JS fallback beyond the basic form GET action.
- Streaming the markdown response (output is buffered until scan completes).
