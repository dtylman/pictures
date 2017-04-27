package model

//PageItem represents a paging item
type PageItem struct {
	Start   int
	Active  bool
	Caption string
}

type Pages []PageItem

func (p Pages) ActivePage() int {
	for i := range p {
		if p[i].Active {
			return i
		}
	}
	return -1
}
