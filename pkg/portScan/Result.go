package portScan

type Result struct {
	OpenPortsPerIp map[string][]uint16

	openPorts               []uint16
	openPortInconsistencies []string

	isPort80Open  bool
	isPort443Open bool
}
