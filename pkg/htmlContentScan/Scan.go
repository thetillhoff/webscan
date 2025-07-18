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

	"github.com/PuerkitoBio/goquery"
	"github.com/thetillhoff/webscan/v3/pkg/cachedHttpGetClient"
	"github.com/thetillhoff/webscan/v3/pkg/status"
	"github.com/thetillhoff/webscan/v3/pkg/types"
)

// TODO return image / media / video / audio / svgs / ... as well

type scanConfig struct {
	client         cachedHttpGetClient.Client
	overrideSchema types.Schema
}

// ConfigOption represents a configuration option for DNS scanning
type ConfigOption func(*scanConfig)

// WithClient sets the client
func WithClient(client cachedHttpGetClient.Client) ConfigOption {
	return func(sc *scanConfig) {
		sc.client = client
	}
}

// WithSchemaOverride sets the schema override
func WithSchemaOverride(schema types.Schema) ConfigOption {
	return func(sc *scanConfig) {
		sc.overrideSchema = schema
	}
}

func Scan(status *status.Status, target types.Target, options ...ConfigOption) (Result, error) {
	var (
		err    error
		result = Result{
			httpContentStylesheetSizes: map[string]int{},
			httpContentScriptSizes:     map[string]int{},
			httpContentRecommendations: []string{},
		}

		response     *http.Response
		responseBody []byte

		body      []byte
		messages  []string
		message   string
		document  *goquery.Document
		parsedUrl *url.URL

		stylesheetSources []string
		scriptSources     []string
	)

	slog.Debug("htmlContentScan: Scan started")

	config := &scanConfig{}
	for _, option := range options {
		option(config)
	}

	target.OverrideSchema(config.overrideSchema)

	status.SpinningUpdate(fmt.Sprintf("Scanning %s content...", target.Schema().String()))

	response, responseBody, err = config.client.Get(target.UrlString())
	if err != nil {
		return result, err
	}

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
				parsedUrl.Host = target.Hostname()                        // Add hostname
				parsedUrl.Path = path.Join(target.Path(), parsedUrl.Path) // Add path prefix
			}

			if !path.IsAbs(parsedUrl.Path) { // If not leading '/' in path
				parsedUrl.Path = "/" + parsedUrl.Path // Add leading '/'
			}

			_, body, err = config.client.Get(parsedUrl.String())
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
				parsedUrl.Host = target.Hostname()                        // Add hostname
				parsedUrl.Path = path.Join(target.Path(), parsedUrl.Path) // Add path prefix
			}

			if !path.IsAbs(parsedUrl.Path) { // If not leading '/' in path
				parsedUrl.Path = "/" + parsedUrl.Path // Add leading '/'
			}

			_, body, err = config.client.Get(parsedUrl.String())
			if err != nil {
				return result, err
			}
			result.httpContentScriptSizes[scriptSource] = len(body)
		}
	}

	result.httpContentRecommendations = messages

	status.SpinningComplete(fmt.Sprintf("Scan of %s content complete.", target.Schema().String()))

	slog.Debug("htmlContentScan: Scan completed")

	return result, nil
}
