# CHANGELOG

## v0.2.0
- Status code 308 added to valid status codes (only 301 before)
- Now properly checks certificate validity at correct step (not when checking status codes)
- Allow IP addresses (IPv4 and IPv6) as input. If that's the case, dns checking and dns entry retrieval is skipped. Also ipv4 & ipv6 compatibility checks are skipped then.

## v0.1.0
- initial release
- added github actions release workflow
