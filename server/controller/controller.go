package controller

import (
	"fmt"
	"github.com/blevesearch/bleve"
	"github.com/dtylman/pictures/db"
	"github.com/dtylman/pictures/indexer"
	"github.com/dtylman/pictures/indexer/runningindexer"
	"github.com/dtylman/pictures/server/session"
	"github.com/dtylman/pictures/server/view"
	goriilacontext "github.com/gorilla/context"
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

// Index displays the index page
func Index(w http.ResponseWriter, r *http.Request) {
	// Display the view
	v := view.New(r)
	v.Name = "index/index"
	v.Vars["index_running"] = runningindexer.IsRunning()

	v.Render(w)
}

// Index displays the index page
func IndexStart(w http.ResponseWriter, r *http.Request) {
	// Display the view
	v := view.New(r)
	v.Name = "index/index"
	err := indexer.Start(runningindexer.Options{IndexLocation: false, ReIndex: false})
	if err != nil {
		flash(r, view.FlashError, err.Error())
	}
	v.Vars["index_running"] = runningindexer.IsRunning()

	v.Render(w)
}

// Index displays the index page
func IndexStop(w http.ResponseWriter, r *http.Request) {
	// Display the view
	v := view.New(r)
	v.Name = "index/index"
	v.Vars["index_running"] = runningindexer.IsRunning()

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

func getParamByName(r *http.Request, name string) string {
	var params httprouter.Params
	params = goriilacontext.Get(r, "params").(httprouter.Params)
	return params.ByName(name)
}
