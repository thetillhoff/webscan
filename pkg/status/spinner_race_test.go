package status

import (
	"bytes"
	"sync"
	"testing"
)

// Run with -race: concurrent SpinningXOfUpdate calls (as the port/tls/cipher
// scans do) must not race on spinner state or panic on the ticker channel.
func TestSpinnerConcurrentUpdates(t *testing.T) {
	var mu sync.Mutex
	var buf bytes.Buffer
	s := NewStatus(true, &mu, &buf)
	s.isTTY = true // force the spinner path without a real terminal

	s.SpinningXOfInit(100, "scanning")

	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			s.SpinningXOfUpdate()
		}()
	}
	wg.Wait()

	s.SpinningComplete("done")
}
