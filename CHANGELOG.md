# CHANGELOG

## WIP

### Bug Fixes

- Fixed bug where install.sh would fail to verify the checksum of the downloaded file

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
