package dnsScan

func generateCartesianProduct(lists ...[]string) []string {
	var (
		results = []string{}

		firstList  []string
		secondList []string
		otherLists [][]string
	)

	firstList, otherLists = lists[0], lists[1:]

	if len(otherLists) > 1 {
		secondList = generateCartesianProduct(otherLists...)
	} else if len(otherLists) == 1 {
		secondList = otherLists[0]
	}

	for _, firstListEntry := range firstList {
		for _, secondListEntry := range secondList {
			results = append(results, firstListEntry+secondListEntry)
		}
	}

	return results
}
