package httpProtocolScan

// Takes url
// and whether http and/or https should be checked
// and checks which HTTP versions the server can speak for each
// Return available versions on http, then https and finally a potential error
func CheckHttpVersions(url string, checkHttp bool, checkHttps bool) ([]string, []string, error) {
	var (
		err                    error
		availableHttpVersions  = []string{}
		availableHttpsVersions = []string{}

		httpVersion1 string
		httpVersion2 string
		httpVersion3 string

		httpsVersion1 string
		httpsVersion2 string
		httpsVersion3 string
	)

	// HTTP versions:
	// 0.9 -> obsolete
	// 1.0 -> obsolete
	// 1.1
	// 2
	// 3 QUIC

	if checkHttp {
		protocol := "http://"

		httpVersion1, err = checkHttp1(protocol + url)
		if err == nil {
			availableHttpVersions = append(availableHttpVersions, httpVersion1)
		}

		httpVersion2, err = checkHttp2(protocol + url)
		if err == nil {
			availableHttpVersions = append(availableHttpVersions, httpVersion2)
		}

		httpVersion3, err = checkHttp3(protocol + url)
		if err == nil {
			availableHttpVersions = append(availableHttpVersions, httpVersion3)
		}
	}

	if checkHttps {
		protocol := "https://"

		httpsVersion1, err = checkHttp1(protocol + url)
		if err == nil {
			availableHttpsVersions = append(availableHttpsVersions, httpsVersion1)
		}

		httpsVersion2, err = checkHttp2(protocol + url)
		if err == nil {
			availableHttpsVersions = append(availableHttpsVersions, httpsVersion2)
		}

		httpsVersion3, err = checkHttp3(protocol + url)
		if err == nil {
			availableHttpsVersions = append(availableHttpsVersions, httpsVersion3)
		}
	}

	return availableHttpVersions, availableHttpsVersions, nil
}
