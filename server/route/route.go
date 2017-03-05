package route

import (
	"net/http"

	"github.com/dtylman/pictures/server/controller"
	hr "github.com/dtylman/pictures/server/route/middleware/httprouterwrapper"
	"github.com/gorilla/context"
	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
)

// LoadHTTP returns the HTTPS routes and middleware
func LoadHTTP() http.Handler {
	return middleware(routes())
}

func routes() *httprouter.Router {
	r := httprouter.New()

	// Set 404 handler
	r.NotFound = alice.
		New().
		ThenFunc(controller.Error404)

	// Serve static files, no directory browsing
	r.GET("/static/*filepath", hr.Handler(alice.New().ThenFunc(controller.Static)))

	r.GET("/", hr.Handler(alice.New().ThenFunc(controller.Search)))
	r.POST("/", hr.Handler(alice.New().ThenFunc(controller.Search)))
	r.GET("/index", hr.Handler(alice.New().ThenFunc(controller.Index)))
	r.GET("/index/start", hr.Handler(alice.New().ThenFunc(controller.IndexStart)))
	r.GET("/index/stop", hr.Handler(alice.New().ThenFunc(controller.IndexStop)))
	r.GET("/backup", hr.Handler(alice.New().ThenFunc(controller.Backup)))
	r.GET("/about", hr.Handler(alice.New().ThenFunc(controller.About)))
	r.GET("/settings", hr.Handler(alice.New().ThenFunc(controller.Settings)))
	r.POST("/settings", hr.Handler(alice.New().ThenFunc(controller.Settings)))
	r.POST("/settings/removesourcefolder/:folder", hr.Handler(alice.New().ThenFunc(controller.RemoveSourceFolder)))
	return r
}

func middleware(h http.Handler) http.Handler {
	// Clear handler for Gorilla Context
	h = context.ClearHandler(h)
	return h
}
