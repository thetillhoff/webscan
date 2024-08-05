package status

import "fmt"

func (status *Status) SpinningXOfComplete(message string) {

	status.SpinningComplete(fmt.Sprintf("(%d/%d) %s", status.spinningXOfCurrent, status.spinningXOfTotal, message))
}
