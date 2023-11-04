package htmlContentScan

import "net/http"

func CheckCompression(response *http.Response) string {

	// Compression (not a header)
	if response.Uncompressed { // Response was decompressed
		// messages = append(messages, "Compression: enabled") // TODO should this be printed?
		// Due to chunked encoding it is very often not possible to determine body length in a straight forward way. Therefore it's not possible to calculate a compression ratio in an easy way.
	} else { // Response wasn't compressed in the first place
		return "Compression is not enabled"
	}

	return ""
	// TODO if compression was enabled, download the page a second time with compression disabled and compare body sizes -> see comments above
}
