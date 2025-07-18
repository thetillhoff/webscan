package webscan

import (
	"github.com/thetillhoff/webscan/v3/pkg/status"
	"github.com/thetillhoff/webscan/v3/pkg/types"
)

type Scanner[C types.ScanConfig, R any] interface {
	Scan(target types.Target, status *status.Status, config C) (R, error)
	PrintResult(result R)
}
