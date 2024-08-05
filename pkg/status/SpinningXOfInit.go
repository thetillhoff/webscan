package status

import (
	"fmt"
)

func (status *Status) SpinningXOfInit(total int, message string) {

	status.spinningXOfTotal = total // Set total
	status.spinningXOfCurrent = 0   // Reset
	status.spinningXOfMessage = message

	status.SpinningUpdate(fmt.Sprintf("(%d/%d) %s", status.spinningXOfCurrent, status.spinningXOfTotal, status.spinningXOfMessage))
}
