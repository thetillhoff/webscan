# CHANGELOG

## v5.1.0

### New Features

- **Domain and IP Blocklist**: `DOMAIN_BLOCKLIST` and `IP_BLOCKLIST` env vars block scan targets matching specified domain suffixes or IP CIDR ranges; prevents scanning of internal cluster services and private IP space

### Improvements

- **Redirect section headers**: Each followed redirect now starts with a `---` separator and a header showing the target URL and the HTTP status code that triggered the follow (e.g. `HTTP 301 redirect`)

## v5.0.0

### Breaking Changes

- **Multi-Binary Layout**: CLI entry point moved from main.go to cmd/webscan/main.go
- **DNS Resolver**: Fails explicitly when no system resolver is found instead of silently falling back to a public resolver

### New Features

- **Web Server Mode**: New browser-based scanning interface with Redis-backed job queue (cmd/webscan-web/, pkg/webserver/); distributed as a Docker image, not a standalone binary
- **Docker Compose**: Full stack deployment with Redis and web server
- **Well-Known Files Scan**: New --files flag checks for standard files (robots.txt, sitemap.xml, security.txt, llms.txt, AI plugin manifest) and warns on exposed sensitive files (.htaccess, .env, .git/config, wp-config.php, server-status)
- **Unified Timeout Flag**: New --timeout flag (default 5s) controls all network requests (DNS, port scan, HTTP, RDAP, blacklist checks); previously hardcoded per scan type

### Improvements

- **HTTP Redirect Detection**: Protocol scan now correctly detects HTTP redirects (301/302/303/307/308) and displays the redirect target and status code
- **Follow Redirect Chain**: --follow now runs web scans (protocol, headers, content) on each redirect target in sequence, without repeating DNS/IP/port/TLS scans
- **Parallel HTTP Protocol Checks**: HTTP/1, HTTP/2, HTTP/3 version checks and redirect detection all run concurrently per port and IP
- **Structured Header Output**: HTTP header scan results now use a structured HeaderEntry type; each header shows its value and recommendation inline with a -> prefix; multi-line values (e.g. CSP) are indented cleanly
- **HTTP & HTTPS Section Deduplication**: Content scan and well-known file scan results are printed once as "HTTP & HTTPS" when both protocols return identical results
- **TLS Labeled Certificate Names**: Certificate names are now prefixed with SN: (Subject Name) or SAN: (Subject Alternative Name) and always shown when the TLS scan runs
- **TLS Cipher Rule Titles**: Each cipher rule now has a short title (e.g. "RC4 ciphers") printed before its description, making sections immediately identifiable
- **TLS Cipher Rule Ordering**: Cipher rules are now printed in definition order (deterministic) instead of map iteration order
- **TLS Cipher Cross-Rule Deduplication**: Ciphers already flagged by a specific rule (RC4, 3DES, CBC) no longer also appear under "Ciphers deemed insecure by Golang"
- **Well-Known File Analysis**: robots.txt, security.txt, and sitemap.xml are parsed and analyzed — robots.txt checks for sitemap directives and overly broad disallow rules; security.txt validates required fields and expiry; sitemap.xml reports URL and sub-sitemap counts
- **TLS Output Deduplication**: Shared certificate info printed once; only per-IP differences shown
- **IP Blacklist Codes**: Comprehensive Spamhaus return code handling with human-readable descriptions; error codes (rate limiting, public resolver) logged as warnings instead of confusing the user
- **crt.sh Error Handling**: Certificate transparency log failures (timeout, rate limit, server error, unreachable) shown as a brief note instead of raw errors
- **Status Non-TTY Fix**: Status messages print to output when not running in a terminal (previously silently dropped)
- **Playwright Tests**: E2E and API test scaffolding in tests/

### Bug Fixes

- **Referrer-Policy Scan**: Fixed incorrect check of Referer (a request header) — now correctly checks Referrer-Policy (the response header that controls referrer behaviour)
- **HSTS for HTTP**: No longer recommends HSTS headers for plain HTTP responses (only HTTPS)
- **TLS Cipher Deduplication**: Cipher rules that are identical across IPs are no longer printed per-IP
- **External Resource Fetch**: HTML content scan now handles unreachable external stylesheets and scripts gracefully instead of aborting the scan; error reason (e.g. DNS returning 0.0.0.0) is included in the output
- **Redirect Output**: Protocol scan no longer emits a redundant "traffic is redirected to" line; redirect info is shown in the scan summary
- **Subdomain Scan Flags**: --subdomains no longer triggers advanced DNS or prints TLS results; TLS scan runs internally for SANs only
- **Protocol Scan Flags**: --protocol no longer triggers advanced DNS scan

