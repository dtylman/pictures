package httprouterhandler

import (
	"net/http"

	"github.com/gorilla/context"
	"github.com/julienschmidt/httprouter"
)

//LastAccess stores when was the last time server accessed
var LastAccess accessTime

// Handler accepts a handler to make it compatible with http.HandlerFunc
func Handler(h http.Handler) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		//update the last access time
		LastAccess.update()
		context.Set(r, "params", p)
		h.ServeHTTP(w, r)
	}
}
