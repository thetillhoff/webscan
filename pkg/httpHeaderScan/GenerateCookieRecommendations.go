package httpHeaderScan

import (
	"net/http"
)

func GenerateCookieRecommendations(response *http.Response) (map[string][]string, []string) {
	var (
		allCookieRecommendations = map[string][]string{}
		otherRecommendations     = []string{}

		cookieRecommendations []string
	)

	// Cookies
	for _, cookie := range response.Cookies() {
		cookieRecommendations = []string{} // Initialize for every cookie

		if cookie.Domain != "" {
			// Domain should not be set, because if it's set cookies will be passed to subdomains as well
			cookieRecommendations = append(cookieRecommendations, "shouldn't have a Domain set. It's recommended to not set it, as it will send the cookie to subdomains if the domain is set.")
		}

		if !cookie.Secure {
			// Secure should be set, else the cookie will also be sent with http calls
			cookieRecommendations = append(cookieRecommendations, "should have Secure set. It's recommended for all cookies.")
		}

		if cookie.MaxAge == 0 && cookie.Expires.IsZero() {
			// Expiration aka MaxAge should be set, else it will be removed when the session is removed, which might be restored when the browser restarts
			cookieRecommendations = append(cookieRecommendations, "should have a max age or expiration set. It's recommended for all cookies. The default (delete after session) might restore it at browser restart.")
		}

		if !cookie.HttpOnly {
			// HttpOnly should be set, if not used otherwise with JS
			cookieRecommendations = append(cookieRecommendations, "shouldn't be accessible via javascript. It's recommended to use localStorage or server-side storage and api-calls for javascript-related data instead.")
		}

		if cookie.SameSite == http.SameSiteDefaultMode {
			// SameSite should be set (to lax or strict), because the default differ per browser
			cookieRecommendations = append(cookieRecommendations, "should have an explicit value for SameSite set. It's recommended to be set explicitly, because the default differs per browser")
		} else if cookie.SameSite == http.SameSiteNoneMode {
			cookieRecommendations = append(cookieRecommendations, "shouldn't have SameSite set to None. It's recommended to set it to another value (Lax or Strict).")
		}

		// Path is irrelevant for security & there's no best-practice

		// if strings.Contains(strings.ToLower(cookie.Name), "id") {
		// 	// Cookie name shouldn't lead on it's purpose. 'id' in it's name is already telling too much.
		// 	cookieRecommendations = append(cookieRecommendations, "shouldn't contain 'id' in it's name. It's recommended to obscure what a cookie is used for by using non-descriptive names.")
		// }

		if len(cookie.Value) < 16 {
			// Cookie value length should be >=128 bits which is >= 16 chars to prevent against brute force attacks
			cookieRecommendations = append(cookieRecommendations, "should have a length of >=16 characters if it's an ID to prevent against brute force attacks.")
		}

		if len(cookie.Value) > 450 {
			// Cookie value length close to 512 chars, which is the upper limit
			cookieRecommendations = append(cookieRecommendations, "should have a length of <450 characters, because that's quite close to 512 (4k), the upper limit in many browsers.")
		}

		if len(cookieRecommendations) > 0 {
			allCookieRecommendations[cookie.Name] = cookieRecommendations
		}
	}

	if len(response.Cookies()) > 10 {
		otherRecommendations = append(otherRecommendations, "More then 10 cookies detected. Are all of them really necessary? Think about storing session information on server-side and localStorage on client-side.")
	}

	return allCookieRecommendations, otherRecommendations
}
