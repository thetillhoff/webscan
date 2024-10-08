package htmlContentScan

import (
	"fmt"
	"log/slog"
	"strings"
)

func PrintResult(result Result, schemaName string) {
	var (
		sizeMessages = []string{}

		stylesheetFileCount int = 0
		totalStylesheetSize int = 0

		scriptFileCount int = 0
		totalScriptSize int = 0
	)

	slog.Debug("htmlContentScan: Printing result started")

	// First calculate all displayed values, then display them in one go

	// TODO

	//

	fmt.Printf("\n\n## %s content scan results\n\n", strings.ToUpper(schemaName))

	// TODO include images, custom fonts

	for _, message := range result.httpContentRecommendations {
		fmt.Println(message)
	}

	// HTML

	sizeMessages = append(sizeMessages, "HTML size: "+printByteSize(result.httpContentHtmlSize))

	if result.httpContentHtmlSize > 0 { // Only print more information if len(body) > 0

		// Size of html
		if result.httpContentHtmlSize > 200000 { // Size is larger than 200kB
			sizeMessages = append(sizeMessages, "  It's recommended to be smaller than 200kB.")
		}

		// Size of inline style
		if result.httpContentInlineStyleSize > 0 {
			sizeMessages = append(sizeMessages, "  Of this are inline Stylesheet (!= inline styles): "+printByteSize(result.httpContentInlineStyleSize))
		}

		// Size of inline script
		if result.httpContentInlineScriptSize > 0 {
			sizeMessages = append(sizeMessages, "  Of this are inline Script: "+printByteSize(result.httpContentInlineScriptSize))
		}

		// Size of external stylesheets

		if len(result.httpContentStylesheetSizes) > 0 {

			for _, size := range result.httpContentStylesheetSizes {
				stylesheetFileCount = stylesheetFileCount + 1
				totalStylesheetSize = totalStylesheetSize + size
			}
			sizeMessages = append(sizeMessages, "Total size of external CSS files: "+printByteSize(totalStylesheetSize))

		}

		// Size of external scripts

		if len(result.httpContentScriptSizes) > 0 {

			for _, size := range result.httpContentScriptSizes {
				scriptFileCount = scriptFileCount + 1
				totalScriptSize = totalScriptSize + size
			}
			sizeMessages = append(sizeMessages, "total size of external JS files: "+printByteSize(totalScriptSize))

		}

		// Total size

		totalSize := result.httpContentHtmlSize + totalStylesheetSize + totalScriptSize
		sizeMessages = append(sizeMessages, "Total download size (without media): "+printByteSize(totalSize))

		if len(result.httpContentRecommendations) > 0 { // Intermediate newline only needed if other text was already written
			fmt.Println()
		}

		for _, sizeMessage := range sizeMessages {
			fmt.Println(sizeMessage)
		}

	}

	slog.Debug("htmlContentScan: Printing result completed")

}
