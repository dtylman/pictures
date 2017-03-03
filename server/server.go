package server

import (
	"fmt"
	"github.com/dtylman/pictures/server/session"
	"github.com/dtylman/pictures/server/view"
	"github.com/dtylman/pictures/server/view/plugin"
	"log"
	"net"
	"net/http"
	"os/exec"
	"time"
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
		plugin.PrettyTime())
}

// Run starts the HTTP and/or HTTPS listener
func Start(httpHandlers http.Handler) error {
	address, err := getAddress()
	if err != nil {
		return err
	}
	go startServer(address, httpHandlers)
	time.Sleep(time.Second)
	cmd := exec.Command("google-chrome", fmt.Sprintf("http://%s", address))
	return cmd.Run()

}

func getAddress() (string, error) {
	addr, err := net.ResolveTCPAddr("tcp", "127.0.0.1:0")
	if err != nil {
		return "", err
	}
	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return "", err
	}
	defer l.Close()
	return l.Addr().String(), nil
}

func startServer(address string, httpHandlers http.Handler) {
	err := http.ListenAndServe(address, httpHandlers)
	if err != nil {
		log.Fatal(err)
	}
}
