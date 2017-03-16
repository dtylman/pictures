package httprouterhandler

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestAccessTime_Elapsed(t *testing.T) {
	LastAccess.update()
	time.Sleep(time.Second)
	assert.True(t, LastAccess.Elapsed(time.Millisecond*500))
	assert.False(t, LastAccess.Elapsed(time.Second*2))
}
