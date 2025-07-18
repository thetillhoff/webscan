package status

func (status *Status) SpinningComplete(message string) {

	if status.isTTY {

		status.stopTicking() // Stop time from triggering further updates
		// status.Complete("  " + status.spinnerCompleteChar + "  " + message) // Set last message
		status.Complete("\033[K\033[1A") // Set last message, which resets the line // TODO testing if this prints everything apart from the completion message
		// More info at https://tldp.org/HOWTO/Bash-Prompt-HOWTO/x361.html
		status.spinnerMessage = "" // Reset
		status.spinnerIndex = 0    // Reset

	}
}