### Dependencies

- Added github.com/redis/go-redis/v9 for job queue
- Updated dependencies

## v4.4.0

### Bug Fixes

- Fixed silent exit when no target argument is provided: `log.Fatalln` was silently dropped because Go 1.21+ bridges `log.*` through slog, and the message level (Info) was below the default threshold (Warn); the check is now handled in `types.NewTarget` and surfaces as a visible error message

## v4.3.0

### Improvements

- Unified functional options pattern across all scan packages via `types.Option[C]` and `types.ApplyOptions` generics; each package now uses a type alias `ConfigOption = types.Option[scanConfig]` instead of a redundant inline definition

## v4.2.0

### Breaking Changes

- Renamed exported functions to follow Go naming conventions for acronyms (HTTP, TLS, DNS, IP, HTML):
  - `EnableTlsScan()` → `EnableTLSScan()`
  - `EnableHttpProtocolScan()` → `EnableHTTPProtocolScan()`
  - `EnableHttpHeaderScan()` → `EnableHTTPHeaderScan()`
  - `EnableHtmlContentScan()` → `EnableHTMLContentScan()`
  - `EnableIpScan()` → `EnableIPScan()`
  - `cachedHttpGetClient`: `GetHttpResponse()` → `HTTPResponse()`, `GetBody()` → `Body()`, `GetError()` → `Err()`, `DoesVerifyTls()` → `VerifyTLS()`, `GetTimeout()` → `Timeout()`
  - `dnsScan`: `CheckDkim()` → `CheckDKIM()`, `CheckDmarc()` → `CheckDMARC()`, `CheckSpf()` → `CheckSPF()`, `CheckIpVersionCompatibility()` → `CheckIPVersionCompatibility()`
  - `htmlContentScan`: `ValidateHtml()` → `ValidateHTML()`
  - `httpProtocolScan`: `CheckHttpRedirects()` → `CheckHTTPRedirects()`, `CheckHttpVersions()` → `CheckHTTPVersions()`
  - `tlsScan`: internal functions aligned with the same convention

### Bug Fixes

- Fixed crash when `-vv` verbosity flag was used (any verbosity ≥ 2 now enables debug level)
- Fixed CLI flag name mismatches that caused several scan flags to be silently ignored
- Fixed TLS connection leak in `tlsScan`: connections were opened to probe cipher/version support but never closed
- Fixed HTTP response body leak in `cachedHttpGetClient`: body was read but not closed
- Fixed `package-level sync.WaitGroup` in `tlsScan` that was unsafe for concurrent use; replaced with a local WaitGroup

### Improvements

- Consolidated `IsIPv4` and `IsIPv6` helpers into `pkg/types` (removed duplicate copies in `dnsScan` and `ipScan`)
- Removed unused `Scanner` generic interface from `pkg/webscan`
- Made `GetHttpProtocolRecommendationsForResult` unexported (only used within `httpProtocolScan`)
- Added Windows portability: DNS config no longer attempts to read `/etc/resolv.conf` on Windows
- Standardised all `slog` messages to `"<package>: <Description>"` format with structured key-value pairs


## v4.1.0

### New Features

- Added docker image
- Automatic dependency updates when available.
  They'll trigger a new patch-release.

## v4.0.4

Updated dependencies.

## v4.0.3

### Improvements

- Add a basic test that the binary to be released executes correctly and prints the correct version.

## v4.0.2

### Improvements

- Adjusted install.sh to use temporary directory for downloads
- Renovate now uses branch type automerges to reduce notifications to watchers.
- Added github action that automatically triggers a patch release when renovate automerges minor or patch versions of dependencies.
- Updated dependencies.

## v4.0.1

### Improvements

- Adjusted some help messages

### Bug Fixes

