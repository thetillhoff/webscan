package webscan

import (
	"fmt"
	"log/slog"

	"github.com/thetillhoff/webscan/v3/pkg/dnsScan"
	"github.com/thetillhoff/webscan/v3/pkg/htmlContentScan"
	"github.com/thetillhoff/webscan/v3/pkg/httpHeaderScan"
	"github.com/thetillhoff/webscan/v3/pkg/httpProtocolScan"
	"github.com/thetillhoff/webscan/v3/pkg/ipScan"
	"github.com/thetillhoff/webscan/v3/pkg/portScan"
	"github.com/thetillhoff/webscan/v3/pkg/subDomainScan"
	"github.com/thetillhoff/webscan/v3/pkg/tlsScan"
	"github.com/thetillhoff/webscan/v3/pkg/types"
)

func (engine *Engine) Scan(input string) error {
	var (
		err error
	)

	// TODO If tty supports color, use custom logger, else use structured logger with zerolog or slog

	// TODO move the following test lines into a test case for `status` package. Ensure the final output looks the same as expected.
	// engine.status.Update("working on a")
	// time.Sleep(time.Second)
	// engine.status.Update("working on b")
	// time.Sleep(time.Second)
	// engine.status.Update("working on c")
	// time.Sleep(time.Second)
	// engine.status.Complete("complete")

	// engine.status.SpinningUpdate("working on a")
	// time.Sleep(3 * time.Second)
	// engine.status.SpinningUpdate("working on b")
	// time.Sleep(3 * time.Second)
	// engine.status.SpinningUpdate("working on c")
	// time.Sleep(3 * time.Second)
	// engine.status.SpinningComplete("complete")

	// engine.status.SpinningXOfInit(3, "things completed")
	// time.Sleep(3 * time.Second)
	// engine.status.SpinningXOfUpdate()
	// time.Sleep(3 * time.Second)
	// engine.status.SpinningXOfUpdate()
	// time.Sleep(3 * time.Second)
	// engine.status.SpinningXOfUpdate()
	// engine.status.SpinningXOfComplete()

	// Debug

	slog.Debug("webscan config",
		"followRedirects", engine.followRedirects,
		"advancedDnsScan", engine.advancedDnsScan,
		"ipScan", engine.ipScan,
		"advancedPortScan", engine.advancedPortScan,
		"tlsScan", engine.tlsScan,
		"httpProtocolScan", engine.httpProtocolScan,
		"httpHeaderScan", engine.httpHeaderScan,
		"htmlContentScan", engine.htmlContentScan,
		"mailConfigScan", engine.mailConfigScan,
		"subDomainScan", engine.subDomainScan)

	// Input

	slog.Debug("raw", "input", input)
	engine.target, err = types.NewTarget(input)
	if err != nil {
		return err
	}

	if _, err := fmt.Fprintf(engine.stdout, "# webscan results for %s\n\n", engine.target.RawTarget()); err != nil {
		slog.Debug("webscan: Error writing to output", "error", err)
	}

	// DNS

	engine.dnsScanResult, err = dnsScan.Scan(engine.target, &engine.status, dnsScan.WithAdvanced(engine.advancedDnsScan), dnsScan.WithFollowRedirects(engine.followRedirects))
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

		tlsScan.PrintResult(engine.tlsScanResult, engine.stdout)
	}

	// HTTP protocol scan

	if engine.httpProtocolScan {
		engine.httpProtocolScanResult, err = httpProtocolScan.Scan(
			engine.target,
			&engine.status,
			httpProtocolScan.WithClient(engine.client),
			httpProtocolScan.WithIsAvailableViaPort80(engine.portScanResult.IsPortOpen(80)),
			httpProtocolScan.WithIsAvailableViaPort443(engine.portScanResult.IsPortOpen(443)),
		)
		if err != nil {
			return err
		}
	}

	httpProtocolScan.PrintResult(engine.httpProtocolScanResult, engine.stdout)

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

		if engine.portScanResult.IsPortOpen(80) && engine.httpProtocolScanResult.IsAvailableViaHttp() {
			engine.httpHtmlContentScanResult, err = htmlContentScan.Scan(&engine.status, engine.target, htmlContentScan.WithClient(engine.client), htmlContentScan.WithSchemaOverride(types.HTTP))
			if err != nil {
				return err
			}

			htmlContentScan.PrintResult(engine.httpHtmlContentScanResult, "http", engine.stdout)
		}

		if engine.portScanResult.IsPortOpen(443) && engine.httpProtocolScanResult.IsAvailableViaHttps() {
			engine.httpsHtmlContentScanResult, err = htmlContentScan.Scan(&engine.status, engine.target, htmlContentScan.WithClient(engine.client), htmlContentScan.WithSchemaOverride(types.HTTPS))
			if err != nil {
				return err
			}

			htmlContentScan.PrintResult(engine.httpsHtmlContentScanResult, "https", engine.stdout)
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
		)

		subDomainScan.PrintResult(engine.subDomainScanResult, engine.stdout)
	}

	// TODO if followRedirects is true, CNAMEs should be followed (scan them, too)
	// TODO if followRedirects is true, http and https redirects should be followed (scan them, too)

	return err
}
