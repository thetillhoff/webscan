package httpHeaderScan

import (
	"errors"
	"strings"
)

func scanCookieChars(cookie string) error {
	var (
		err          error
		invalidChars = []string{ // Taken from https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Set-Cookie#attributes
			"(",
			")",
			"<",
			">",
			"@",
			",",
			";",
			":",
			"\\",
			"\"",
			"/",
			"[",
			"]",
			"?",
			"=",
			"{",
			"}",
		}

		detectedInvalidChars = ""
	)

	for _, invalidChar := range invalidChars {
		if strings.Contains(cookie, invalidChar) {
			detectedInvalidChars = detectedInvalidChars + invalidChar
		}
	}

	switch len(detectedInvalidChars) {
	case 0:
		// Do nothing, as all characters are valid
	case 1:
		err = errors.New("Invalid character is `" + detectedInvalidChars + "`.")
	default:
		err = errors.New("Invalid characters are `" + detectedInvalidChars + "`.")
	}

	return err
}
