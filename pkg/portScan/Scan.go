package portScan

import (
	"time"

	"github.com/thetillhoff/webscan/v5/pkg/status"
	"github.com/thetillhoff/webscan/v5/pkg/types"
)

type scanConfig struct {
	aRecords    []string
	aaaaRecords []string
	advanced    bool
	fullRange   bool
	timeout     time.Duration
}

type ConfigOption = types.Option[scanConfig]

// WithARecords sets the A records to scan
func WithARecords(aRecords []string) ConfigOption {
	return func(sc *scanConfig) {
		sc.aRecords = aRecords
	}
}

// WithAAAARecords sets the AAAA records to scan
func WithAAAARecords(aaaaRecords []string) ConfigOption {
	return func(sc *scanConfig) {
		sc.aaaaRecords = aaaaRecords
	}
}

// WithAdvanced enables advanced port scanning
func WithAdvanced(advanced bool) ConfigOption {
	return func(sc *scanConfig) {
		sc.advanced = advanced
	}
}

// WithFullRange scans all 65535 TCP ports instead of the common-port list. Only
// takes effect together with WithAdvanced(true).
func WithFullRange(fullRange bool) ConfigOption {
	return func(sc *scanConfig) {
		sc.fullRange = fullRange
	}
}

// WithTimeout sets the per-connection timeout for port scanning
func WithTimeout(timeout time.Duration) ConfigOption {
	return func(sc *scanConfig) {
		sc.timeout = timeout
	}
}

func Scan(target types.Target, status *status.Status, options ...ConfigOption) (Result, error) {
	config := &scanConfig{
		timeout: 5 * time.Second,
	}
	types.ApplyOptions(config, options)

	switch {
	case config.advanced && target.Port() == "" && target.Schema() == types.NONE:
		return AdvancedScan(status, config.aRecords, config.aaaaRecords, config.timeout, config.fullRange)
	default:
		return SimpleScan(target, status, config.aRecords, config.aaaaRecords, config.timeout)
	}
}
