# webscan

[![Go Report Card](https://goreportcard.com/badge/thetillhoff/webscan)](https://goreportcard.com/report/thetillhoff/webscan)

**webscan** tries to gather as much information from domains, IPs, and URLs as possible from an external perspective.

## Quick Start

### CLI Mode

```bash
# Install
curl -s https://raw.githubusercontent.com/thetillhoff/webscan/main/install.sh | sh

# Scan a domain
webscan google.com
```

### Web Mode

```bash
# Build web version
go build -o webscan-web ./cmd/webscan-web/

# Run web server
./webscan-web --port 8080

# Open http://localhost:8080 in your browser
```

Open http://localhost:8080, enter any domain, IP, or URL and click "Scan" — or jump straight to a result:

```bash
# Open scan directly in browser
open "http://localhost:8080/scan?q=example.com"

# Get plain-text output (great for scripts)
curl "http://localhost:8080/scan?q=example.com&md=1"
```

## Installation

### CLI Installation

If you're feeling fancy:

```sh
curl -s https://raw.githubusercontent.com/thetillhoff/webscan/main/install.sh | sh
```

If you have `brew` installed:

```sh
brew install thetillhoff/homebrew-tap/webscan
```

or manually from [releases](https://github.com/thetillhoff/webscan/v3/releases/latest).

### Web Installation

To run the web interface:

1. **Build from source:**

   ```sh
   # Clone repo
   git clone https://github.com/thetillhoff/webscan.git
   cd webscan

   # Build web version
   make build-web
   # or manually:
   go build -o webscan-web ./cmd/webscan-web/
   ```

1. **Run the web server:**

   ```sh
   ./webscan-web --port 8080
   ```

1. **Open in browser:** <http://localhost:8080>

## Usage

### CLI

Basic usage:

```sh
webscan google.com            # Scan domain and website
webscan 192.168.0.1          # Scan IP address
webscan https://github.com/thetillhoff/webscan  # Scan specific URL
webscan http://example.com:8080                # Scan specific port
```

### Web Interface

The web interface provides a Google-like search experience across two pages:

- **`/`** — Landing page with centered search bar. Enter a target and click "Scan".
- **`/scan?q=<target>`** — Scan and results page. The scan starts automatically.

**URL parameters:**

| Parameter | Description |
|---|---|
| `q=<target>` | Domain, IP, or URL to scan |
| `follow=1` | Follow CNAMEs and HTTP redirects |
| `md=1` | Return plain-text output instead of the HTML page |

Share or bookmark any scan directly:

```
https://your-instance/scan?q=example.com
https://your-instance/scan?q=example.com&follow=1
```

**Scripting / curl** — `?md=1` returns raw scan output as `text/plain`:

```bash
curl "https://your-instance/scan?q=example.com&md=1"
curl "https://your-instance/scan?q=example.com&follow=1&md=1"
```

All scan features are automatically enabled — no flags needed.

## Project Structure

```text
cmd/
  webscan/           # CLI entry point
  webscan-web/       # Web server entry point
pkg/
  webscan/           # Core scan engine (library)
  webserver/         # Web server implementation
  dnsScan/           # DNS scanning
  tlsScan/           # TLS/SSL scanning
  ...                # Other scan packages
tests/
  e2e/               # End-to-end browser tests (Playwright)
  api/               # API integration tests (Playwright)
  start-test-stack.sh
  playwright.config.ts
```

## Architecture

The web interface uses a hybrid approach:

- **Backend**: Go web server performing all scans
- **Frontend**: Two-page HTML/JS/CSS app served by the Go backend (`/` landing, `/scan` results)
- **API**: `/api/scan` (async, Redis-backed) for browser polling; `/scan?md=1` for synchronous plain-text output

This approach ensures maximum compatibility and feature parity with the CLI while providing a user-friendly web interface.

## Features

### DNS

Display comprehensive DNS information and improvement recommendations:

- WHOIS via RDAP
- Nameserver ownership verification
- CNAME resolution
- Records overview (A, AAAA, MX, TXT, etc.)
- IPv6 readiness checks
- Domain blacklist detection
- DNS validators and key values

### IP Analysis

- IPv4 and IPv6 address ownership via RDAP
- Hoster identification (AWS, Azure, GCP, etc.)
- Blacklist checks for all discovered IPs

### Port Scanning

- TCP port scanning for common services
- Parallel checking with configurable timeout per port
- Cross-IP port consistency verification
- Special handling for HTTP/HTTPS ports

### TLS/SSL Checks

- Certificate validity and chain verification
- Cipher suite recommendations
- TLS version checks (1.2, 1.3 recommended)
- Perfect Forward Secrecy verification
- Deduplicated output when multiple IPs share the same certificate

### HTTP Detection

- HTTP/HTTPS availability checks
- Redirect behavior analysis
- Protocol version detection (HTTP/1.1, HTTP/2, HTTP/3)
- Parallel protocol checks with configurable timeout

### HTTP Headers

- Security header analysis (HSTS, CSP, etc.)
- Cookie inspection
- Host header recommendations
- Compression verification

### HTML Content

- HTML validation
- Resource analysis (CSS, JS, images)
- Size recommendations
- Accessibility hints
- Standards compliance checks

### Well-Known Files

- Checks for standard files: robots.txt, sitemap.xml, security.txt, llms.txt, AI plugin manifest
- Detects exposed sensitive files: .htaccess, .env, .git/config, wp-config.php, server-status

## Development

### Build Commands

```bash
# Build CLI version
make build
# or: go build -o webscan ./cmd/webscan/

# Build web version
make build-web
# or: go build -o webscan-web ./cmd/webscan-web/

# Run web version
./webscan-web --port 8080
```

### Makefile Targets

- `run <target>`: Run CLI with arguments (e.g. `make run thetillhoff.de`)
- `run-web <args>`: Run web server with arguments
- `build`: Standard CLI binary
- `build-web`: Web server with interface
- `test`: Run unit tests
- `lint`: Run Go vet, build, and markdownlint
- `format`: Code formatting
- `upgrade`: Dependency updates
- `compose-start`: Start Docker Compose stack
- `compose-stop`: Stop Docker Compose stack

## Contributing

See [DEVELOPMENT.md](DEVELOPMENT.md) for development and contribution guidance.

## License

MIT License - See [LICENSE](LICENSE) for details.

## Contact

Questions, issues, or contributions welcome:

- GitHub: <https://github.com/thetillhoff/webscan>
- Issues: <https://github.com/thetillhoff/webscan/v3/issues>

## Roadmap

See [TODO.md](TODO.md) for the full list of planned features, known bugs, and improvement ideas.
