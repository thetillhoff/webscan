package webscan

func (engine Engine) PrintScanResults() {

	engine.PrintDnsScanEngines()

	engine.PrintIpScanResults()

	engine.PrintPortScanResults()

	engine.PrintTlsScanResults()

	engine.PrintHttpProtocolScanResults()

	engine.PrintHttpHeaderScanResults()

	engine.PrintHttpContentScanResults()

	engine.PrintMailConfigScanResults()

	engine.PrintSubdomainScanResults()

}
