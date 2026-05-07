package subDomainScan

import (
	"encoding/json"
	"fmt"
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

type crtShStatus int

const (
	crtShOK          crtShStatus = iota
	crtShTimeout     crtShStatus = iota
	crtShRateLimited crtShStatus = iota
	crtShServerError crtShStatus = iota
	crtShUnreachable crtShStatus = iota
)

func (s crtShStatus) String() string {
	switch s {
	case crtShTimeout:
		return "crt.sh timed out"
	case crtShRateLimited:
		return "crt.sh rate limited"
	case crtShServerError:
		return "crt.sh server error"
	case crtShUnreachable:
		return "crt.sh unreachable"
	default:
		return ""
	}
}

// Queries crt.sh for any related certificates in the transparent certificate logs (https://certificate.transparency.dev/)
func CheckCertLogs(target types.Target, timeout time.Duration) (map[string]struct{}, crtShStatus) {
	domainNames := map[string]struct{}{}

	httpClient := http.Client{Timeout: timeout}

	slog.Debug("subDomainScan: Checking cert logs started")

	resp, err := httpClient.Get("https://crt.sh/?output=json&q=" + target.Hostname())
	if err != nil {
		if os.IsTimeout(err) {
			slog.Info("subDomainScan: crt.sh request timed out", "error", err)
			return domainNames, crtShTimeout
		}
		slog.Info("subDomainScan: crt.sh unreachable", "error", err)
		return domainNames, crtShUnreachable
	}

	defer func() { _ = resp.Body.Close() }()

	switch {
	case resp.StatusCode == http.StatusOK:
		// success, parse below
	case resp.StatusCode == http.StatusTooManyRequests:
		slog.Info("subDomainScan: crt.sh rate limited", "status", resp.StatusCode)
		return domainNames, crtShRateLimited
	case resp.StatusCode >= 500:
		slog.Info("subDomainScan: crt.sh server error", "status", resp.StatusCode)
		return domainNames, crtShServerError
	default:
		slog.Info("subDomainScan: crt.sh unexpected status", "status", resp.StatusCode)
		return domainNames, crtShServerError
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		slog.Debug("subDomainScan: error reading crt.sh response", "error", err)
		return domainNames, crtShServerError
	}

	var certs []Cert
	if err := json.Unmarshal(body, &certs); err != nil {
		slog.Debug("subDomainScan: error parsing crt.sh response", "error", err, "body", fmt.Sprintf("%.200s", string(body)))
		return domainNames, crtShServerError
	}

	for _, cert := range certs {
		if strings.HasSuffix(cert.CommonName, target.ParsedUrl().Host) {
			domainNames[cert.CommonName] = struct{}{}
		}
	}

	slog.Debug("subDomainScan: Checking cert logs completed", "count", len(domainNames))

	return domainNames, crtShOK
}
