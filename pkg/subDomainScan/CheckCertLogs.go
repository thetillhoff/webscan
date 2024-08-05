package subDomainScan

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"time"
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
func CheckCertLogs(url string) (map[string]struct{}, error) {
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

	resp, err = httpClient.Get("https://crt.sh?output=json&q=" + url) // Make the request
	if os.IsTimeout(err) {
		// A timeout error occurred
		return domainNames, errors.New("timeout while fetching subdomain data from crt.sh")
	}
	if err != nil {
		log.Fatalln("error retrieving the response from crt.sh:", err)
	}

	log.Println(resp.StatusCode)

	body, err = io.ReadAll(resp.Body) // Read the response
	if err != nil {
		log.Fatalln("error reading the response from crt.sh:", err)
	}
	resp.Body.Close()

	err = json.Unmarshal(body, &certs) // Parse the json
	if err != nil {
		log.Fatalln("error parsing the response from crt.sh:", err, "\n", string(body))
	}

	for _, cert := range certs {
		if strings.HasSuffix(cert.CommonName, url) { // Clear out third-parties
			if _, ok := domainNames[cert.CommonName]; !ok { // Clear out duplicates
				domainNames[cert.CommonName] = struct{}{}
			}
		}
	}

	slog.Debug("subDomainScan: Checking cert logs completed")

	return domainNames, nil
}
