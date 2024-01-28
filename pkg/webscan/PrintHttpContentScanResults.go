package webscan

import (
	"fmt"
	"strconv"
)

func (engine Engine) PrintHttpContentScanResults() {
	var (
		sizeMessages = []string{}

		stylesheetFileCount int = 0
		totalStylesheetSize int = 0

		scriptFileCount int = 0
		totalScriptSize int = 0
	)

	if engine.HttpContentScan {

		fmt.Printf("\n\n--- HTTP content scan results ---\n")

		// TODO include images, custom fonts

		if len(engine.httpContentRecommendations) > 0 {
			for _, message := range engine.httpContentRecommendations {
				fmt.Println(message)
			}
		}

		// HTML

		sizeMessages = append(sizeMessages, "HTML size: "+strconv.Itoa(engine.httpContentHtmlSize/1000)+"kB")
		if engine.httpContentHtmlSize > 200000 { // Size is larger than 200kB
			sizeMessages = append(sizeMessages, "  It's recommended to be smaller than 200kB.")
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
		sizeMessages = append(sizeMessages, "External CSS size: "+strconv.Itoa(totalStylesheetSize/1000)+"kB")

		// Scripts

		for _, size := range engine.httpContentScriptSizes {
			scriptFileCount = scriptFileCount + 1
			totalScriptSize = totalScriptSize + size
		}
		sizeMessages = append(sizeMessages, "External JS size: "+strconv.Itoa(totalScriptSize/1000)+"kB")

		// Total

		if engine.httpContentHtmlSize > 0 {
			totalSize := engine.httpContentHtmlSize + totalStylesheetSize + totalScriptSize
			sizeMessages = append(sizeMessages, "Total download size (without media): "+strconv.Itoa(totalSize/1000)+"kB")

			if len(engine.httpContentRecommendations) > 0 { // Intermediate newline only needed if other text was already written
				fmt.Println()
			}

			for _, sizeMessage := range sizeMessages {
				fmt.Println(sizeMessage)
			}
		}
	}

}
