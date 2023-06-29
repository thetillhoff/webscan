package portScan

func ScanPortRangeOfIps(ips []string, ports []uint16) map[string][]uint16 {
	var (
		openPortsPerIp = map[string][]uint16{}
	)

	// TODO parallelize this

	for _, ip := range ips {
		openPortsPerIp[ip] = ScanPortRangeOfIp(ip, ports)
	}

	return openPortsPerIp
}
