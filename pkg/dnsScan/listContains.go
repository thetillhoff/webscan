package dnsScan

func listContains(list []string, element string) bool {

	for _, listElement := range list {
		if listElement == element {
			return true
		}
	}

	return false
}
