package portScan

import (
	"log/slog"
	"slices"

	"github.com/thetillhoff/webscan/v3/pkg/status"
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
			6443, // HTTPS alternative, Kubernetes Control Plane
			8080, // HTTP alternative
			8443, // HTTPS alternative
		}

		result = Result{
			OpenPortsPerIp:          map[string][]uint16{},
			openPorts:               []uint16{},
			openPortInconsistencies: []string{},
			isPort80Open:            false,
			isPort443Open:           false,
		}
	)

	slog.Debug("portScan: Advanced scan started")

	result.OpenPortsPerIp = scanPortRangeOfIps(status, append(aRecords, aaaaRecords...), scanPorts)

	result.openPorts, result.openPortInconsistencies = CompareOpenPortsOfIps(result.OpenPortsPerIp)

	// Check if port 80 / 443 are open
	result.isPort80Open = slices.Contains(result.openPorts, 80)
	result.isPort443Open = slices.Contains(result.openPorts, 443)

	slog.Debug("portScan: Advanced scan completed")

	return result, nil
}
