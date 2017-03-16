package server

import (
	"github.com/dtylman/pictures/server/session"
	"github.com/dtylman/pictures/server/view"
	"github.com/dtylman/pictures/server/view/plugin"
	"log"
	"net/http"
)

func init() {
	s := session.Session{}
	s.SecretKey = "@r4B?EThaSEh_drudR7P_hub=s#s2Pah"
	s.Name = "gosess"
	s.Options.Path = "/"
	s.Options.Domain = ""
	s.Options.MaxAge = 28800
	s.Options.Secure = false
	s.Options.HttpOnly = true
	session.Configure(s)
	v := view.View{
		BaseURI:   "/",
		Caching:   true,
		Folder:    "template",
		Extension: "tmpl",
		Name:      "blank",
	}
	view.Configure(v)
	view.LoadTemplates("base", []string{"partial/menu", "partial/footer"})
	view.LoadPlugins(
		plugin.TagHelper(v),
		plugin.NoEscape(),
		plugin.PrettyTime(),
		plugin.Base64(),
	)
}

//Start starts the server
func Start(address string, httpHandlers http.Handler) error {
	log.Printf("Starting server, browse here: http://%s", address)
	return http.ListenAndServe(address, httpHandlers)
}
