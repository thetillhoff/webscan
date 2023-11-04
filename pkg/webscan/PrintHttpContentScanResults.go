package webscan

import (
	"fmt"
	"strconv"
)

func (engine Engine) PrintHttpContentScanResults() {
	var (
		sizeMessages = []string{}

		stylesheetFileCount int     = 0
		totalStylesheetSize float64 = 0

		scriptFileCount int     = 0
		totalScriptSize float64 = 0
	)

	fmt.Printf("\n\n--- HTTP content scan results ---\n")

	// TODO include images, custom fonts

	if len(engine.httpContentRecommendations) > 0 {
		for _, message := range engine.httpContentRecommendations {
			fmt.Println(message)
		}
	}

	// HTML

	sizeMessages = append(sizeMessages, "HTML size: "+strconv.FormatFloat(engine.httpContentHtmlSizekB, 'f', -1, 64)+"kb")
	if engine.httpContentHtmlSizekB > 200 { // Size is larger than 200kb
		sizeMessages = append(sizeMessages, "  It's recommended to be smaller than 200kb.")
	}

	// Inline style
	if engine.httpContentInlineStyleSize > 0 {
		sizeMessages = append(sizeMessages, "  Of this are inline Stylesheet (!= inline styles): "+strconv.Itoa(engine.httpContentInlineStyleSize/1000)+"kb")
	}

	// Inline script
	if engine.httpContentInlineScriptSize > 0 {
		sizeMessages = append(sizeMessages, "  Of this are inline Script: "+strconv.Itoa(engine.httpContentInlineScriptSize/1000)+"kb")
	}

	// Stylesheets

	for _, size := range engine.httpContentStylesheetSizes {
		stylesheetFileCount = stylesheetFileCount + 1
		totalStylesheetSize = totalStylesheetSize + size
	}
	sizeMessages = append(sizeMessages, "External CSS size: "+strconv.FormatFloat(totalStylesheetSize, 'f', -1, 64)+"kb")

	// Scripts

	for _, size := range engine.httpContentScriptSizes {
		scriptFileCount = scriptFileCount + 1
		totalScriptSize = totalScriptSize + size
	}
	sizeMessages = append(sizeMessages, "External JS size: "+strconv.FormatFloat(totalScriptSize, 'f', -1, 64)+"kb")

	// Total

	if engine.httpContentHtmlSizekB > 0 {
		totalsize := engine.httpContentHtmlSizekB + totalStylesheetSize + totalScriptSize
		sizeMessages = append(sizeMessages, "Total download size (without media): "+strconv.FormatFloat(totalsize, 'f', -1, 64)+"kb")

		fmt.Println()
		for _, sizeMessage := range sizeMessages {
			fmt.Println(sizeMessage)
		}
	}

}
