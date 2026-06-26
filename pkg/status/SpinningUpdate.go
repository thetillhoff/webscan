package status

import (
	"time"
)

func (status *Status) SpinningUpdate(message string) {

	if status.isTTY {

		spinnerMutex.Lock()
		defer spinnerMutex.Unlock()

		status.spinnerMessage = message // Set message internally (so updates keep displaying it)
		status.updateSpinnerLocked()    // Display initial message

		if !status.spinning { // Make sure there is only ever one ticker routine active
			status.spinning = true
			status.startTickingLocked() // Start timer
		}
		return
	}

	status.Println(message)
}

// Display or update displayed message. Caller must hold spinnerMutex.
func (status *Status) updateSpinnerLocked() {
	status.Update("  " + status.nextSpinner() + "  " + status.spinnerMessage)
}

// Start timer to trigger updateSpinner(). Caller must hold spinnerMutex.
func (status *Status) startTickingLocked() {
	// Buffered so stopTickingLocked never blocks while holding spinnerMutex.
	stop := make(chan struct{}, 1)
	status.spinnerStop = stop
	ticker := time.NewTicker(status.SpinnerUpdateInterval)

	go func() {
		for {
			select {
			case <-stop:
				ticker.Stop()
				return
			case <-ticker.C:
				spinnerMutex.Lock()
				status.updateSpinnerLocked()
				spinnerMutex.Unlock()
			}
		}
	}()
}

// Stop timer to trigger updateSpinner(). Caller must hold spinnerMutex.
func (status *Status) stopTickingLocked() {
	if !status.spinning {
		return
	}
	status.spinnerStop <- struct{}{} // buffered(1), won't block
	status.spinning = false
}
