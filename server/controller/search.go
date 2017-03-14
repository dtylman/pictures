package controller

import (
	"github.com/blevesearch/bleve"
	"github.com/blevesearch/bleve/search/query"
	"github.com/dtylman/pictures/server/model"
	"github.com/dtylman/pictures/server/view"
	"net/http"
	"strconv"
)

var mySearch *model.Search

//Search performs new search
func Search(w http.ResponseWriter, r *http.Request) {
	queryString := r.FormValue("query")
	var query query.Query
	if queryString != "" {
		query = bleve.NewQueryStringQuery(queryString)
	} else {
		query = bleve.NewMatchAllQuery()
	}
	var err error
	mySearch, err = model.NewSearch(queryString, query)
	if err != nil {
		flashError(r, err)
	}
	SearchResults(w, r)
}

func Page(w http.ResponseWriter, r *http.Request) {
	fromString := r.URL.Query().Get("from")
	var err error
	if fromString == "" {
		err = mySearch.StartFrom(0)
	} else {
		from, err := strconv.Atoi(fromString)
		if err == nil {
			err = mySearch.StartFrom(from)
		}
	}
	if err != nil {
		flashError(r, err)
	}
	SearchResults(w, r)
}

// Search displays the search page
func SearchResults(w http.ResponseWriter, r *http.Request) {
	v := view.New(r)
	v.Name = "index/search"
	if mySearch != nil {
		v.Vars["query"] = mySearch.QueryString
		v.Vars["facets"] = mySearch.Facets
		v.Vars["thumbs"] = mySearch.Thumbs
		v.Vars["total"] = mySearch.Result.Total
		v.Vars["pages"] = mySearch.Pages
	}
	v.Render(w)
}
