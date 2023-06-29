package webscan

import "errors"

func (engine Engine) ScanMailConfig() (Engine, error) {

	if engine.SubdomainScan {
		if engine.DkimSelector != "" {
			engine.mailConfigRecommendations = engine.DnsScanEngine.CheckMailSecurity(engine.url, engine.DkimSelector)
		} else {
			return engine, errors.New("DKIM selector required")
		}
	}

	return engine, nil
}
