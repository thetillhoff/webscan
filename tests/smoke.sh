#!/usr/bin/env bash
#
# Live happy-path smoke test for the webscan CLI.
#
# Builds the binary and runs it against real public targets, asserting that
# every scan phase renders and that the HTTP/HTTPS dual-schema output behaves
# correctly (header printed per-schema, content/known-files merged when equal,
# redirect targets scanned via the follow path). Requires outbound network.
#
# Usage: tests/smoke.sh

set -euo pipefail

REPO_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
BIN="$(mktemp -d)/webscan"
FULL_TARGET="${SMOKE_FULL_TARGET:-example.com}"
REDIRECT_TARGET="${SMOKE_REDIRECT_TARGET:-http://cloudflare.com}"

fail() { echo "SMOKE FAIL: $*" >&2; exit 1; }

assert_contains() {
	local file="$1" needle="$2"
	grep -qF "$needle" "$file" || fail "missing '$needle' in scan of $3"
}

echo "building webscan..."
( cd "$REPO_ROOT" && go build -o "$BIN" ./cmd/webscan/ )

# --- Full happy path -------------------------------------------------------
full_out="$(mktemp)"
echo "scanning $FULL_TARGET (full happy path)..."
timeout 120 "$BIN" --quiet "$FULL_TARGET" > "$full_out" 2>/dev/null || fail "scan of $FULL_TARGET exited non-zero"

for section in \
	"## DNS scan results" \
	"## IP scan results" \
	"## TCP port scan results" \
	"## TLS scan results" \
	"## HTTP protocol scan results" \
	"## HTTP header scan results" \
	"## HTTPS header scan results" \
	"## Subdomain scan results" \
	; do
	assert_contains "$full_out" "$section" "$FULL_TARGET"
done

# Header is printed per-schema (never merged); content/known-files merge when
# both schemas return equal results. These guard the emitDualSchema helper.
assert_contains "$full_out" "## HTTP & HTTPS content scan results" "$FULL_TARGET"
assert_contains "$full_out" "Well-known files scan results (HTTP & HTTPS)" "$FULL_TARGET"

# --- Redirect follow path --------------------------------------------------
redir_out="$(mktemp)"
echo "scanning $REDIRECT_TARGET (--follow redirect path)..."
timeout 120 "$BIN" --quiet --follow "$REDIRECT_TARGET" > "$redir_out" 2>/dev/null || fail "scan of $REDIRECT_TARGET exited non-zero"

grep -qF "redirect)" "$redir_out" || fail "expected a redirect result block for $REDIRECT_TARGET"
# The redirect target must be re-scanned via scanWebTarget (protocol + web scans).
awk '/redirect\)/{found=1} found && /## HTTP.* scan results/{ok=1} END{exit ok?0:1}' "$redir_out" \
	|| fail "expected web scans under the redirect block for $REDIRECT_TARGET"

echo "SMOKE PASS"
