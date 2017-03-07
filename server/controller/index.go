package controller

import (
	"github.com/dtylman/pictures/indexer"
	"github.com/dtylman/pictures/server/view"
	"net/http"
)

// Index displays the index page
func Index(w http.ResponseWriter, r *http.Request) {
	// Display the view
	v := view.New(r)
	v.Name = "index/index"
	if r.Method == http.MethodPost {
		action := r.FormValue("action")
		if action == "start" {
			err := indexer.Start(indexer.Options{IndexLocation: false, ReIndex: false})
			if err != nil {
				flashError(r, err)
			}
		} else if action == "stop" {
			indexer.Stop()
		}
	}

	v.Vars["index_running"] = indexer.IsRunning()
	v.Vars["index_progress"] = indexer.GetProgress()

	v.Render(w)
}
