package model

//FacetItem represents facet item in the display
type FacetItem struct {
	Term  string
	Count int
}

type Facets []FacetItem

// Len is the number of elements in the collection.
func (f Facets) Len() int {
	return len(f)
}

// Less reports whether the element with
// index i should sort before the element with index j.
func (f Facets) Less(i, j int) bool {
	return f[i].Count > f[j].Count
}

// Swap swaps the elements with indexes i and j.
func (f Facets) Swap(i, j int) {
	f[i], f[j] = f[j], f[i]
}
