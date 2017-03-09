package controller

import (
	"github.com/blevesearch/bleve"
	"github.com/dtylman/pictures/indexer/db"
	"github.com/dtylman/pictures/server/view"
	"net/http"
)

// Search displays the search page
func Search(w http.ResponseWriter, r *http.Request) {
	// Display the view
	v := view.New(r)
	v.Name = "index/search"

	query := r.FormValue("query")
	v.Vars["query"] = query
	var sr *bleve.SearchResult
	var err error
	if query != "" {
		sr, err = db.Search(bleve.NewSearchRequest(bleve.NewQueryStringQuery(query)))
	} else {
		sr, err = db.Search(bleve.NewSearchRequest(bleve.NewMatchAllQuery()))
	}
	if err != nil {
		flash(r, view.FlashError, err.Error())
	}
	v.Vars["hits"] = sr.Hits
	v.Vars["total"] = sr.Total
	v.Render(w)
}
