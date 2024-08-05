package status

import (
	"github.com/jwalton/gchalk"
)

func (status *Status) nextSpinner() string {
	status.spinnerIndex = (status.spinnerIndex + 1) % len(status.spinnerChars)
	return gchalk.Dim(status.spinnerChars[status.spinnerIndex])
}
