# CHANGELOG

##
- Upgrade go version from 1.20 to 1.21
- `webscan version` now prints the currently used version

## v0.2.1
- `webscan` now displays status messages during scans
- Fix bug where scanning ips would trigger dns & ip scan, and vice versa for domain scans

## v0.2.0
- Status code 308 added to valid status codes (only 301 before)
- Now properly checks certificate validity at correct step (not when checking status codes)
- Allow IP addresses (IPv4 and IPv6) as input. If that's the case, dns checking and dns entry retrieval is skipped. Also ipv4 & ipv6 compatibility checks are skipped then.

## v0.1.0
- initial release
- added github actions release workflow
