# webscan

[![Go Report Card](https://goreportcard.com/badge/thetillhoff/webscan)](https://goreportcard.com/report/thetillhoff/webscan)

Webscan tries to gather as much information from domains, IPs, and URLs as possible from an external perspective.
It covers

- DNS configuration
- Domain and Nameserver ownerships
- IPv4 and IPv6 availability
- IP address ownerships
- Blacklisting status
- Open ports
- TLS validity
- TLS configuration safety
- HTTP/HTTPS configuration with redirects
- host-headers
- cookies
- html, js, css sizes
- subdomains
- and much more!

of a specified url or ip and gives improvement recommendations based on best-practices.

## Installation

If you're feeling fancy:

```sh
curl -s https://raw.githubusercontent.com/thetillhoff/webscan/main/install.sh | sh
```

If you have `brew` installed:

```sh
brew install thetillhoff/homebrew-tap/webscan
# This tap is updated once per day. So it can take up to 24h for new versions to propagate to it.
```

or manually from <https://github.com/thetillhoff/webscan/v3/releases/latest>.

## Usage

```sh
webscan google.com # Scan domain and website
webscan 192.168.0.1 # Scan IP and website
webscan https://github.com/thetillhoff/webscan # Scan website at default port according to schema
webscan http://example.com:8080 # Scan website at specific port

webscan --help # Learn more about running specific scans
```

## Features

### DNS

Display dns information about the provided URL, and give improvement recommendations.

- [x] This is skipped if the input is an ipv4 or ipv6 address
- [x] Check who is the owner of the Domain via RDAP (not supported for country-TLDs)
- [x] Check who is the owner of the DNS zone (== nameserver owner)
- [x] Follow CNAMEs
- [x] DNS records overview
- [x] Long DNS name (requires `--opinionated` flag)
- [ ] DNS best practices (TTL values, SOA, ...)
- [ ] DNSSEC
- [ ] Warn if there are more than one CNAME redirects
- [ ] Detect CNAME loops
- [x] Domain blacklist detection
- [x] Scan DNS of domain even if input is domain with path (like "github.com/webscan")
- [x] Specify a custom dns server with the `--dns <dns server location>` option.
- [x] Use system DNS configuration on Unix systems (resolv.conf) by default with fallback to public DNS
- [x] Use DNS library `miekg/dns`, as it's more versatile than the default one. It can for example NOT follow CNAMEs when looking up other records.

### DNS mail security

- [~] SPF
  - [x] Verify usage of allowed mechanisms, modifiers, qualifiers.
  - [x] Verify contained ipv4 addresses.
  - [x] Verify contained ipv6 addresses.
  - [ ] Verify `<domain-spec>` as described in spec.
  - [ ] Recursive check of any referenced external information.
- [~] DKIM
  - [x] TXT variant detection
  - [x] CNAME variant detection
  - [ ] TXT variant verification
  - [ ] CNAME variant recursive check
- [~] DMARC
  - [x] TXT variant detection
  - [x] CNAME variant detection
  - [ ] TXT variant verification
  - [ ] CNAME variant recursive check
- [ ] MX blacklist detection

### Subdomain finder

- [x] Search for subdomains of the provided domains and provide a list of them.
  - [x] Search for subdomains listed in certificates logs (explanation at <https://certificate.transparency.dev/>, searched at <https://crt.sh>).
- [x] Search for subdomains in the subject and alternate name list of the original domain tls certificate.
- [ ] Check other DNS entries (like PTR), certificate pointers, SPF record, certificate logs, reverse-ip-lookups
  - [ ] reverse ip lookup with <https://hackertarget.com/reverse-ip-lookup/>
  - [ ] Bing reverse ip lookup
  - [ ] <https://www.nmmapper.com/sys/tools/subdomainfinder/>
  - [ ] <https://dnsdumpster.readthedocs.io/en/latest/#>
  - [ ] crawl the website itself and the links it has for subdomains
  - [ ] second source for cert-transparency-logs, f.e. <https://github.com/SSLMate/certspotter>
  - [ ] other 3rd party databases like virustotal, securitytrails, and dnsdumpster
  - [ ] 3rd party subdomain apis like Amass, PassiveTotal, and Shodan

### IPv6 readiness

- [x] Check if both ipv4 and ipv6 entries were found.
  - IPv4 is necessary to stay backwards compatible.
  - IPv6 is recommended to be IPv6 ready.

### IP analysis

- [x] Check who is the hoster of the IP address via RDAP (successor of whois) - like AWS, Azure, GCP, ...
- Check if any IP (v4 and v6) of the domain is blacklisted.
  - [x] IPv4
  - [x] IPv6

### Open ports

- [x] Check all found ipv4 and ipv6 entries for relevant open ports. Examples for relevant ports are SSH, FTP, SMB, SMTP, HTTP, POSTGRES
  - [ ] Check whether FTP is disabled completely (only use SFTP or FTPS)
  - [ ] Check whether SSH has password auth disabled and uses a secure configuration
- [x] Check ports in parallel, since the connection timeout is set to 2s, which can add up quite much.
- [x] Check if open ports match across all IPs.
- [x] If http detection feature is enabled, check HTTP and HTTPS ports even if this feature is not enabled.

### SSL/TLS check

- [x] Validate certificate only if port 443 is open.
- [x] Check the validity of the ssl certificate. Subject, date, chain, ciphers, tls-min-version (if valid, but not recommended, don't fail, but print warning/s instead).
- [ ] Write tests against badssl.com

- [ ] SSL is not recommended
- [x] TLS 1.0 and TLS 1.1 are not recommended, only TLS 1.2 & 1.3 are okay
- [ ] TLS 1.3 should be supported

- cipher recommendations like

  - [x] Recommending against 3DES, as it's vulnerable to birthday attacks (<https://sweet32.info/>).
  - [x] Recommending against RC4, as it's exploitable biases can lead to plaintext recovery without side channels (<https://www.rc4nomore.com/>).
  - [x] Recommending against CBC, as it seems fundamentally flawed since the Lucky13 vulnerability was discovered (<https://en.wikipedia.org/wiki/Lucky_Thirteen_attack>).
  - [x] Keep in mind `ECDH_` ciphers don't support Perfect Forward Secrecy and shouldn't be used after 2026.
  - [x] Keep in mind `DH_` ciphers don't support Perfect Forward Secrecy and shouldn't be used after 2026.

- [ ] check who's the issuer of the certificate. If it's one of the most known paid providers, recommend to use a free one, like letsencrypt.

### HTTP detection

By default `webscan` assumes you're using https. Yet, it will check whether it's available via http as well.

- [ ] Optionally follow HTTP redirects (status codes 30x)
- [x] If http is available, it should be used for redirects only.
- [x] If https is availabe, it should either redirect or respond with 200 status code.
- [x] If both http and https are available
  - [x] and https redirects, check that either http redirects to https or to the same location that https redirects to
  - [x] and both are redirecting, the destination should be a https location
  - [x] and https does not redirect, http should redirect to it with 301 or 308 status code. (<https://kubernetes.github.io/ingress-nginx/user-guide/nginx-configuration/configmap/#http-redirect-code>)
- [x] Check which http versions are supported by the webserver (like HTTP/1.1, HTTP/2, HTTP/3 aka QUIC)

### HTTP headers

- [x] Analyze host-headers of the response and recommend best-practices.
- [x] Check HTTP Strict Transport Security (HSTS / STS). This defeats attacks such as SSL Stripping, and also avoids the round-trip cost of the 301 redirect to Redirect HTTP to HTTPS.
- [x] CSP header settings (have one is the minimum requirement here)
  - [ ] use nonce or hash for script
  - [ ] use self, or at least https, warn for everything else -> <https://storage.googleapis.com/pub-tools-public-publication-data/pdf/45542.pdf>
- [x] Scan cookies
  - [x] amount
  - [x] length
  - [x] used characters

### HTML content

Print recommendations on the html code.

- [x] Scan HTML contents of path even if input is domain with path (like "github.com/webscan")
- [x] Check if compression is enabled
- [x] HTML validation -> It has to be parsable to look further into the details
- [x] HTML5 check if size of body > 0
  - [x] `<!DOCTYPE html>` is first node
  - [x] `<html lang="en-US">` `html` node has `lang` attribute set (it's value isn't validated)
  <!-- - [ ] `<meta charset="utf-8"/>` charset is set -->
- [~] HTML parsing - requires go-colly for finding references to other files

  - [~] check html
    - [x] size < 200kb
    - [x] validation
    - [ ] minification
    - [x] html5 validation
  - [~] check css (if size of html > 0)
    - [x] size
    - [ ] validation
    - [ ] minification
    - [ ] feature support recommendation
  - [~] check js (if size of html > 0)
    - [x] size
    - [ ] validation
    - [ ] minification
    - [ ] feature support recommendation
    - [ ] Outdated Javascript dependencies referenced
  - [ ] check images (if size of html > 0)
    - [ ] size (< 500kb)
    - [ ] image format -> webp
    - [ ] images one-by-one, shouldn't be too large each, plus not too large in total
  - [ ] HTML accessibility check
    - [ ] Don't use fully-qualified URLs in links on the same page ("<https://example.com/b>" when you're already on "<https://example.com/a>" -> use "/b" instead as it's smaller and less error-prone).
    - [ ] Check that all links on the page also use https (or relative to the current one).
      - [ ] But due to mixed content security concerns, an HTTP `<iframe>` doesn't work in an HTTPS page.
    - [ ] amount of custom fonts & their size
    - [ ] viewport configured correctly
      - [ ] `<meta name="viewport" content="width=device-width"/>` set viewport - Die Seite enthält keine viewport-Angabe, die es Browsern ermöglicht, die Größe und Skalierung der Seite an den jeweiligen Bildschirm anzupassen. <https://web.dev/responsive-web-design-basics/#viewport> - Darstellungsbereich nicht auf „device-width“ festgelegt: Die Seite weist eine feste Breite für den viewport auf und kann sich dadurch nicht an verschiedene Bildschirmgrößen anpassen.
    - [ ] Inhalt breiter als Bildschirm: Horizontales Scrollen ist notwendig, um Text und Bilder auf der Seite sehen zu können Dies ist dann der Fall, wenn Seiten absolute Werte in CSS-Deklarationen verwenden oder Bilder nutzen, die für eine bestimmte Browserbreite optimiert sind, zum Beispiel 980 Pixel. Problem beheben: Achten Sie darauf, dass für die Seite relative Breiten- und Positionswerte für CSS-Elemente verwendet werden und Bilder ebenfalls skaliert werden können. Hier erfahren Sie, wie Sie die Größe von Inhalten an den Darstellungsbereich anpassen. <https://web.dev/responsive-web-design-basics/#size-content>
    - [ ] Text ist zu klein zum Lesen: Ein bedeutender Teil des Textes wird im Verhältnis zur Breite der Seite zu klein dargestellt. Dadurch wird der Text auf einem Mobilgerät schwer lesbar. Werfen Sie einen Blick auf den Test-Screenshot Ihres Geräts, um den problematischen Text zu finden. Problem beheben: Legen Sie einen Darstellungsbereich für Ihre Webseiten fest und wählen Sie die Schriftgrößen so aus, dass der Text innerhalb des Darstellungsbereichs entsprechend skaliert wird und auf dem Gerät sichtbar ist. Hier erfahren Sie mehr über Best Practices für die Schriftgröße. <https://web.dev/font-size/>
    - [ ] Anklickbare Elemente liegen zu dicht beieinander: Touch-Elemente wie Schaltflächen und Navigationslinks sind so dicht nebeneinander, dass mobile Nutzer beim Tippen auf ein gewünschtes Element versehentlich auch das benachbarte Element berühren. Problem beheben: Werfen Sie einen Blick auf den Testscreenshot, um alle betroffenen Schaltflächen, Links und Berührungszielbereiche zu finden. Achten Sie darauf, dass Ihre Berührungszielbereiche nicht näher beieinander liegen als die durchschnittliche Breite einer Fingerspitze, oder dass sich nicht mehr als ein Zielbereich berühren lässt. Weitere Informationen zur optimalen Größe von Tippzielen finden Sie hier. <https://web.dev/accessible-tap-targets/>
  - [ ] media embedding recommendations

- [ ] headless rendering with <https://github.com/chromedp/chromedp>

  - [ ] performance index (download speed ( first time, second time), render speed)
  - [ ] check for console errors

- [ ] Check if there's a /.well-known/security.txt file

### SEO recommendations

- [ ] check if robots.txt exists
- [ ] check if sitemap exists
- [ ] check if incompatible plugins like flash are used

## Open todos or feature ideas

### Bugfixes

- [ ] `WRN  Couldn't check ip blacklisting because of error code: ip="46.62.145.190", response="[127.255.255.254]"`
      What happened: ip blacklisting returns error code.
      Expected: try to parse the error code and print helpful text.
- [ ] Subdomain scan should print helpful text in case of common errors, like too many requests etc. instead of displaying the error code.
- [ ] List of subdomains should be filtered to show only subdomains, not all domains listed in the certificates.
- [ ] HTTP header scan results contain `Recommended action for Strict-Transport-Security: max-age value should be increased in stages from 15552000 to 63072000 (two years)`, which doens't make sense for http.
      Expected: don't print this for http.
- [ ] find solution for crt.sh error - don't show it or whatever
- [ ] `--subdomains triggers subdomainscan but not tls check, even though that's required for checking SANs of cert
- [ ] HTTP content size is 5kB even though it just redirects to https? Is there a follow-redirect set?

### Feature ideas

- find better way to pass results of previous scans to the next scan.
  Ideas:
  - target contains results, so it's not necessary to pass them as arguments to each scan.
    Advantage: Simple, easy to maintain.
    Risk: Circular dependencies.
  - results as shared package, so they can be read by all scans.
    Advantage: Strong typing, easy to maintain.
    Risk: all packages depend on this, even if it's not used.
  - use simple variables, which are then passed to the next scan.
    Advantage: Simple, works everywhere.
    Risk: complex, hard to maintain.
  - use result variables in webscan package only, and use them to pass the "simple" result variables around.
    Advantage: Strong typing, easy to maintain.
    Risk: Is this possible?
- Instead of fixed texts for recommendations and the likes, use types/enums/constants/functions for them.
  Make them comparable to each other, so f.e. http and https recommendations are comparable.
- add reverse ip lookup to subdomain scan with `resolver.LookupAddr`.
- add reverse dns lookup for ip addresses.
<!-- - [ ] List of subdomains should be sorted alphabetically. Not considered right now, as length or amount of subdomains might be more relevant for display. -->
- Search for `TODO` in code
- Check README of `pkg/status` for todos
- Check README of `pkg/logger` for todos
- Check github issues <https://github.com/thetillhoff/webscan/v3/issues>
- Ensure github actions build always has the correct version as output of `webscan version`
  - add buildargs to example usage sections in all three repos for the actions
- portscan on ipv4 and ipv6 might result in consistency-warning if your local machine only supports one of them! But sometimes it works anyway :shrug: add a note to the output after diving deeper
- HTTP header scan results and HTTPS header scan results could be merged into one if they are equal. Also, if they redirect, they should not be displayed (don't follow redirects by default).
- HTTP content scan results and HTTPS content scan results could be merged into one if they are equal. Also, if they redirect, they should not be displayed (don't follow redirects by default).
- subdomainscan should check body of all responses in cache of httpClient for referenced subdomains
- use ruleset approach for http/https and forwarding evaluation
- add repo url to help text, maybe even Issues link
- add support to read $1 arg from stdin
- add support to structure output as json with `--json-output`
- subdomainScan should list in alphabetical order, and include `*.` as dedicated host, but maybe italic or dimmed
- Check readme of thetillhoff.de for insights (accessibility, other features, plus caddyfile, ...)
- Check TTL for dns records and html caching
- urls with or without ending slash / filename & extension
  input-path and redirect locations should either end with filename.ext or `/`. `netlight.com` is "wrong" while `netlight.com/` or `netlight.com/index.html` are correct.
  The reason is that the part after the last slash might be tried to parse as filename by some applications.
  This is only a recommendation though.
- Check favicons (<https://css-tricks.com/favicons-how-to-make-sure-browsers-only-download-the-svg-version/>, <https://evilmartians.com/chronicles/how-to-favicon-in-2021-six-files-that-fit-most-needs>)
- Add `webscan status` or additional functionality to `webscan version` that checks if a new version is available (`status` could be used to check if internet connectivity is available, plus maybe scanning the local machine with it's public IP)
- RE: timeout timing: as soon as the first response came, wait for one more second, then stop waiting and continue
- httpHeaderScan results should be shown as list of items formatted like

  ```sh
  - <Problem statement like "No CSP header set">;
    <Recommendation like "To incease the security of the website, implement a CSP header.">
    More information: <link>
  ```

- Also, key-value pairs contains in CSP etc should be in separate lines each for better readability.
- tlsScan results should be shown as list of items formatted like

  ```sh
  - <Problem statement like "TLS ciphers using RC4 were prohibited to use by the IETF in 2015">;
    <Recommendation like "Remove the affected ciphers from your TLS terminator.">
    More information: <link>
    Affected ciphers:
    - <cipher name like "TLS_ECDHE_ECDSA_WITH_RC4_128_SHA">
  ```

- tlsScan could compare against known configurations like aws-tls-configs with different versions.
  Then the recommendation can be more specific ("Are you using AWS? Find out how to change the available ciphers at `<link>`")
- try to identify CMS system by checking known urls like `/wp-includes/...`
- If a request fails (timeout, etc), skip it and the resulting scan.
  Example: `netlight.com` with `https://netlight.com/wp-includes/css/dist/block-library/style.min.css\?ver\=5.3.2`
  Print warning in such a case
- think about whether, and if yes, where to add <https://ssl-config.mozilla.org/>
- add unit tests
  - add tests for redirecting output to a file or pipe it and then write it into file
    - test stdout result
    - test stderr result
      - should include logs and error messages
- Print IP RDAP info in pretty mode, depending on longest ip address
- HTTP scan should check for latency, hops, download speed
- HTML content scan should depend on content type of response; for example it should verify if it's valid json for application/json
- list all domains that are referenced (like fonts.google.com, ...)
- Outside of webscan: Add buildargs to example usage sections in all three repos for the gh-actions
  In other words: Ensure github actions build always has the correct version as output of `webscan version`
  Or: Add buildargs to example usage sections in all three repos for the actions
- check if both ipv4 and ipv6 mx records exist (follow cnames on mx records automatically)
- add functional tests with expected results on example website (github-pages?)
- add check of version in tcp greeting / header message. openssh tells the client about it's version there.
- check FTP headers
- dns checks and dials time out on windows
- add integration test to release pipeline, which runs webscan with a few sample runs on multiple platforms
- add estimated location to ips and ASN info
- check if quic uses udp on port 443, and incorporate it into scans
- add post-quantum cipher verification to tls scan
- create a new docs structure, as this one is getting too long
