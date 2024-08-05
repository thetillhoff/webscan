package httpHeaderScan

import (
	"errors"
	"log/slog"
	"strconv"
	"strings"
)

func validateHeaderStrictTransportSecurity(content string) error {
	var (
		err error

		directives = strings.Split(content, ";")

		maxAge int = -1

		includeSubdomains bool
		preload           bool
	)

	slog.Debug("httpHeaderScan: Validating HSTS started")

	for _, directive := range directives {
		directive = strings.TrimSpace(directive)
		if strings.HasPrefix(directive, "max-age=") {
			maxAge, err = strconv.Atoi(strings.TrimPrefix(directive, "max-age="))
			if err != nil {
				return err
			}
			if maxAge < 15768000 {
				return errors.New("max-age value should be increased in stages from " + strconv.Itoa(maxAge) + " to 63072000 (two years)")
			}
		} else if directive == "includesubdomains" {
			includeSubdomains = true
		} else if directive == "preload" {
			preload = true
		} else {
			return errors.New("Unknown directive '" + directive + "' should be removed")
		}
	}

	if maxAge == -1 {
		return errors.New("max-age value is required")
	}

	if includeSubdomains && !preload {
		return errors.New("subdomains are included, so preload can be enabled as well")
	}

	if !includeSubdomains && preload {
		return errors.New("'preload' requires 'includeSubdomains' to be set")
	}

	if !includeSubdomains && !preload {
		return errors.New("'includeSubdomains' and 'preload' should be enabled when you are sure your https setup works fine")
	}

	slog.Debug("httpHeaderScan: Validating HSTS completed")

	return nil
}
