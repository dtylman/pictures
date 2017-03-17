package controller

import (
	"fmt"
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
			err := indexer.Start(
				indexer.Options{IndexLocation: isChecked(r, "location"), ReIndex: isChecked(r, "reindex")})
			if err != nil {
				flashError(r, err)
			}
		} else if action == "stop" {
			indexer.Stop()
		}
	}

	v.Vars["index_running"] = indexer.IsRunning()
	v.Vars["reindex"] = r.FormValue("reindex")
	v.Vars["location"] = r.FormValue("location")
	v.Render(w)
}

func IndexStatus(w http.ResponseWriter, r *http.Request) {
	obj, err := indexer.GetProgressStatus().ToJSON()
	if err != nil {
		Error500(w, r, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, obj)
}
