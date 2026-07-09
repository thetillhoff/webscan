package ipScan

import (
	"io"
	"sync"
	"testing"

	"github.com/thetillhoff/webscan/v5/pkg/status"
	"github.com/thetillhoff/webscan/v5/pkg/types"
)

// No A/AAAA records must not be an error - it's a valid state (e.g. mail-only
// domains). The zero-IP path returns early before any network call.
func TestScanNoIPsIsNotError(t *testing.T) {
	target, err := types.NewTarget("example.com")
	if err != nil {
		t.Fatalf("NewTarget: %v", err)
	}
	st := status.NewStatus(true, &sync.Mutex{}, io.Discard)

	_, err = Scan(target, &st) // no WithARecords / WithAAAARecords
	if err != nil {
		t.Fatalf("expected nil error for zero IPs, got %v", err)
	}
}
