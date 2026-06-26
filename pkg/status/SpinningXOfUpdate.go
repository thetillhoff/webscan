package status

import "fmt"

func (status *Status) SpinningXOfUpdate() {

	countMutex.Lock()
	status.spinningXOfCurrent = status.spinningXOfCurrent + 1
	current := status.spinningXOfCurrent
	countMutex.Unlock()

	status.SpinningUpdate(fmt.Sprintf("(%d/%d) %s", current, status.spinningXOfTotal, status.spinningXOfMessage))

}
