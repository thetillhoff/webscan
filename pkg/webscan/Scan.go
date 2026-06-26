package webscan

import (
	"fmt"
	"log/slog"

	"github.com/thetillhoff/webscan/v5/pkg/dnsScan"
	"github.com/thetillhoff/webscan/v5/pkg/htmlContentScan"
	"github.com/thetillhoff/webscan/v5/pkg/httpHeaderScan"
	"github.com/thetillhoff/webscan/v5/pkg/httpProtocolScan"
	"github.com/thetillhoff/webscan/v5/pkg/ipScan"
	"github.com/thetillhoff/webscan/v5/pkg/knownFilesScan"
	"github.com/thetillhoff/webscan/v5/pkg/portScan"
	"github.com/thetillhoff/webscan/v5/pkg/subDomainScan"
	"github.com/thetillhoff/webscan/v5/pkg/tlsScan"
	"github.com/thetillhoff/webscan/v5/pkg/types"
)

func (engine *Engine) Scan(input string) error {
	var (
		err error
	)

	// TODO If tty supports color, use custom logger, else use structured logger with zerolog or slog

	slog.Debug("webscan config",
		"followRedirects", engine.followRedirects,
		"advancedDnsScan", engine.advancedDnsScan,
		"ipScan", engine.ipScan,
		"advancedPortScan", engine.advancedPortScan,
		"tlsScan", engine.tlsScan,
		"httpProtocolScan", engine.httpProtocolScan,
		"httpHeaderScan", engine.httpHeaderScan,
		"htmlContentScan", engine.htmlContentScan,
		"knownFilesScan", engine.knownFilesScan,
		"mailConfigScan", engine.mailConfigScan,
		"subDomainScan", engine.subDomainScan)

	// Input

	slog.Debug("webscan: Raw input", "input", input)
	engine.target, err = types.NewTarget(input)
	if err != nil {
		return err
	}

	if _, err := fmt.Fprintf(engine.stdout, "# webscan results for %s\n\n", engine.target.RawTarget()); err != nil {
		slog.Debug("webscan: Error writing to output", "error", err)
	}

	// DNS

	engine.dnsScanResult, err = dnsScan.Scan(
		engine.target,
		&engine.status,
		dnsScan.WithCustomNameServer(engine.dnsServer),
		dnsScan.WithAdvanced(engine.advancedDnsScan),
		dnsScan.WithFollowRedirects(engine.followRedirects),
		dnsScan.WithTimeout(engine.timeout),
	)
	if err != nil {
		return err
	}

	dnsScan.PrintResult(engine.dnsScanResult, engine.stdout)

	// IP scan

	if engine.ipScan {

		engine.ipScanResult, err = ipScan.Scan(
			engine.target,
			&engine.status,
			ipScan.WithARecords(engine.dnsScanResult.ARecords),
			ipScan.WithAAAARecords(engine.dnsScanResult.AAAARecords),
			ipScan.WithTimeout(engine.timeout),
		)
		if err != nil {
			return err
		}

		ipScan.PrintResult(
			engine.ipScanResult,
			engine.dnsScanResult.ARecords,
			engine.dnsScanResult.AAAARecords,
			engine.stdout,
		)
	}

	// Port scan

	engine.portScanResult, err = portScan.Scan(
		engine.target,
		&engine.status,
		portScan.WithARecords(engine.dnsScanResult.ARecords),
		portScan.WithAAAARecords(engine.dnsScanResult.AAAARecords),
		portScan.WithAdvanced(engine.advancedPortScan),
		portScan.WithTimeout(engine.timeout),
	)
	if err != nil {
		return err
	}

	portScan.PrintResult(engine.portScanResult, engine.stdout)

	// TLS scan

	// TODO only run tls scan if protocol is tls, https or not specified.
	// In cast of tls or https, run it either on 443 or another port if one is specified.

	if len(engine.dnsScanResult.ARecords) > 0 || len(engine.dnsScanResult.AAAARecords) > 0 {

		if engine.tlsScan || engine.subDomainScan {

			engine.tlsScanResult, err = tlsScan.Scan(
				engine.target,
				&engine.status,
				engine.dnsScanResult.ARecords,
				engine.dnsScanResult.AAAARecords,
				engine.portScanResult,
			)
			if err != nil {
				return err
			}
		}

		if engine.tlsScan {
			tlsScan.PrintResult(engine.tlsScanResult, engine.stdout)
		}
	}

	// HTTP protocol scan (required by header, content, and known files scans)

	if engine.httpProtocolScan || engine.httpHeaderScan || engine.htmlContentScan || engine.knownFilesScan {
		engine.httpProtocolScanResult, err = httpProtocolScan.Scan(
			engine.target,
			&engine.status,
			httpProtocolScan.WithIsAvailableViaPort80(engine.portScanResult.IsPortOpen(80)),
			httpProtocolScan.WithIsAvailableViaPort443(engine.portScanResult.IsPortOpen(443)),
			httpProtocolScan.WithTimeout(engine.timeout),
		)
		if err != nil {
			return err
		}
	}

	if engine.httpProtocolScan {
		httpProtocolScan.PrintResult(engine.httpProtocolScanResult, engine.stdout)
	}

	// HTTP header scan

	if engine.httpHeaderScan {

		if engine.portScanResult.IsPortOpen(80) && engine.httpProtocolScanResult.IsAvailableViaHttp() {
			engine.httpHeaderScanResult, err = httpHeaderScan.Scan(&engine.status, engine.target, httpHeaderScan.WithClient(engine.client), httpHeaderScan.WithSchemaOverride(types.HTTP))
			if err != nil {
				return err
			}

			httpHeaderScan.PrintResult(engine.httpHeaderScanResult, "http", engine.stdout)
		}

		if engine.portScanResult.IsPortOpen(443) && engine.httpProtocolScanResult.IsAvailableViaHttps() {
			engine.httpsHeaderScanResult, err = httpHeaderScan.Scan(&engine.status, engine.target, httpHeaderScan.WithClient(engine.client), httpHeaderScan.WithSchemaOverride(types.HTTPS))
			if err != nil {
				return err
			}

			httpHeaderScan.PrintResult(engine.httpsHeaderScanResult, "https", engine.stdout)
		}
	}

	if engine.htmlContentScan {
		httpAvail := engine.portScanResult.IsPortOpen(80) && engine.httpProtocolScanResult.IsAvailableViaHttp()
		httpsAvail := engine.portScanResult.IsPortOpen(443) && engine.httpProtocolScanResult.IsAvailableViaHttps()

		if httpAvail {
			engine.httpHtmlContentScanResult, err = htmlContentScan.Scan(&engine.status, engine.target, htmlContentScan.WithClient(engine.client), htmlContentScan.WithSchemaOverride(types.HTTP))
			if err != nil {
				return err
			}
		}

		if httpsAvail {
			engine.httpsHtmlContentScanResult, err = htmlContentScan.Scan(&engine.status, engine.target, htmlContentScan.WithClient(engine.client), htmlContentScan.WithSchemaOverride(types.HTTPS))
			if err != nil {
				return err
			}
		}

		if httpAvail && httpsAvail && engine.httpHtmlContentScanResult.Equal(engine.httpsHtmlContentScanResult) {
			htmlContentScan.PrintResult(engine.httpHtmlContentScanResult, "HTTP & HTTPS", engine.stdout)
		} else {
			if httpAvail {
				htmlContentScan.PrintResult(engine.httpHtmlContentScanResult, "http", engine.stdout)
			}
			if httpsAvail {
				htmlContentScan.PrintResult(engine.httpsHtmlContentScanResult, "https", engine.stdout)
			}
		}
	}

	// Known files scan

	if engine.knownFilesScan {
		httpAvail := engine.portScanResult.IsPortOpen(80) && engine.httpProtocolScanResult.IsAvailableViaHttp()
		httpsAvail := engine.portScanResult.IsPortOpen(443) && engine.httpProtocolScanResult.IsAvailableViaHttps()

		if httpAvail {
			engine.httpKnownFilesScanResult = knownFilesScan.Scan(
				engine.target,
				&engine.status,
				types.HTTP,
				knownFilesScan.WithTimeout(engine.timeout),
			)
		}

		if httpsAvail {
			engine.httpsKnownFilesScanResult = knownFilesScan.Scan(
				engine.target,
				&engine.status,
				types.HTTPS,
				knownFilesScan.WithTimeout(engine.timeout),
			)
		}

		if httpAvail && httpsAvail && engine.httpKnownFilesScanResult.EqualContent(engine.httpsKnownFilesScanResult) {
			knownFilesScan.PrintResult(engine.httpKnownFilesScanResult, "HTTP & HTTPS", engine.stdout)
		} else {
			if httpAvail {
				knownFilesScan.PrintResult(engine.httpKnownFilesScanResult, "http", engine.stdout)
			}
			if httpsAvail {
				knownFilesScan.PrintResult(engine.httpsKnownFilesScanResult, "https", engine.stdout)
			}
		}
	}

	// if engine.MailConfigScan {
	// 	engine, err = engine.ScanMailConfig(input)
	// 	if err != nil {
	// 		return engine, err
	// 	}
	// }

	if engine.subDomainScan {
		engine.subDomainScanResult = subDomainScan.Scan(
			engine.target,
			&engine.status,
			subDomainScan.WithCertNames(engine.tlsScanResult.ListAllCertNames()),
			subDomainScan.WithTimeout(engine.timeout),
		)

		subDomainScan.PrintResult(engine.subDomainScanResult, engine.stdout)
	}

	// Follow HTTP redirects: run web scans (protocol, header, content) on redirect targets
	if engine.followRedirects {
		type redirectJob struct {
			loc    string
			reason string
		}

		visited := map[string]bool{
			engine.target.UrlString(): true,
		}

		var redirectTargets []redirectJob
		if loc := engine.httpProtocolScanResult.HttpRedirectLocation(); loc != "" {
			redirectTargets = append(redirectTargets, redirectJob{loc, fmt.Sprintf("HTTP %d", engine.httpProtocolScanResult.HttpStatusCode())})
		}
		if loc := engine.httpProtocolScanResult.HttpsRedirectLocation(); loc != "" {
			redirectTargets = append(redirectTargets, redirectJob{loc, fmt.Sprintf("HTTPS %d", engine.httpProtocolScanResult.HttpsStatusCode())})
		}

		for len(redirectTargets) > 0 {
			job := redirectTargets[0]
			redirectTargets = redirectTargets[1:]
			loc := job.loc

			if visited[loc] {
				continue
			}
			visited[loc] = true

			redirectTarget, parseErr := types.NewTarget(loc)
			if parseErr != nil {
				slog.Debug("webscan: could not parse redirect target", "location", loc, "error", parseErr)
				continue
			}

			if _, err := fmt.Fprintf(engine.stdout, "\n---\n\n# webscan results for %s (%s redirect)\n\n", loc, job.reason); err != nil {
				slog.Debug("webscan: Error writing to output", "error", err)
			}

			// Protocol scan on redirect target
			schema := redirectTarget.Schema()
			if schema == types.NONE {
				schema = types.HTTPS
			}

			redirectProtocolResult, protocolErr := httpProtocolScan.Scan(
				redirectTarget,
				&engine.status,
				httpProtocolScan.WithIsAvailableViaPort80(schema == types.HTTP),
				httpProtocolScan.WithIsAvailableViaPort443(schema == types.HTTPS),
				httpProtocolScan.WithTimeout(engine.timeout),
			)
			if protocolErr != nil {
				slog.Debug("webscan: redirect protocol scan failed", "target", loc, "error", protocolErr)
				continue
			}

			httpProtocolScan.PrintResult(redirectProtocolResult, engine.stdout)

			// Header scan on redirect target
			if engine.httpHeaderScan {
				headerResult, headerErr := httpHeaderScan.Scan(&engine.status, redirectTarget, httpHeaderScan.WithClient(engine.client), httpHeaderScan.WithSchemaOverride(schema))
				if headerErr == nil {
					httpHeaderScan.PrintResult(headerResult, schema.String(), engine.stdout)
				}
			}

			// Content scan on redirect target
			if engine.htmlContentScan {
				contentResult, contentErr := htmlContentScan.Scan(&engine.status, redirectTarget, htmlContentScan.WithClient(engine.client), htmlContentScan.WithSchemaOverride(schema))
				if contentErr == nil {
					htmlContentScan.PrintResult(contentResult, schema.String(), engine.stdout)
				}
			}

			// Known files scan on redirect target
			if engine.knownFilesScan {
				filesResult := knownFilesScan.Scan(
					redirectTarget,
					&engine.status,
					schema,
					knownFilesScan.WithTimeout(engine.timeout),
				)
				knownFilesScan.PrintResult(filesResult, schema.String(), engine.stdout)
			}

			// Queue further redirects from this target
			if newLoc := redirectProtocolResult.HttpRedirectLocation(); newLoc != "" && !visited[newLoc] {
				redirectTargets = append(redirectTargets, redirectJob{newLoc, fmt.Sprintf("HTTP %d", redirectProtocolResult.HttpStatusCode())})
			}
			if newLoc := redirectProtocolResult.HttpsRedirectLocation(); newLoc != "" && !visited[newLoc] {
				redirectTargets = append(redirectTargets, redirectJob{newLoc, fmt.Sprintf("HTTPS %d", redirectProtocolResult.HttpsStatusCode())})
			}
		}
	}

	return nil
}