- Fixed bug where install.sh would fail to verify the checksum of the downloaded file
- Fixed bug where help would not show usage samples
- Fixed bug where `--follow` flag would be shown twice in help
- Fixed bug where `completion` subcommand would be shown twice in help

## v4.0.0

### Breaking Changes

- **Instant Flag Removal**: Made the `--instant` flag default and removed it, as there is no disadvantage on using it
- **Quiet Flag Handling**: Moved handling of the `--quiet` flag to the cli, as it makes no sense to have the library handle (to just print debug messages, but no results).
  The `--quiet` flag now discards all output that would be printed to stdout.

### New Features

- **CNAME Following**: Added `--follow` flag to follow CNAMEs
- **Schema Support**: Added support for schema prefixes `http://` and `https://` in input
- **Port Specification**: Added support for port suffixes like `:80`, `:443` or `:8080` in input
- **Path Support**: Added support for path suffixes like `/path` in input
- **Multi-Record TLS Scan**: tlsScan now checks all A/AAAA records of the target, not just the first one

### Improvements

- **DNS System Integration**: DNS client now uses system nameservers from resolv.conf on Unix systems with fallback to public DNS
- **Cross-Platform DNS Support**: Added support for macOS and Windows DNS configuration
- **Nameserver Owner Detection**: Fixed RDAP lookup for nameserver hostnames
- **Subdomain Scan**: Added filters to only show subdomains of the actual target domain when going through the certificate SAN list
- **IPv6 Nameserver Support**: Fixed IPv6 address formatting in DNS queries
- **Library I/O Architecture**: Configured the libraries to use io.Writer instead of fmt.Println, so the libraries are more versatile and can be used in other projects
  The cli still hands in os.Stdout and os.Stderr by default
- **API Consistency**: Adjusted scan-modules so they have a more consistent api
- **Output Formatting**: Adjusted formatting of result outputs, so they are more consistent and easier to read
- **Debug Logging**: Added more debug log statements

### Bug Fixes

- **Codebase Restructuring**: Fixed several large and small bugs, restructured most of the codebase to make it more consistent and easier to maintain

### Dependencies

- **Package Updates**: Updated dependencies

## v3.0.10

### Bug Fixes

- Fix formatting of logs of http protocol scan

## v3.0.9

### Bug Fixes

- Fixed bug, where ip blacklisting error/warning would break formatting with its error message
- Fixed bug, where cli-args were not picked up correctly, and the old `--all` was implicitely active at all times
- Fixed bug, where results of http-protocol scan were not printed
- Fixed bug, where the http-protocol scan would be wrong at all times for multiple reasons

## v3.0.8

### Bug Fixes

- Fixed support for ipv4 and ipv6 as input
- Fixed bug, where subdomainResults contained ip addresses if they were in the certificate SAN list

### Dependencies

- Updated dependencies

## v3.0.7

### Bug Fixes

- Fixed bug where subdomain scan would fail silently if crt.sh was not reachable

## v3.0.6

### Bug Fixes

- Revert change from v3.0.5 as the bug did not exist and was a local issue

## v3.0.5

### Bug Fixes

- Fixed bug where version was not printed to output of `webscan version`

## v3.0.4

### Documentation

- Updated readme

### Dependencies

- Updated dependencies

## v3.0.3

### Bug Fixes

- Fixed bug, where tlsScan failed for urls with hostname and path, like `abc.de/path`
- Fixed bug, where htmlContentScan couldn't read the response body correctly

## v3.0.2

### Breaking Changes

- Adjusted module path to conform to go.mod spec, where the module path needs to contain the major version (aka `.../webscan/v3` instead of just `.../webscan`)

## v3.0.1

### Improvements

- Reduced output if everything is alright, but a scan could not check for information that does not exist

## v3.0.0

### Breaking Changes

- Moved from `spf13/cobra` and `spf13/viper` to `urfave/cli` as cli-library due to maintenance issues

### New Features

- Added progress updates with spinner and - depending on type of scan `X/total` status updates with numbers
- Added list of SN & SAN to output of TLS scan
- Added display of Server header if response contains it

### Improvements

- Don't show other sizes if html body size is 0 already
- Automatically format file sizes to use kB if >1000B or B otherwise
- Don't show scan results for web scans on http or https if the respective ports are not open - there will not be anything to show either way
- Grouped TLS cipher issues by rule instead of just tuples of cipher and rule as before
- Set version variable during build of release

