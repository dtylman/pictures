package controller

import (
	"fmt"
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
		v.Vars["first_name"] = fmt.Sprintf("%v", r.FormValue("text"))
	} else {
		v.Vars["first_name"] = "Bart Simpson"
	}

	v.Render(w)
}

// About displays the About page
func Settings(w http.ResponseWriter, r *http.Request) {
	// Display the view
	v := view.New(r)
	v.Name = "settings/settings"
	v.Render(w)
}
