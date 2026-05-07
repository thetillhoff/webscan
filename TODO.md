# TODO

## CI / Build

- GitHub Actions workflows use thetillhoff/action-golang-build which always builds from repo root. Needs a PACKAGE input or inline go build ./cmd/webscan/ instead.
- Add version buildargs to example usage in all three GitHub Actions repos
- Add integration tests to the release pipeline (sample scans on multiple platforms)

## Bugs

- Subdomain list shows all domains in certificates, not just actual subdomains of the target
- DNS checks and dials time out on Windows

## DNS / Mail Security

- DNSSEC validation
- SPF: verify domain-spec per spec; recursive check of external references
- DKIM: TXT variant verification; CNAME variant recursive check
- DMARC: TXT variant verification; CNAME variant recursive check
- MX blacklist detection
- Warn on multiple CNAME redirects; detect CNAME loops
- DNS best practices (TTL values, SOA)
- Check both IPv4 and IPv6 MX records (follow CNAMEs on MX automatically)

## TLS / SSL

- Recommend TLS 1.3 support (currently only warns on 1.0/1.1)
- Add certificate expiry date to scan results
- Check certificate issuer — recommend free CAs like Let's Encrypt over paid ones
- Tests against badssl.com
- Compare against known TLS configs (AWS TLS policies) for specific recommendations
- Post-quantum cipher verification
- Perfect Forward Secrecy ciphers: enable the 2026 warning (code exists but is time-gated)

## HTTP Protocol

- Check that redirect locations end with / or a filename
- Note that HTTP/1.1 is considered insecure (site redirection attacks)

## HTTP Headers

- CSP: check for nonce/hash for scripts; verify self/https usage
- Format header recommendations as structured problem/recommendation/link items

## HTML Content

- CSS/JS validation and minification checks
- Image scanning (size < 500kB, webp format recommendation)
- HTML accessibility checks (relative URLs, mixed content, viewport, font sizes, tap targets)
- Media embedding recommendations
- Check for /.well-known/security.txt
- Headless rendering with chromedp (performance index, console errors)
- Content type-aware validation (e.g. JSON for application/json responses)
- List all referenced external domains (fonts.google.com, etc.)

## SEO

- robots.txt check
- sitemap check
- Detect incompatible plugins (Flash)

## Subdomain Scan

- Reverse IP lookup (multiple sources: hackertarget, Bing, nmmapper)
- Additional cert transparency sources (certspotter)
- Third-party databases (virustotal, securitytrails, dnsdumpster)
- Third-party APIs (Amass, PassiveTotal, Shodan)
- Check response bodies in HTTP cache for referenced subdomains
- Crawl the website itself for subdomain references

## Port Scan

- Check FTP disabled (only SFTP/FTPS)
- Check SSH password auth disabled / secure configuration
- Scan relevant UDP ports
- Check version strings in TCP greeting/headers (e.g. OpenSSH)
- Check FTP headers
- Note about IPv4/IPv6 port inconsistencies when local machine only supports one

## IP Analysis

- Add estimated geolocation and ASN info to IPs
- Reverse DNS lookup for IP addresses
- Check if QUIC uses UDP on port 443 and incorporate into scans

## Output and UX

- --json-output flag for structured JSON output
- Merge HTTP/HTTPS header results when identical; skip if redirecting
- Merge HTTP/HTTPS content results when identical; skip if redirecting
- Format TLS results as structured problem/recommendation/link/affected-ciphers items
- Support reading target from stdin
- Add repo URL and issues link to help text
- Print IP RDAP info formatted by longest IP address width
- Subdomain list in alphabetical order; show wildcard entries distinctly (dimmed/italic)
- Check HTTP latency, hops, download speed

## Logger / Status Packages

- Colored completion indicators (green/yellow/red checkmarks)
- Compatibility with piping/stdin/redirect to file
- FAIL level with error code, message, remediation, and link
- Respect ALL_PROXY, HTTP_PROXY, HTTPS_PROXY, NO_PROXY
- Use $LINES and $COLUMNS for output formatting

## Architecture / Code Quality

- Remove "opinionated" flag — either it's a valid recommendation or not
- Find better way to pass results between scans (shared context vs arguments)
- Use types/enums/constants for recommendations (make http/https comparable)
- Add unit tests (stdout result, stderr result, pipe/file output)
- Target type needs tests for each parsing case including error cases

## Web Interface

- Rate limiting (IP-based)
- Proper CORS configuration
- CSRF protection
- Security headers (HSTS, XSS prevention, etc.)
- Kubernetes Helm chart
- Configuration file / env var support
