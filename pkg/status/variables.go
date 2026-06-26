package status

import "sync"

var (
	countMutex   sync.Mutex // To ensure multi-threaded apps don't mess up the global counting
	spinnerMutex sync.Mutex // Serializes all spinner state mutations (terminal is a shared global resource)
)
