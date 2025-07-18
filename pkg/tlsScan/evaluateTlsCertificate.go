package tlsScan

import (
	"crypto/tls"
	"errors"
	"log/slog"
	"os"

	"github.com/thetillhoff/webscan/v3/pkg/types"
)

// checks whether the certificate is valid
func evaluateTlsCertificate(target types.Target, ip string) ([]certInfo, error, error) {
	var (
		targetEndpoint = ip + ":" + target.Port()
		err            error

		certInfos = []certInfo{}
		tlsErr    error
	)

	slog.Debug("tlsScan: Evaluating TLS certificate started", "targetEndpoint", targetEndpoint, "ServerName", target.Hostname())

	_, tlsErr = tls.Dial("tcp", targetEndpoint, &tls.Config{
		MinVersion: tls.VersionTLS10,
		MaxVersion: tls.VersionTLS13,
		ServerName: target.Hostname(),
	})

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
			names:   []string{},
			issuers: []string{},
		}

		// TODO: Add expiry date to result
		// Sample: fmt.Printf("Expiry: %s \n", cert.NotAfter.Format("2006-11-02"))

		slog.Debug("tlsScan", "idx", idx, "common name", cert.Subject.CommonName)
		certInfo.names = append(certInfo.names, cert.Subject.CommonName)

		slog.Debug("tlsScan", "idx", idx, "dns names", cert.DNSNames)
		certInfo.names = append(certInfo.names, cert.DNSNames...)

		slog.Debug("tlsScan", "idx", idx, "email addresses", cert.EmailAddresses)
		certInfo.names = append(certInfo.names, cert.EmailAddresses...)

		slog.Debug("tlsScan", "idx", idx, "ip addresses", cert.IPAddresses)
		for _, ipAddress := range cert.IPAddresses {
			certInfo.names = append(certInfo.names, ipAddress.String())
		}

		slog.Debug("tlsScan", "idx", idx, "uris", cert.URIs)
		for _, uri := range cert.URIs {
			certInfo.names = append(certInfo.names, uri.Host)
		}

		certInfo.issuers = append(certInfo.issuers, cert.Issuer.String())
		slog.Debug("tlsScan", "idx", idx, "issuer", cert.Issuer.String())

		certInfos = append(certInfos, certInfo)

		slog.Debug("tlsScan: certInfo", "idx", idx, "cert names", certInfo.names, "cert issuers", certInfo.issuers)
	}

	slog.Debug("tlsScan: Evaluating TLS certificate completed")

	return certInfos, tlsErr, nil
}
