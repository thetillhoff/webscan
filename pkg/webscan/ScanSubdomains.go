package webscan

import subdomainfinder "github.com/thetillhoff/webscan/pkg/subDomainFinder"

func (engine Engine) ScanSubdomains() (Engine, error) {
	var (
		err error
	)

	engine.subdomains, err = subdomainfinder.CheckCertLogs(engine.url)

	if err != nil {
		return engine, err
	}

	return engine, nil
}
