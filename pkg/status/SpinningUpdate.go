package status

import (
	"time"
)

func (status *Status) SpinningUpdate(message string) {

	if status.isTTY {

		status.spinnerMessage = message // Set message internally (so updates keep displaying it)
		status.updateSpinner()          // Display initial message

		if !status.spinning { // Make sure there is only ever one ticker routine active
			status.spinning = true
			status.startTicking() // Start timer
		}
	}
}

// Display or update displayed message
func (status *Status) updateSpinner() {
	status.Update("  " + status.nextSpinner() + "  " + status.spinnerMessage)
}

// Start timer to trigger updateSpinner()
func (status *Status) startTicking() {
	status.spinnerStop = make(chan struct{})
	ticker := time.NewTicker(status.SpinnerUpdateInterval)

	go func() {
		for {
			select {
			case <-status.spinnerStop:
				ticker.Stop()
				return
			case <-ticker.C:
				status.updateSpinner()
			}
		}
	}()
}

// Stop timer to trigger updateSpinner()
func (status *Status) stopTicking() {
	status.spinnerStop <- struct{}{}
	close(status.spinnerStop)
	status.spinning = false
}
