package ipScan

import (
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/openrdap/rdap"
)

func GetIPOwnerViaRDAP(ip string, timeout time.Duration) (string, error) {
	var (
		err    error
		client = &rdap.Client{
			HTTP: &http.Client{Timeout: timeout},
		}
		rdapIPNetwork *rdap.IPNetwork

		response           string
		emailDomainsUnique = map[string]struct{}{}
		emailDomains       = []string{}
	)

	slog.Debug("ipScan: Checking for ip owner via rdap started")

	rdapIPNetwork, err = client.QueryIP(ip)
	if err != nil {
		return "", err
	}

	for _, entity := range rdapIPNetwork.Entities {
		if entity.VCard != nil {

			if entity.VCard.Email() != "" {
				email := entity.VCard.Email()
				emailParts := strings.Split(email, "@")

				emailDomainsUnique[emailParts[len(emailParts)-1]] = struct{}{}
			}

			for _, subEntity := range entity.Entities {
				if subEntity.VCard != nil {

					if subEntity.VCard.Email() != "" {
						email := subEntity.VCard.Email()
						emailParts := strings.Split(email, "@")

						emailDomainsUnique[emailParts[len(emailParts)-1]] = struct{}{}
					}
				}
			}
		}
	}

	for emailDomain := range emailDomainsUnique {
		emailDomains = append(emailDomains, emailDomain)
	}
	response = strings.Join(emailDomains, ", ")

	slog.Debug("ipScan: Checking for ip owner via rdap completed")

	return response, nil
}
