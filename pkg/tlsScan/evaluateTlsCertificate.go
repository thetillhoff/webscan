package tlsScan

import (
	"crypto/tls"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"time"
)

// checks whether the certificate is valid
func evaluateTlsCertificate(url string) (map[string]struct{}, error) {
	var (
		err error

		uniqueCertNames = map[string]struct{}{}
	)

	slog.Debug("tlsScan: Evaluating TLS certificate started")

	client := &http.Client{
		Timeout: 5 * time.Second, // TODO 5s might be a bit long?
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		}, // Don't follow redirects
	}

	_, err = client.Get("https://" + url)

	if err != nil {
		if os.IsTimeout(err) { // If err is timeout, tell user about it
			return uniqueCertNames, errors.New("http call exceeded 5s timeout")
		} else {
			return uniqueCertNames, errors.New("Invalid TLS certificate used for HTTPS: " + errors.Unwrap(errors.Unwrap(err)).Error())
		}
	}

	//

	conn, err := tls.Dial("tcp", url+":443", &tls.Config{
		InsecureSkipVerify: true,
	})
	if err != nil {
		slog.Error("could not load remote certificate", "error", err)
	}
	defer conn.Close()
	certs := conn.ConnectionState().PeerCertificates
	for _, cert := range certs {
		// fmt.Printf("Issuer Name: %s\n", cert.Issuer)
		// fmt.Printf("Expiry: %s \n", cert.NotAfter.Format("2006-11-02"))

		uniqueCertNames[cert.Subject.String()] = struct{}{}
		// fmt.Println("subject:", cert.Subject)

		for _, dnsName := range cert.DNSNames {
			uniqueCertNames[dnsName] = struct{}{}
		}
		// fmt.Println("dnsnames:", cert.DNSNames)

		for _, emailAddress := range cert.EmailAddresses {
			uniqueCertNames[emailAddress] = struct{}{}
		}
		// fmt.Println("mail addresses:", cert.EmailAddresses)

		for _, ipAddress := range cert.IPAddresses {
			uniqueCertNames[ipAddress.String()] = struct{}{}
		}
		// fmt.Println("ip addresses:", cert.IPAddresses)

		for _, uri := range cert.URIs {
			uniqueCertNames[uri.Host] = struct{}{}
		}
		// fmt.Println("uris:", cert.URIs)

	}

	slog.Debug("tlsScan: Result", "cert names", uniqueCertNames)

	slog.Debug("tlsScan: Evaluating TLS certificate completed")

	return uniqueCertNames, nil
}
