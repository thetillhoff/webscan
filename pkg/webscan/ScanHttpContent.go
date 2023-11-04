package webscan

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"net/url"
	"path"
	"path/filepath"

	"github.com/PuerkitoBio/goquery"
	htmlContentScan "github.com/thetillhoff/webscan/pkg/htmlContentScan"
)

// TODO return image / media / video / audio / svgs / ... as well

func (engine Engine) ScanHttpContent(inputUrl string) (Engine, error) {
	var (
		err error

		body     []byte
		messages []string
		message  string
		document *goquery.Document

		stylesheetSources []string
		scriptSources     []string

		parsedUrl *url.URL
	)

	fmt.Println("Scanning HTTP content...")

	message = htmlContentScan.CheckCompression(engine.response) // Check whether compression was used
	if message != "" {
		messages = append(messages, message)
	}

	body, err = io.ReadAll(engine.response.Body) // Read the body
	if err != nil {
		return engine, err
	}
	defer engine.response.Body.Close()

	engine.httpContentHtmlSizekB = float64(len(body)) / 1000 // Get html size

	// HTML

	if len(body) > 0 {

		message, err = htmlContentScan.ValidateHtml(body)
		if err != nil { // HTML content is not valid HTML
			return engine, errors.New("invalid html: " + err.Error())
		}
		if message != "" { // HTML content is valid, but not HTML5
			messages = append(messages, message)
		}

		// Load the HTML document
		document, err = goquery.NewDocumentFromReader(bytes.NewReader(body))
		if err != nil {
			return engine, err
		}

		engine.httpContentInlineStyleSize, err = htmlContentScan.GetInlineStyleSize(document) // Inline style
		if err != nil {
			return engine, err
		}

		engine.httpContentInlineScriptSize, err = htmlContentScan.GetInlineScriptSize(document) // Inline script
		if err != nil {
			return engine, err
		}

		// Stylesheets

		stylesheetSources = htmlContentScan.GetExternalStylesheetSources(document)
		scriptSources = htmlContentScan.GetExternalScriptSources(document)

		if len(stylesheetSources) > 0 {
			engine.httpContentStylesheetSizes = map[string]float64{}
		}

		for _, stylesheetSource := range stylesheetSources {
			parsedUrl, err = url.Parse(stylesheetSource)
			if err != nil {
				return engine, err
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
				parsedUrl.Host = inputUrl // Add hostname
			}

			if !filepath.IsAbs(parsedUrl.Path) { // If not leading '/' in path
				parsedUrl.Path = "/" + parsedUrl.Path // Add leading '/'
			}
			engine.httpContentStylesheetSizes[stylesheetSource], err = htmlContentScan.GetSizeInKb(parsedUrl.String())
			if err != nil {
				return engine, err
			}
		}

		// Scripts

		if len(scriptSources) > 0 {
			engine.httpContentScriptSizes = map[string]float64{}
		}

		for _, scriptSource := range scriptSources {
			parsedUrl, err = url.Parse(scriptSource)
			if err != nil {
				return engine, err
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
				parsedUrl.Host = inputUrl // Add hostname
			}

			if !filepath.IsAbs(parsedUrl.Path) { // If not leading '/' in path
				parsedUrl.Path = "/" + parsedUrl.Path // Add leading '/'
			}
			engine.httpContentScriptSizes[scriptSource], err = htmlContentScan.GetSizeInKb(parsedUrl.String())
			if err != nil {
				return engine, err
			}
		}
	}

	engine.httpContentRecommendations = messages

	return engine, nil
}
