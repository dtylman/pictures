package controller

import (
	"fmt"
	"github.com/dtylman/pictures/server/session"
	"github.com/dtylman/pictures/server/view"
	gorillacontext "github.com/gorilla/context"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
	"strconv"
)

// About displays the About page
func About(w http.ResponseWriter, r *http.Request) {
	// Display the view
	v := view.New(r)
	v.Name = "about/about"
	v.Render(w)
}

func flash(r *http.Request, messageType string, message string, args ...interface{}) {
	sess := session.Instance(r)
	sess.AddFlash(view.Flash{fmt.Sprintf(message, args...), messageType})
}

func flashError(r *http.Request, err error) {
	flash(r, view.FlashError, err.Error())
}

func getRouterParam(r *http.Request, name string) string {
	var params httprouter.Params
	params = gorillacontext.Get(r, "params").(httprouter.Params)
	return params.ByName(name)
}

func getRouterParamInt(r *http.Request, name string) (int, error) {
	s := getRouterParam(r, name)
	return strconv.Atoi(s)
}

func isChecked(r *http.Request, name string) bool {
	log.Println(name, r.FormValue(name))
	return r.FormValue(name) == "on"
}
