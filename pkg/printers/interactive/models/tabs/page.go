package tabs

type page struct {
	startIndex int
	endIndex   int
}

type joinFunc func(...string) string

type boundaryFunc func(string) bool

type renderFunc func(int, string) string

func pagesForTabs(tabs []Tab, join joinFunc, boundary boundaryFunc, render renderFunc) []*page {
	pages := []*page{}
	tempPage := &page{startIndex: 0}
	tabStr := ""
	for i, tab := range tabs {
		rendered := render(i, tab.Name)
		tempTabStr := join(tabStr, rendered)

		if boundary(tempTabStr) {
			tempPage.endIndex = i - 1
			pages = append(pages, tempPage)
			tempPage = &page{startIndex: i}
			tabStr = join("", rendered)
			continue
		}

		tabStr = tempTabStr
		tempPage.endIndex = i
	}

	// capture any lingering page
	pages = append(pages, tempPage)

	return pages
}
