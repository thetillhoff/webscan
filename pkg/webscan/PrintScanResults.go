package webscan

import "net"

func (engine Engine) PrintScanResults(inputUrl string) {

	netIP := net.ParseIP(inputUrl)
	if netIP == nil { // If input was an IPaddress, nothing to see here
		engine.PrintDnsScanResults(inputUrl)

		engine.PrintIpScanResults()
	}

	engine.PrintPortScanResults()

	engine.PrintTlsScanResults()

	engine.PrintProtocolScanResults(inputUrl)

	engine.PrintHttpHeaderScanResults()

	engine.PrintHttpContentScanResults()

	engine.PrintMailConfigScanResults()

	engine.PrintSubdomainScanResults()

}
