package controller

import (
	"fmt"
	"github.com/dtylman/pictures/server/session"
	"github.com/dtylman/pictures/server/view"
	gorillacontext "github.com/gorilla/context"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
)

// About displays the About page
func About(w http.ResponseWriter, r *http.Request) {
	// Display the view
	v := view.New(r)
	v.Name = "about/about"
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
	log.Println(name, r.FormValue(name))
	return r.FormValue(name) == "on"
}
