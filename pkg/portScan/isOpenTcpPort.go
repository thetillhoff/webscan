package portScan

import (
	"log/slog"
	"net"
	"strconv"
	"time"

	"github.com/thetillhoff/webscan/v3/pkg/status"
)

// Checks whether a tcp connection can be established to the specified <ip>:<port> combo.
// Has a timeout of 2s
func isOpenTcpPort(status *status.Status, ipPortTuple IpPortTuple, portChannel chan IpPortTuple) {
	var (
		err error
	)

	defer status.SpinningXOfUpdate()
	defer wgPortScan.Done()

	slog.Debug("portScan: Checking tcp port started", "ip", ipPortTuple.Ip, "port", ipPortTuple.Port)

	_, err = net.DialTimeout("tcp", ipPortTuple.Ip+":"+strconv.FormatUint(uint64(ipPortTuple.Port), 10), 5*time.Second)

	if err == nil {

		portChannel <- ipPortTuple // Write ip & port combo to channel only if it was open

		slog.Debug("portScan: Checking tcp port completed", "ip", ipPortTuple.Ip, "port", ipPortTuple.Port, "isOpen", true)

	} else {

		slog.Debug("portScan: Checking tcp port completed", "ip", ipPortTuple.Ip, "port", ipPortTuple.Port, "isOpen", false)

	}
}
