package controller

import (
	"fmt"
	"github.com/dtylman/pictures/indexer"
	"github.com/dtylman/pictures/indexer/runningindexer"
	"github.com/dtylman/pictures/server/session"
	"github.com/dtylman/pictures/server/view"
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
	if r.Method == http.MethodPost {
		v.Vars["first_name"] = r.FormValue("text")
	} else {
		v.Vars["first_name"] = "Bart Simpson"
	}

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
