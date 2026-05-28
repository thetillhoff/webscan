package knownFilesScan

import (
	"bufio"
	"bytes"
	"strings"
	"time"
)

func analyzeSecurityTxt(body []byte) (observations []string, recommendations []string) {
	var (
		hasContact bool
		hasExpires bool
	)

	scanner := bufio.NewScanner(bytes.NewReader(body))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.SplitN(line, ":", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		switch strings.ToLower(key) {
		case "contact":
			hasContact = true
			observations = append(observations, "Contact: "+value)
		case "expires":
			hasExpires = true
			// RFC 9116 specifies ISO 8601 / RFC 3339 format
			t, err := time.Parse(time.RFC3339, value)
			if err != nil {
				recommendations = append(recommendations, "Expires field has invalid format (expected RFC 3339, e.g. 2006-01-02T15:04:05Z): "+value)
			} else if t.Before(time.Now()) {
				recommendations = append(recommendations, "Expires date is in the past: "+value)
			} else {
				observations = append(observations, "Expires: "+value)
			}
		}
	}

	if !hasContact {
		recommendations = append(recommendations, "Required field 'Contact' is missing")
	}
	if !hasExpires {
		recommendations = append(recommendations, "Required field 'Expires' is missing")
	}

	return observations, recommendations
}
