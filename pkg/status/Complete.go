package status

func (status *Status) Complete(message string) {

	if status.isTTY {

		status.Update(message + "\n")

	}
}
