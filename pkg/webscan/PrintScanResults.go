package webscan

import "net"

func (engine Engine) PrintScanResults() {

	netIP := net.ParseIP(engine.url)
	if netIP == nil { // If input was an IPaddress, nothing to see here
		engine.PrintDnsScanResults()

		engine.PrintIpScanResults()
	}

	engine.PrintPortScanResults()

	engine.PrintTlsScanResults()

	engine.PrintProtocolScanResults()

	engine.PrintHttpHeaderScanResults()

	engine.PrintHttpContentScanResults()

	engine.PrintMailConfigScanResults()

	engine.PrintSubdomainScanResults()

}
