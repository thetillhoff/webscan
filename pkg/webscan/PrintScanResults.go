package webscan

func (engine Engine) PrintScanResults() {

	// check if dns scan results are empty, and if not print them
	// check if ip scan results are empty, and if not print them
	// then remove the part below

	if engine.inputType == Domain {
		engine.PrintDnsScanEngines()
	}

	engine.PrintIpScanResults()

	engine.PrintPortScanResults()

	engine.PrintTlsScanResults()

	engine.PrintProtocolScanResults()

	engine.PrintHttpHeaderScanResults()

	engine.PrintHttpContentScanResults()

	engine.PrintMailConfigScanResults()

	engine.PrintSubdomainScanResults()

}
