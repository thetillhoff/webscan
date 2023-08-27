package portScan

import (
	"net"
	"strconv"
	"time"
)

// Checks whether a tcp connection can be established to the specified <ip>:<port> combo.
// Has a timeout of 2s
func isOpenTcpPort(ip net.IP, port uint16, portChannel chan<- uint16) {
	var (
		err error
	)
	defer wgPortScan.Done()

	_, err = net.DialTimeout("tcp", ip.String()+":"+strconv.FormatUint(uint64(port), 10), 5*time.Second)
	if err == nil {
		portChannel <- port
	} // Write ip & port combo to channel if it was open
}
