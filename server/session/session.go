package session

import (
	"net/http"

	"github.com/gorilla/sessions"
	"log"
)

var (
	// Store is the cookie store
	store *sessions.CookieStore
	// Name is the session name
	name string
	//Cache holds all sessions server side state
	Cache State
)

// Session stores session level information
type Session struct {
	Options   sessions.Options `json:"Options"`   // Pulled from: http://www.gorillatoolkit.org/pkg/sessions#Options
	Name      string           `json:"Name"`      // Name for: http://www.gorillatoolkit.org/pkg/sessions#CookieStore.Get
	SecretKey string           `json:"SecretKey"` // Key for: http://www.gorillatoolkit.org/pkg/sessions#CookieStore.New
}

// Configure the session cookie store
func Configure(s Session) {
	store = sessions.NewCookieStore([]byte(s.SecretKey))
	store.Options = &s.Options
	name = s.Name
	Cache.init()
}

// Instance returns a session
func Instance(r *http.Request) *sessions.Session {
	sess, err := store.Get(r, name)
	if err != nil {
		log.Fatal(err)
		return nil
	}
	if sess.IsNew {
		Cache.addSession(sess)
	}
	return sess
}

// Empty deletes all the current session values
func Empty(sess *sessions.Session) {
	// Clear out all stored values in the cookie
	for k := range sess.Values {
		delete(sess.Values, k)
	}
	Cache.deleteSession(sess)

}
