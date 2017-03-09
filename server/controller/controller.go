package controller

import (
	"fmt"
	"github.com/blevesearch/bleve"
	"github.com/dtylman/pictures/indexer/db"
	"github.com/dtylman/pictures/server/session"
	"github.com/dtylman/pictures/server/view"
	gorillacontext "github.com/gorilla/context"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

// About displays the About page
func About(w http.ResponseWriter, r *http.Request) {
	// Display the view
	v := view.New(r)
	v.Name = "about/about"
	v.Render(w)
}

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
		sr, err = db.Query(bleve.NewTermQuery(query), 0, 25)
	} else {
		sr, err = db.QueryAll(0, 25)
	}
	if err != nil {
		flash(r, view.FlashError, err.Error())
	}
	v.Vars["hits"] = sr.Hits
	v.Vars["total"] = sr.Total
	v.Render(w)
}

// Backup displays the backup page
func Backup(w http.ResponseWriter, r *http.Request) {
	// Display the view
	v := view.New(r)
	v.Name = "backup/backup"
	v.Render(w)
}
func flash(r *http.Request, messageType string, message string, args ...interface{}) {
	session.Instance(r).AddFlash(view.Flash{fmt.Sprintf(message, args...), messageType})
}

func flashError(r *http.Request, err error) {
	flash(r, view.FlashError, err.Error())
}

func getRouterParam(r *http.Request, name string) string {
	var params httprouter.Params
	params = gorillacontext.Get(r, "params").(httprouter.Params)
	return params.ByName(name)
}

func isChecked(r *http.Request, name string) bool {
	return r.FormValue(name) == "on"
}
