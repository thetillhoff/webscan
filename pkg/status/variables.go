package status

import "sync"

var (
	countMutex sync.Mutex // To ensure multi-threaded apps don't mess up the global counting
)
