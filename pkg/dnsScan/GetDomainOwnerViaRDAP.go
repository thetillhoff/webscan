package dnsScan

import (
	"log/slog"
	"strings"

	"github.com/openrdap/rdap"
)

func GetDomainOwnerViaRDAP(url string) ([]string, error) {
	var (
		err        error
		client     *rdap.Client = &rdap.Client{}
		rdapDomain *rdap.Domain

		emailDomainsUnique = map[string]struct{}{}
		emailDomains       = []string{}
	)

	slog.Debug("dnsScan: Getting domain owner via rdap started")

	rdapDomain, err = client.QueryDomain(url)
	if err != nil {
		return emailDomains, err
	}

	for _, entity := range rdapDomain.Entities {
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

	slog.Debug("dnsScan: Getting domain owner via rdap completed")

	return emailDomains, nil
}
