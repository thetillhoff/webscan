package portScan

import (
	"log/slog"

	"github.com/thetillhoff/webscan/pkg/status"
)

func AdvancedScan(status *status.Status, aRecords []string, aaaaRecords []string) (Result, error) {
	var (
		scanPorts = []uint16{
			15,   // Netstat
			20,   // FTP
			21,   // FTP
			22,   // SSH
			25,   // SMTP
			53,   // DNS
			67,   // BOOTP
			68,   // BOOTP
			69,   // TFTP
			80,   // HTTP
			88,   // Kerberos
			110,  // POP3
			111,  // PortMap
			119,  // NNTP
			123,  // NTP
			137,  // NetBIOS (SMB related)
			138,  // NetBIOS (SMB related)
			139,  // NetBIOS
			143,  // IMAP
			179,  // BGP
			389,  // LDAP
			443,  // HTTPS
			445,  // SMB
			465,  // SMTP SSL
			587,  // SMTP TLS
			636,  // SLDAP
			1433, // Microsoft SQL
			2525, // SMTP TLS alternative
			3306, // MySQL/MariaDB
			3389, // RDP
			5432, // Postgres
			8080, // HTTP alternative
			8443, // HTTP alternative
		}

		result = Result{
			openPorts:               []uint16{},
			openPortInconsistencies: []string{},
			isAvailableViaHttp:      false,
			isAvailableViaHttps:     false,
		}

		openPortsPerIp map[string][]uint16
	)

	slog.Debug("portScan: Advanced scan started")

	openPortsPerIp = scanPortRangeOfIps(status, append(aRecords, aaaaRecords...), scanPorts)

	result.openPorts, result.openPortInconsistencies = CompareOpenPortsOfIps(openPortsPerIp)

	// Check if HTTP / HTTPS are available
	for _, openPort := range result.openPorts {
		if openPort == 80 {
			result.isAvailableViaHttp = true
		} else if openPort == 443 {
			result.isAvailableViaHttps = true
		}
	}

	slog.Debug("portScan: Advanced scan completed")

	return result, nil
}
