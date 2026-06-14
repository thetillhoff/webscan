# Two-Page Layout with URL Query Parameter

## Overview

Split the single-page web UI into two distinct pages with clean URL semantics. The `?q=` parameter carries the scan target and drives routing decisions on both pages.

## Routes

| Route | Condition | Behaviour |
|---|---|---|
| `GET /` | `?q=` absent or empty | Render landing page |
| `GET /` | `?q=` present and non-empty | Redirect 302 → `/scan?q=<value>` |
| `GET /scan` | `?q=` absent or empty | Redirect 302 → `/` |
| `GET /scan` | `?q=` present and non-empty | Render scan/results page |

## Landing Page (`/`)

Minimal centered layout — search bar vertically and horizontally centered in the viewport, similar to a standard search engine home page.

**Elements:**
- Product name / logo above the input
- Text input for the target (domain, IP, or URL)
- "Scan" submit button
- "Follow redirects" checkbox (state persisted in `localStorage`)
- Footer with version and GitHub link (same as today)

**Behaviour:**
- Form `action="/scan"` `method="GET"`, input `name="q"` — no JavaScript required for the basic submit; the browser navigates to `/scan?q=<value>`.
- If `?q=` is non-empty when the page loads, JavaScript immediately redirects to `/scan?q=<value>` (handles the case where someone pastes a query-bearing `/` URL).

## Scan Page (`/scan?q=<target>`)

Functionally identical to today's single page after a scan is submitted, with two additions.

**Elements:**
- Same search bar at the top (not centered — left/top-aligned in a compact header area)
- Input pre-filled with the current `?q=` value
- Spinner, live logs, results, and error sections — unchanged from today
- Footer — unchanged

**Behaviour:**
- On page load, JavaScript reads `q` from `window.location.search` and immediately fires the scan (POST `/api/scan`, then polls — same flow as today's form submit handler).
- If the user edits the input and clicks Scan, the page navigates to `/scan?q=<new-value>` (i.e. `window.location.href` assignment or form submit with `method="GET"`), triggering a fresh page load and a new auto-scan.
- No in-page history manipulation (`pushState`) — navigation is always a full page load so the browser back button reliably returns to `/`.

## Backend Changes

- `server.go`: register `GET /scan` handler.
- `handlers.go`: add `scanPageHandler` — reads `q`, redirects to `/` if empty, otherwise renders the scan page template.
- Update `indexHandler` — reads `q`, redirects to `/scan?q=<value>` if non-empty, otherwise renders the landing page template.
- Two templates: `index.html` (landing) and `scan.html` (scan/results). Shared markup (head, footer) extracted into a `base.html` layout or Go template blocks.

## Static Assets

`style.css` and `script.js` are shared between both pages, served from `/static/` as today.

CSS additions:
- Landing page: `.landing-container` — full-viewport flex column, `justify-content: center`, `align-items: center`.
- Scan page: existing `.scan-container` layout unchanged.

JS changes:
- On `/` load: if `?q=` non-empty, redirect.
- On `/scan` load: read `?q=`, pre-fill input, auto-fire scan.
- Form submit on scan page: navigate to `/scan?q=<value>` instead of calling the API directly.

## Out of Scope

- Preserving `follow` checkbox state in the URL (remains `localStorage` only).
- Any server-side rendering of scan results (results are still fetched client-side via polling).
- Progressive enhancement / no-JS fallback beyond the basic form GET action.
