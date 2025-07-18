package subDomainScan

import (
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/thetillhoff/webscan/v3/pkg/types"
)

type Cert struct {
	IssuerCaID     int    `json:"issuer_ca_id"`
	IssuerName     string `json:"issuer_name"`
	CommonName     string `json:"common_name"`
	NameValue      string `json:"name_value"`
	ID             int64  `json:"id"`
	EntryTimestamp string `json:"entry_timestamp"`
	NotBefore      string `json:"not_before"`
	NotAfter       string `json:"not_after"`
	SerialNumber   string `json:"serial_number"`
}

// Queries crt.sh for any related certificates in the transparent certificate logs (https://certificate.transparency.dev/)
func CheckCertLogs(target types.Target) (map[string]struct{}, error) {
	var (
		err         error
		domainNames = map[string]struct{}{}

		resp  *http.Response
		body  []byte
		certs = []Cert{}

		httpClient = http.Client{
			Timeout: 5 * time.Second,
		}
	)

	slog.Debug("subDomainScan: Checking cert logs started")

	resp, err = httpClient.Get("https://crt.sh?output=json&q=" + target.Hostname()) // Make the request
	if err != nil {
		if os.IsTimeout(err) {
			// A timeout error occurred
			return domainNames, errors.New("timeout while fetching subdomain data from crt.sh")
		}
		return domainNames, errors.New("error retrieving the response from crt.sh: " + err.Error())
	}

	defer func() {
		if closeErr := resp.Body.Close(); closeErr != nil {
			slog.Debug("subDomainScan: Error closing response body", "error", closeErr)
		}
	}()

	if resp.StatusCode == http.StatusOK {

		slog.Debug("subDomainScan: Cert logs response received", "status", resp.StatusCode)

		body, err = io.ReadAll(resp.Body) // Read the response
		if err != nil {
			return domainNames, errors.New("error reading the response from crt.sh: " + err.Error())
		}
		if closeErr := resp.Body.Close(); closeErr != nil {
			slog.Debug("subDomainScan: Error closing response body", "error", closeErr)
		}

		err = json.Unmarshal(body, &certs) // Parse the json
		if err != nil {
			return domainNames, errors.New("error parsing the response from crt.sh: " + err.Error() + "\n" + string(body))
		}

		for _, cert := range certs {
			if strings.HasSuffix(cert.CommonName, target.ParsedUrl().Host) { // Clear out third-parties
				if _, ok := domainNames[cert.CommonName]; !ok { // Clear out duplicates
					domainNames[cert.CommonName] = struct{}{}
				}
			}
		}
	} else {
		slog.Error("subDomainScan: Cert logs response not 200", "status", resp.StatusCode)
	}

	slog.Debug("subDomainScan: Checking cert logs completed")

	return domainNames, nil
}
