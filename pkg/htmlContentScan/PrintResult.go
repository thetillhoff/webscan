package htmlContentScan

import (
	"fmt"
	"io"
	"log/slog"
	"strings"
)

func PrintResult(result Result, schemaName string, out io.Writer) {
	var (
		messages = []string{}

		stylesheetFileCount = 0
		totalStylesheetSize = 0

		scriptFileCount = 0
		totalScriptSize = 0
	)

	slog.Debug("htmlContentScan: Printing result started")

	// TODO include images, custom fonts

	messages = append(messages, result.httpContentRecommendations...)

	// HTML

	messages = append(messages, "HTML size: "+printByteSize(result.httpContentHtmlSize))

	if result.httpContentHtmlSize > 0 { // Only print more information if len(body) > 0

		// Size of html
		if result.httpContentHtmlSize > 200000 { // Size is larger than 200kB
			messages = append(messages, "  It's recommended to be smaller than 200kB.")
		}

		// Size of inline style
		if result.httpContentInlineStyleSize > 0 {
			messages = append(messages, "  Of this are inline Stylesheet (!= inline styles): "+printByteSize(result.httpContentInlineStyleSize))
		}

		// Size of inline script
		if result.httpContentInlineScriptSize > 0 {
			messages = append(messages, "  Of this are inline Script: "+printByteSize(result.httpContentInlineScriptSize))
		}

		// Size of external stylesheets

		if len(result.httpContentStylesheetSizes) > 0 {

			for _, size := range result.httpContentStylesheetSizes {
				stylesheetFileCount = stylesheetFileCount + 1
				totalStylesheetSize = totalStylesheetSize + size
			}
			messages = append(messages, "Total size of external CSS files: "+printByteSize(totalStylesheetSize))

		}

		// Size of external scripts

		if len(result.httpContentScriptSizes) > 0 {

			for _, size := range result.httpContentScriptSizes {
				scriptFileCount = scriptFileCount + 1
				totalScriptSize = totalScriptSize + size
			}
			messages = append(messages, "total size of external JS files: "+printByteSize(totalScriptSize))

		}

		// Total size

		totalSize := result.httpContentHtmlSize + totalStylesheetSize + totalScriptSize
		messages = append(messages, "Total download size (without media): "+printByteSize(totalSize))
	}

	if len(messages) > 0 {
		if _, err := fmt.Fprintf(out, "\n## %s content scan results\n\n", strings.ToUpper(schemaName)); err != nil {
			slog.Debug("htmlContentScan: Error writing to output", "error", err)
		}

		for _, message := range messages {
			if _, err := fmt.Fprintf(out, "%s\n", message); err != nil {
				slog.Debug("htmlContentScan: Error writing to output", "error", err)
			}
		}
	} else {
		slog.Debug("htmlContentScan: No information found")
	}

	slog.Debug("htmlContentScan: Printing result completed")
}
