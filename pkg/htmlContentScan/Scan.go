package htmlContentScan

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"net/url"
	"path"
	"path/filepath"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/thetillhoff/webscan/v3/pkg/httpClient"
	"github.com/thetillhoff/webscan/v3/pkg/status"
)

// TODO return image / media / video / audio / svgs / ... as well

func Scan(status *status.Status, inputUrl string, response *http.Response, responseBody []byte, client httpClient.Client, schemaName string) (Result, error) {
	var (
		err    error
		result = Result{
			httpContentStylesheetSizes: map[string]int{},
			httpContentScriptSizes:     map[string]int{},
			httpContentRecommendations: []string{},
		}

		body     []byte
		messages []string
		message  string
		document *goquery.Document

		stylesheetSources []string
		scriptSources     []string

		parsedUrl *url.URL
	)

	slog.Debug("htmlContentScan: Scan started")

	status.SpinningUpdate(fmt.Sprintf("Scanning %s content...", schemaName))

	message = CheckCompression(response) // Check whether compression was used
	if message != "" {
		messages = append(messages, message)
	}

	result.httpContentHtmlSize = len(responseBody)

	// HTML

	if len(responseBody) > 0 {

		message, err = ValidateHtml(responseBody)
		if err != nil { // HTML content is not valid HTML
			return result, errors.New("invalid html: " + err.Error())
		}
		if message != "" { // HTML content is valid, but not HTML5
			messages = append(messages, message)
		}

		// Load the HTML document
		document, err = goquery.NewDocumentFromReader(bytes.NewReader(responseBody))
		if err != nil {
			return result, err
		}

		result.httpContentInlineStyleSize, err = GetInlineStyleSize(document) // Inline style
		if err != nil {
			return result, err
		}

		result.httpContentInlineScriptSize, err = GetInlineScriptSize(document) // Inline script
		if err != nil {
			return result, err
		}

		// Stylesheets

		stylesheetSources = GetExternalStylesheetSources(document)
		scriptSources = GetExternalScriptSources(document)

		if len(stylesheetSources) > 0 {
			result.httpContentStylesheetSizes = map[string]int{}
		}

		for _, stylesheetSource := range stylesheetSources {
			parsedUrl, err = url.Parse(stylesheetSource)
			if err != nil {
				return result, err
			}

			if path.Ext(parsedUrl.Path) != ".css" {
				messages = append(messages, "External stylesheets should have `.css` set as file extension. Got "+stylesheetSource)
			}
			if parsedUrl.IsAbs() { // Includes a scheme
				if parsedUrl.Scheme != "https" {
					log.Println(parsedUrl.Scheme)
					messages = append(messages, "External stylesheets should only be referenced via HTTPS. Got "+stylesheetSource)
				}
			} else { // Doesn't include a scheme
				parsedUrl.Scheme = "https" // Add scheme
			}

			if parsedUrl.Host == "" { // Doesn't include hostname
				if strings.Contains(inputUrl, "/") {
					inputUrlParts := strings.SplitN(inputUrl, "/", 2)
					parsedUrl.Host, parsedUrl.Path = inputUrlParts[0], inputUrlParts[1]+parsedUrl.Path
				} else {
					parsedUrl.Host = inputUrl // Add hostname
				}
			}

			if !filepath.IsAbs(parsedUrl.Path) { // If not leading '/' in path
				parsedUrl.Path = "/" + parsedUrl.Path // Add leading '/'
			}

			body, err = client.GetBodyForRequest("GET", parsedUrl.String())
			if err != nil {
				return result, err
			}
			result.httpContentStylesheetSizes[stylesheetSource] = len(body)
		}

		// Scripts

		if len(scriptSources) > 0 {
			result.httpContentScriptSizes = map[string]int{}
		}

		for _, scriptSource := range scriptSources {
			parsedUrl, err = url.Parse(scriptSource)
			if err != nil {
				return result, err
			}

			if path.Ext(parsedUrl.Path) != ".js" {
				messages = append(messages, "External scripts should have `.js` set as file extension. Got "+scriptSource)
			}
			if parsedUrl.IsAbs() { // Includes a scheme
				if parsedUrl.Scheme != "https" {
					log.Println(parsedUrl.Scheme)
					messages = append(messages, "External scripts should only be referenced via HTTPS. Got "+scriptSource)
				}
			} else { // Doesn't include a scheme
				parsedUrl.Scheme = "https" // Add scheme
			}

			if parsedUrl.Host == "" { // Doesn't include hostname
				if strings.Contains(inputUrl, "/") {
					inputUrlParts := strings.SplitN(inputUrl, "/", 2)
					parsedUrl.Host, parsedUrl.Path = inputUrlParts[0], inputUrlParts[1]+parsedUrl.Path
				} else {
					parsedUrl.Host = inputUrl // Add hostname
				}
			}

			if !filepath.IsAbs(parsedUrl.Path) { // If not leading '/' in path
				parsedUrl.Path = "/" + parsedUrl.Path // Add leading '/'
			}

			body, err = client.GetBodyForRequest("GET", parsedUrl.String())
			if err != nil {
				return result, err
			}
			result.httpContentScriptSizes[scriptSource] = len(body)
		}
	}

	result.httpContentRecommendations = messages

	status.SpinningComplete(fmt.Sprintf("Scan of %s content complete.", schemaName))

	slog.Debug("htmlContentScan: Scan completed")

	return result, nil
}
