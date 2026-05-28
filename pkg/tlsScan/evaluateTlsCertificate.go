package tlsScan

import (
	"crypto/tls"
	"errors"
	"log/slog"
	"os"

	"github.com/thetillhoff/webscan/v3/pkg/types"
)

// checks whether the certificate is valid
func evaluateTLSCertificate(target types.Target, ip string) ([]certInfo, error, error) {
	var (
		targetEndpoint = ip + ":" + target.Port()
		err            error

		certInfos = []certInfo{}
		tlsErr    error
	)

	slog.Debug("tlsScan: Evaluating TLS certificate started", "targetEndpoint", targetEndpoint, "ServerName", target.Hostname())

	strictConn, tlsErr := tls.Dial("tcp", targetEndpoint, &tls.Config{
		MinVersion: tls.VersionTLS10,
		MaxVersion: tls.VersionTLS13,
		ServerName: target.Hostname(),
	})
	if strictConn != nil {
		strictConn.Close()
	}

	if os.IsTimeout(tlsErr) {
		return certInfos, errors.New("http call exceeded 5s timeout"), nil
	}

	conn, err := tls.Dial("tcp", targetEndpoint, &tls.Config{
		MinVersion:         tls.VersionTLS10,
		MaxVersion:         tls.VersionTLS13,
		ServerName:         target.Hostname(),
		InsecureSkipVerify: true,
	})
	if err != nil {
		slog.Debug("tlsScan: TLS dial failed", "error", err)
		return certInfos, nil, err
	}
	defer func() {
		if closeErr := conn.Close(); closeErr != nil {
			slog.Debug("tlsScan: Error closing connection", "error", closeErr)
		}
	}()

	peerCerts := conn.ConnectionState().PeerCertificates
	for idx, cert := range peerCerts {

		certInfo := certInfo{
			issuers: []string{},
		}

		// TODO: Add expiry date to result
		// Sample: fmt.Printf("Expiry: %s \n", cert.NotAfter.Format("2006-11-02"))

		slog.Debug("tlsScan: Certificate field", "idx", idx, "common_name", cert.Subject.CommonName)
		certInfo.commonName = cert.Subject.CommonName

		slog.Debug("tlsScan: Certificate field", "idx", idx, "dns_names", cert.DNSNames)
		certInfo.sans = append(certInfo.sans, cert.DNSNames...)

		slog.Debug("tlsScan: Certificate field", "idx", idx, "email_addresses", cert.EmailAddresses)
		certInfo.sans = append(certInfo.sans, cert.EmailAddresses...)

		slog.Debug("tlsScan: Certificate field", "idx", idx, "ip_addresses", cert.IPAddresses)
		for _, ipAddress := range cert.IPAddresses {
			certInfo.sans = append(certInfo.sans, ipAddress.String())
		}

		slog.Debug("tlsScan: Certificate field", "idx", idx, "uris", cert.URIs)
		for _, uri := range cert.URIs {
			certInfo.sans = append(certInfo.sans, uri.Host)
		}

		certInfo.issuers = append(certInfo.issuers, cert.Issuer.String())
		slog.Debug("tlsScan: Certificate field", "idx", idx, "issuer", cert.Issuer.String())

		certInfos = append(certInfos, certInfo)

		slog.Debug("tlsScan: Certificate parsed", "idx", idx, "common_name", certInfo.commonName, "sans", certInfo.sans, "cert_issuers", certInfo.issuers)
	}

	slog.Debug("tlsScan: Evaluating TLS certificate completed")

	return certInfos, tlsErr, nil
}
