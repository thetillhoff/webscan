package portScan

import "slices"

func (result *Result) IsPortOpen(port uint16) bool {
	return slices.Contains(result.openPorts, port)
}

func (result *Result) IsPortOpenOnIp(ip string, port uint16) bool {
	return slices.Contains(result.OpenPortsPerIp[ip], port)
}
