package webscan

import (
	"context"
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

// emitDualSchema runs a web scan over the HTTP and/or HTTPS schema (whichever is
// available), then prints the results. When equal != nil and both schemas are
// available and produce equal results, a single combined result is printed;
// otherwise each available schema is printed separately.
func emitDualSchema[R any](
	httpAvail, httpsAvail bool,
	scan func(schema types.Schema) (R, error),
	equal func(a, b R) bool,
	print func(r R, label string),
) error {
	var httpRes, httpsRes R

	if httpAvail {
		r, err := scan(types.HTTP)
		if err != nil {
			return err
		}
		httpRes = r
	}

	if httpsAvail {
		r, err := scan(types.HTTPS)
		if err != nil {
			return err
		}
		httpsRes = r
	}

	if equal != nil && httpAvail && httpsAvail && equal(httpRes, httpsRes) {
		print(httpRes, "HTTP & HTTPS")
		return nil
	}

	if httpAvail {
		print(httpRes, "http")
	}
	if httpsAvail {
		print(httpsRes, "https")
	}

	return nil
}

// scanWebTarget runs the enabled per-target web scans (header, content, known
// files) for the given available schemas and prints their results. Used both
// for the primary target and for redirect targets.
func (engine *Engine) scanWebTarget(target types.Target, httpAvail, httpsAvail bool) error {
	if engine.httpHeaderScan {
		if err := emitDualSchema(httpAvail, httpsAvail,
			func(s types.Schema) (httpHeaderScan.Result, error) {
				return httpHeaderScan.Scan(&engine.status, target, httpHeaderScan.WithClient(engine.client), httpHeaderScan.WithSchemaOverride(s))
			},
			nil, // header results are never merged
			func(r httpHeaderScan.Result, label string) { httpHeaderScan.PrintResult(r, label, engine.stdout) },
		); err != nil {
			return err
		}
	}

	if engine.htmlContentScan {
		if err := emitDualSchema(httpAvail, httpsAvail,
			func(s types.Schema) (htmlContentScan.Result, error) {
				return htmlContentScan.Scan(&engine.status, target, htmlContentScan.WithClient(engine.client), htmlContentScan.WithSchemaOverride(s))
			},
			func(a, b htmlContentScan.Result) bool { return a.Equal(b) },
			func(r htmlContentScan.Result, label string) { htmlContentScan.PrintResult(r, label, engine.stdout) },
		); err != nil {
			return err
		}
	}

	if engine.knownFilesScan {
		if err := emitDualSchema(httpAvail, httpsAvail,
			func(s types.Schema) (knownFilesScan.Result, error) {
				return knownFilesScan.Scan(target, &engine.status, s, knownFilesScan.WithTimeout(engine.timeout)), nil
			},
			func(a, b knownFilesScan.Result) bool { return a.EqualContent(b) },
			func(r knownFilesScan.Result, label string) { knownFilesScan.PrintResult(r, label, engine.stdout) },
		); err != nil {
			return err
		}
	}

	return nil
}

func (engine *Engine) Scan(ctx context.Context, input string) error {
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

	// Derived scan requirements: some scans must run to feed others even when the
	// user didn't request their output. Naming them once documents the hidden
	// dependencies (TLS feeds subdomain cert names; protocol feeds header/content/
	// known-files availability).
	needTLS := engine.tlsScan || engine.subDomainScan
	needProtocol := engine.httpProtocolScan || engine.httpHeaderScan || engine.htmlContentScan || engine.knownFilesScan

	// Input

	slog.Debug("webscan: Raw input", "input", input)
	engine.target, err = types.NewTarget(input)
	if err != nil {
		return err
	}

	if _, err := fmt.Fprintf(engine.stdout, "# webscan results\n\n"); err != nil {
		slog.Debug("webscan: Error writing to output", "error", err)
	}

	// DNS

	if err := ctx.Err(); err != nil {
		return err
	}

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

	// Without A/AAAA records there is nothing to connect to, so all
	// IP-dependent phases (port, TLS, HTTP) are skipped.
	hasIPs := len(engine.dnsScanResult.ARecords) > 0 || len(engine.dnsScanResult.AAAARecords) > 0

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

	if err := ctx.Err(); err != nil {
		return err
	}

	if hasIPs {
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
	}

	// TLS scan

	// TODO only run tls scan if protocol is tls, https or not specified.
	// In cast of tls or https, run it either on 443 or another port if one is specified.

	if hasIPs && needTLS {
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

		if engine.tlsScan {
			tlsScan.PrintResult(engine.tlsScanResult, engine.stdout)
		}
	}

	// HTTP protocol scan (required by header, content, and known files scans)

	if err := ctx.Err(); err != nil {
		return err
	}

	if hasIPs && needProtocol {
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

	// HTTP web scans (header, content, known files) on the primary target.
	httpAvail := hasIPs && engine.portScanResult.IsPortOpen(80) && engine.httpProtocolScanResult.IsAvailableViaHttp()
	httpsAvail := hasIPs && engine.portScanResult.IsPortOpen(443) && engine.httpProtocolScanResult.IsAvailableViaHttps()

	if err := engine.scanWebTarget(engine.target, httpAvail, httpsAvail); err != nil {
		return err
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
			if err := ctx.Err(); err != nil {
				return err
			}

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

			// Web scans on redirect target (only the redirect's own schema is available).
			if err := engine.scanWebTarget(redirectTarget, schema == types.HTTP, schema == types.HTTPS); err != nil {
				slog.Debug("webscan: redirect web scan failed", "target", loc, "error", err)
				continue
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
