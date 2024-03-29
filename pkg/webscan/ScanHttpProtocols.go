package webscan

import (
	"fmt"

	httpProtocolScan "github.com/thetillhoff/webscan/pkg/httpProtocolScan"
)

func (engine Engine) ScanHttpProtocols(inputUrl string) (Engine, error) {
	var (
		err error
	)

	fmt.Println("Scanning HTTP protocols...") // Needed in more cases than just protocol scan

	// Scan HTTP / HTTPS for redirects
	engine.httpStatusCode, engine.httpRedirectLocation, engine.httpsStatusCode, engine.httpsRedirectLocation, err = httpProtocolScan.CheckHttpRedirects(inputUrl, engine.isAvailableViaHttp, engine.isAvailableViaHttps)
	if err != nil {
		return engine, err
	}

	// TODO check redirect from http zone apex to https www. prefix
	// TODO check redirect from https zone apex to https www. prefix
	// TODO check redirect from http www. prefix to https www. prefix
	// TODO check redirects to omit the port (it's unneeded if protocol is set and it's the default 80 or 443)

	// TODO follow redirects if desired -> Probably not here, but in Scan().
	// TODO Only check http versions when there is no redirect happening

	// Scan Http Versions
	engine.httpVersions, engine.httpsVersions, err = httpProtocolScan.CheckHttpVersions(inputUrl, engine.httpRedirectLocation != "", engine.httpsRedirectLocation != "")
	if err != nil {
		return engine, err
	}

	return engine, nil
}
