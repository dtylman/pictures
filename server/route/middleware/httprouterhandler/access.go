package httprouterhandler

import (
	"sync"
	"time"
)

type accessTime struct {
	time  time.Time
	mutex sync.Mutex
}

func init() {
	LastAccess.update()
}

func (at *accessTime) update() {
	at.mutex.Lock()
	defer at.mutex.Unlock()
	at.time = time.Now()
}

//Elapsed returns true if time since access was larger than duration
func (at *accessTime) Elapsed(d time.Duration) bool {
	at.mutex.Lock()
	defer at.mutex.Unlock()
	return time.Now().Sub(at.time) > d
}
