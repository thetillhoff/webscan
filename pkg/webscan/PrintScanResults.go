package webscan

func (engine Engine) PrintScanResults() {

	engine.PrintDnsScanResults()

	engine.PrintIpScanResults()

	engine.PrintPortScanResults()

	engine.PrintTlsScanResults()

	engine.PrintProtocolScanResults()

	engine.PrintHttpHeaderScanResults()

	engine.PrintHttpContentScanResults()

	engine.PrintMailConfigScanResults()

	engine.PrintSubdomainScanResults()

}
