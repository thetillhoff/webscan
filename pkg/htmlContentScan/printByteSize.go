package htmlContentScan

import "strconv"

func printByteSize(size int) string {
	if size > 1000 {
		return strconv.Itoa(size/1000) + "kB"
	}
	return strconv.Itoa(size) + "B"
}
