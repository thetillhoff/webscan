package webscan

import (
	"errors"
	"fmt"
)

func (engine Engine) ScanMailConfig(inputUrl string) (Engine, error) {

	fmt.Println("Scanning mail config...")

	if engine.SubdomainScan {
		if engine.DkimSelector != "" {
			engine.mailConfigRecommendations = engine.dnsScanEngine.CheckMailSecurity(inputUrl, engine.DkimSelector)
		} else {
			return engine, errors.New("DKIM selector required")
		}
	}

	return engine, nil
}
