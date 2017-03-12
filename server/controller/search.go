package controller

import (
	"github.com/blevesearch/bleve"
	"github.com/blevesearch/bleve/search"
	"github.com/blevesearch/bleve/search/query"
	"github.com/dtylman/pictures/indexer/db"
	"github.com/dtylman/pictures/server/view"
	"net/http"
	"time"
)

//FacetItem represents facet item in the display
type FacetItem struct {
	Name  string
	Field string
	Term  string
	Count int
}

// Search displays the search page
func Search(w http.ResponseWriter, r *http.Request) {
	// Display the view
	v := view.New(r)
	v.Name = "index/search"

	queryString := r.FormValue("query")
	v.Vars["query"] = queryString
	var query query.Query
	if queryString != "" {
		query = bleve.NewQueryStringQuery(queryString)
	} else {
		query = bleve.NewMatchAllQuery()
	}
	req := bleve.NewSearchRequest(query)
	req.AddFacet("Location", bleve.NewFacetRequest("location", 4))
	takenFacet := bleve.NewFacetRequest("taken", 3)
	takenFacet.AddDateTimeRange("Old", time.Unix(0, 0), time.Now().AddDate(0, -1, 0))
	takenFacet.AddDateTimeRange("Last Month", time.Now().AddDate(0, -1, 0), time.Now())
	takenFacet.AddDateTimeRange("Last Week", time.Now().AddDate(0, 0, -7), time.Now())
	req.AddFacet("Taken", takenFacet)
	sr, err := db.Search(req)
	if err != nil {
		flashError(r, err)
	} else {
		v.Vars["facets"] = facetItems(sr.Facets)
		v.Vars["hits"] = sr.Hits
		v.Vars["total"] = sr.Total
	}

	v.Render(w)
}

func facetItems(res search.FacetResults) []FacetItem {
	items := make([]FacetItem, 0)
	for fn, fr := range res {
		for _, term := range fr.Terms {
			items = append(items, FacetItem{Name: fn, Field: fr.Field, Term: term.Term, Count: term.Count})
		}
	}
	return items
}
