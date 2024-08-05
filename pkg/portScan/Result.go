package portScan

type Result struct {
	openPorts []uint16

	openPortInconsistencies []string

	isAvailableViaHttp  bool
	isAvailableViaHttps bool
}
