package webscan

import (
	"fmt"
	"strconv"
)

func (engine Engine) PrintHttpContentScanResults() {
	var (
		sizeMessages = []string{}

		stylesheetFileCount int     = 0
		totalStylesheetSize float32 = 0

		scriptFileCount int     = 0
		totalScriptSize float32 = 0
	)

	// TODO include images, custom fonts

	if len(engine.httpContentRecommendations) > 0 {
		fmt.Println()
		for _, message := range engine.httpContentRecommendations {
			fmt.Println(message)
		}
	}

	// HTML

	sizeMessages = append(sizeMessages, "HTML size: "+strconv.Itoa(engine.httpContentHtmlSize/1000)+"kb")
	if engine.httpContentHtmlSize > 200*1000 { // Size is larger than 200kb
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
		// for stylesheetSource, size := range engine.httpContentStylesheetSizes {
		// 	fmt.Println("References external Stylesheet with " + strconv.FormatFloat(float64(size), 'f', 1, 64) + "kb at " + stylesheetSource)
		stylesheetFileCount = stylesheetFileCount + 1
		totalStylesheetSize = totalStylesheetSize + size
	}

	// Scripts

	for _, size := range engine.httpContentScriptSizes {
		// for scriptSource, size := range engine.httpContentScriptSizes {
		// 	fmt.Println("References external script with " + strconv.FormatFloat(float64(size), 'f', 1, 64) + "kb at " + scriptSource)
		scriptFileCount = scriptFileCount + 1
		totalScriptSize = totalScriptSize + size
	}

	// Totals

	if totalStylesheetSize > 0 {
		sizeMessages = append(sizeMessages, "Total size of "+strconv.Itoa(stylesheetFileCount)+" external stylesheets: "+strconv.FormatFloat(float64(totalStylesheetSize), 'f', 1, 64)+"kb")
	}

	if totalScriptSize > 0 {
		sizeMessages = append(sizeMessages, "Total size of "+strconv.Itoa(scriptFileCount)+" external scripts: "+strconv.FormatFloat(float64(totalScriptSize), 'f', 1, 64)+"kb")
	}

	totalsize := float32(engine.httpContentHtmlSize/1000) + totalStylesheetSize + totalScriptSize
	sizeMessages = append(sizeMessages, "Total download size: "+strconv.FormatFloat(float64(totalsize), 'f', 1, 64)+"kb")

	fmt.Println()
	for _, sizeMessage := range sizeMessages {
		fmt.Println(sizeMessage)
	}

}
