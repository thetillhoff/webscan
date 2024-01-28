# CHANGELOG

## v1.2.4
- Update TLS cipher recommendations

## v1.2.3
- Fixed bug where DNS scan results weren't shown with new input type 'domain with path'

## v1.2.2
- Fixed bug where TLS ciphers weren't tested with new input type 'domain with path'
- Fixed bug where newline was printed between http content scan result headline and content

## v1.2.1
- Fixed bug where script urls with new input type 'domain with path' didn't work
- Improved http response/body handling and reduced amount of http clients generated

## v1.2.0
- Updated dependencies
- Added compatibility with new input type 'domain with path' like "github.com/webscan"
- Added valid character scan for cookie headers

## v1.1.0
- Added ipv6 blacklist check
- Added nameserver owner check
- Fixed bug where headline of scan result would be printed without content

## v1.0.0
- Restructured output
- Improved TLS cipher recommendations
- Added domain and ip blacklist search

## v0.3.1
- Upgraded dependencies
- Adjusted pipelines

## v0.3.0
- Upgrade go version from 1.20 to 1.21
- `webscan version` now prints the currently used version
- inputUrl is not stored in webScan.Engine any more, but has to be passed as argument to the Scan functions. It's stored in the Result fields instead.
- Added Verbose flag
  - Verbose mode prints what was the result of identifying the input (domain, ipv4, or ipv6) and other sometimes useful information.
- Following redirects now also applies to following CNAMEs if no A nor AAAA records were detected.
- `PrintScanResults` and all other `Print*` methods no doesn't require any parameters to be called.
- Moved IpVersion compatibility hints from ipScan to dnsScan as that's the correct level of abstraction for such a check.
- Removed `GetCustomDnsServer()` as it's unused after the rework.
- Removed `customDns` variable as it's unused after the rework.
- Only print DNS related information if the input was a domain.
- Moved dnsEngine initialization from scanEngine initialization to webScan initialization.
- Removed duplicate adding of ip address to dnsEngine if input is said ip address.
- Moved httpProtocolScan hint generation to scan method instead of print function.
- Merged `dnsScanEngine` and `dnsScanResults` into one instance of dnsEngine.
- Open ports are now sorted ascending instead of random.
- Inconsistencies of open ports between ip addresses are now detected and printed.
- Scanning open ports of ips is now not only parallelized on ports per ip level, but on ip level, too (all ports of all ips in parallel now).

## v0.2.1
- `webscan` now displays status messages during scans
- Fixed bug where scanning ips would trigger dns & ip scan, and vice versa for domain scans

## v0.2.0
- Status code 308 added to valid status codes (only 301 before)
- Now properly checks certificate validity at correct step (not when checking status codes)
- Allow IP addresses (IPv4 and IPv6) as input. If that's the case, dns checking and dns entry retrieval is skipped. Also ipv4 & ipv6 compatibility checks are skipped then.

## v0.1.0
- initial release
- added github actions release workflow
