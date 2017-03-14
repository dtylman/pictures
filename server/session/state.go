package session

import (
	"github.com/gorilla/sessions"
	"github.com/patrickmn/go-cache"
	"sync"
	"time"
)

//State holds cache for all session
type State struct {
	mutex sync.Mutex
	data  map[string]*cache.Cache
}

func (s *State) init() {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.data = make(map[string]*cache.Cache)
}

func (s *State) addSession(sess *sessions.Session) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.data[sess.ID] = cache.New(time.Second*time.Duration(sess.Options.MaxAge), time.Second*30)
}

func (s *State) deleteSession(sess *sessions.Session) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	delete(s.data, sess.ID)
}

func (s *State) Put(sess *sessions.Session, name string, value interface{}) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if s.data[sess.ID] == nil {
		s.data[sess.ID] = cache.New(time.Second*time.Duration(sess.Options.MaxAge), time.Second*30)
	}
	s.data[sess.ID].SetDefault(name, value)
}

func (s *State) Get(sess *sessions.Session, name string) (interface{}, bool) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	c := s.data[sess.ID]
	if c == nil {
		return nil, false
	}
	return c.Get(name)
}
