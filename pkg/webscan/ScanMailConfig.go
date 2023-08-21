package webscan

import (
	"errors"
	"fmt"
)

func (engine Engine) ScanMailConfig() (Engine, error) {

	fmt.Println("Scanning mail config...")

	if engine.SubdomainScan {
		if engine.DkimSelector != "" {
			engine.mailConfigRecommendations = engine.DnsScanEngine.CheckMailSecurity(engine.url, engine.DkimSelector)
		} else {
			return engine, errors.New("DKIM selector required")
		}
	}

	return engine, nil
}
