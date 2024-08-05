package status

import "fmt"

func (status *Status) SpinningXOfUpdate() {

	countMutex.Lock()
	status.spinningXOfCurrent = status.spinningXOfCurrent + 1
	countMutex.Unlock()

	status.SpinningUpdate(fmt.Sprintf("(%d/%d) %s", status.spinningXOfCurrent, status.spinningXOfTotal, status.spinningXOfMessage))

}
