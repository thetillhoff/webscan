package status

func (status *Status) SpinningComplete(message string) {

	if status.isTTY {

		status.stopTicking()                                                // Stop time from triggering further updates
		status.Complete("  " + status.spinnerCompleteChar + "  " + message) // Set last message
		status.spinnerMessage = ""                                          // Reset
		status.spinnerIndex = 0                                             // Reset

	}
}
