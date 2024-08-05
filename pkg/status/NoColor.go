package status

func (status *Status) DisableColor() {
	status.spinnerCompleteChar = "✓"
}