## v2.0.3

### Documentation

- Added openssf badge - login delayed ... see <https://github.com/coreinfrastructure/best-practices-badge/issues/2150>

### Build

- Added Makefile to support brew in the future

## v2.0.2

### Bug Fixes

- Disabled automatic mail-config and sub-domain scan for implication of `-a` (run all scans)

## v2.0.1

### Bug Fixes

- Fixed bug where implication of `-a` (run all scans) did not work

## v2.0.0

### Breaking Changes

- Removed `-a` (run all scans) argument and made it the default configuration if no specific scans are enabled

## v1.2.7

### Documentation

- Added goreportcard badge with automatic refresh on releases

## v1.2.6

### New Features

- Added support for `arm64` ARCH type

### Build

- Added OS and ARCH verification to `install.sh` script

### Dependencies

- Updated dependencies

## v1.2.5

### Dependencies

- Updated dependencies

## v1.2.4

### Improvements

- Updated TLS cipher recommendations

## v1.2.3

### Bug Fixes

- Fixed bug where DNS scan results weren't shown with new input type 'domain with path'

## v1.2.2

### Bug Fixes

- Fixed bug where TLS ciphers weren't tested with new input type 'domain with path'
- Fixed bug where newline was printed between http content scan result headline and content

## v1.2.1

### Bug Fixes

- Fixed bug where script urls with new input type 'domain with path' didn't work

### Improvements

- Improved http response/body handling and reduced amount of http clients generated

## v1.2.0

### New Features

- Added compatibility with new input type 'domain with path' like "github.com/webscan"
- Added valid character scan for cookie headers

### Dependencies

- Updated dependencies

## v1.1.0

### New Features

- Added ipv6 blacklist check
- Added nameserver owner check

### Bug Fixes

- Fixed bug where headline of scan result would be printed without content

## v1.0.0

### Improvements

- Restructured output
- Improved TLS cipher recommendations

### New Features

- Added domain and ip blacklist search

## v0.3.1

### Dependencies

- Upgraded dependencies

### Build

- Adjusted pipelines

## v0.3.0

### Breaking Changes

- Upgrade go version from 1.20 to 1.21
- inputUrl is not stored in webScan.Engine any more, but has to be passed as argument to the Scan functions. It's stored in the Result fields instead
- `PrintScanResults` and all other `Print*` methods no doesn't require any parameters to be called
- Removed `GetCustomDnsServer()` as it's unused after the rework
- Removed `customDns` variable as it's unused after the rework

### New Features

- `webscan version` now prints the currently used version
- Added Verbose flag
  - Verbose mode prints what was the result of identifying the input (domain, ipv4, or ipv6) and other sometimes useful information
- Following redirects now also applies to following CNAMEs if no A nor AAAA records were detected

### Improvements

- Moved IpVersion compatibility hints from ipScan to dnsScan as that's the correct level of abstraction for such a check
- Only print DNS related information if the input was a domain
- Moved dnsEngine initialization from scanEngine initialization to webScan initialization
- Removed duplicate adding of ip address to dnsEngine if input is said ip address
- Moved httpProtocolScan hint generation to scan method instead of print function
- Merged `dnsScanEngine` and `dnsScanResults` into one instance of dnsEngine
- Open ports are now sorted ascending instead of random
- Inconsistencies of open ports between ip addresses are now detected and printed
- Scanning open ports of ips is now not only parallelized on ports per ip level, but on ip level, too (all ports of all ips in parallel now)

## v0.2.1

### New Features

- `webscan` now displays status messages during scans

### Bug Fixes

- Fixed bug where scanning ips would trigger dns & ip scan, and vice versa for domain scans

## v0.2.0

### New Features

- Status code 308 added to valid status codes (only 301 before)
- Allow IP addresses (IPv4 and IPv6) as input. If that's the case, dns checking and dns entry retrieval is skipped. Also ipv4 & ipv6 compatibility checks are skipped then

### Improvements

- Now properly checks certificate validity at correct step (not when checking status codes)

## v0.1.0

### New Features

- initial release
- added github actions release workflow
