package portScan

import (
	"github.com/thetillhoff/webscan/v3/pkg/status"
	"github.com/thetillhoff/webscan/v3/pkg/types"
)

type scanConfig struct {
	aRecords    []string
	aaaaRecords []string
	advanced    bool
}

// ConfigOption represents a configuration option for port scanning
type ConfigOption func(*scanConfig)

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

func Scan(target types.Target, status *status.Status, options ...ConfigOption) (Result, error) {
	// Apply configuration options
	config := &scanConfig{}
	for _, option := range options {
		option(config)
	}

	switch {
	case config.advanced && target.Port() == "" && target.Schema() == types.NONE:
		return AdvancedScan(status, config.aRecords, config.aaaaRecords)
	default:
		return SimpleScan(target, status, config.aRecords, config.aaaaRecords)
	}
}
