package webscan

import (
	"fmt"
	"log/slog"

	"github.com/thetillhoff/webscan/pkg/dnsScan"
	"github.com/thetillhoff/webscan/pkg/htmlContentScan"
	"github.com/thetillhoff/webscan/pkg/httpHeaderScan"
	"github.com/thetillhoff/webscan/pkg/httpProtocolScan"
	"github.com/thetillhoff/webscan/pkg/ipScan"
	"github.com/thetillhoff/webscan/pkg/portScan"
	"github.com/thetillhoff/webscan/pkg/subDomainScan"
	"github.com/thetillhoff/webscan/pkg/tlsScan"
)

func (engine *Engine2) PrintResults() {

	if engine.instant { // If instant-output is enabled
		slog.Debug("engine.instant is enabled, so printing of results at the end is skipped")
		return // Don't print any results here
	}

	fmt.Printf("# webscan results for %s\n", engine.target.rawTarget)

	if engine.target.isDomain && engine.advancedDnsScan {
		dnsScan.PrintResult(engine.dnsScanResult)
	}

	if engine.ipScan {
		ipScan.PrintResult(engine.ipScanResult, engine.dnsScanResult.ARecords, engine.dnsScanResult.AAAARecords)
	}

	if engine.advancedPortScan {
		portScan.PrintResult(engine.portScanResult)
	}

	if engine.portScanResult.IsAvailableViaHttps() && engine.tlsScan {
		tlsScan.PrintResult(engine.tlsScanResult)
	}

	if engine.portScanResult.IsAvailableViaHttp() || engine.portScanResult.IsAvailableViaHttps() {

		if engine.httpProtocolScan {
			httpProtocolScan.PrintResult(engine.httpProtocolScanResult)
		}

		if engine.httpHeaderScan {
			httpHeaderScan.PrintResult(engine.httpHeaderScanResult, "http")
			httpHeaderScan.PrintResult(engine.httpsHeaderScanResult, "https")
		}

		if engine.htmlContentScan {
			htmlContentScan.PrintResult(engine.httpHtmlContentScanResult, "http")
			htmlContentScan.PrintResult(engine.httpsHtmlContentScanResult, "https")
		}
	}

	// mailConfigScan.PrintResult(engine.mailConfigScanResult) // TODO

	if engine.subDomainScan {
		subDomainScan.PrintResult(engine.subDomainScanResult)
	}

}
