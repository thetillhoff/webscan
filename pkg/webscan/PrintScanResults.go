package webscan

func (engine Engine) PrintScanResults() {

	if len(engine.DnsScanEngine.ARecords) == 0 && len(engine.DnsScanEngine.AAAARecords) == 0 { // If input was an IPaddress, don't even try...
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
