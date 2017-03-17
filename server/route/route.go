package route

import (
	"net/http"

	"github.com/dtylman/pictures/server/controller"
	hr "github.com/dtylman/pictures/server/route/middleware/httprouterhandler"
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
	r.NotFound = alice.New().ThenFunc(controller.Error404)

	// Serve static files, no directory browsing
	r.GET("/static/*filepath", hr.Handler(alice.New().ThenFunc(controller.Static)))
	r.GET("/activeimage/:md5", hr.Handler(alice.New().ThenFunc(controller.ActiveImage)))
	r.GET("/nextimage", hr.Handler(alice.New().ThenFunc(controller.NextImage)))
	r.GET("/previmage", hr.Handler(alice.New().ThenFunc(controller.PrevImage)))
	r.GET("/thumb/:hit/:md5", hr.Handler(alice.New().ThenFunc(controller.Thumb)))
	r.GET("/image/:hit/:md5/view", hr.Handler(alice.New().ThenFunc(controller.ImageView)))
	r.GET("/", hr.Handler(alice.New().ThenFunc(controller.SearchResults)))
	r.POST("/search", hr.Handler(alice.New().ThenFunc(controller.Search)))
	r.GET("/search", hr.Handler(alice.New().ThenFunc(controller.Search)))
	r.GET("/page", hr.Handler(alice.New().ThenFunc(controller.Page)))
	r.GET("/nextpage", hr.Handler(alice.New().ThenFunc(controller.NextPage)))
	r.GET("/prevpage", hr.Handler(alice.New().ThenFunc(controller.PrevPage)))
	r.GET("/index", hr.Handler(alice.New().ThenFunc(controller.Index)))
	r.GET("/index/status", hr.Handler(alice.New().ThenFunc(controller.IndexStatus)))
	r.POST("/index", hr.Handler(alice.New().ThenFunc(controller.Index)))
	r.GET("/backup", hr.Handler(alice.New().ThenFunc(controller.Backup)))
	r.GET("/about", hr.Handler(alice.New().ThenFunc(controller.About)))
	r.GET("/settings", hr.Handler(alice.New().ThenFunc(controller.Settings)))
	r.GET("/backup/status", hr.Handler(alice.New().ThenFunc(controller.BackupStatus)))

	r.POST("/settings", hr.Handler(alice.New().ThenFunc(controller.Settings)))
	r.POST("/settings/removesourcefolder/:folder", hr.Handler(alice.New().ThenFunc(controller.RemoveSourceFolder)))
	return r
}

func middleware(h http.Handler) http.Handler {
	// Clear handler for Gorilla Context
	h = context.ClearHandler(h)
	return h
}
