package lockfile

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestOpen(t *testing.T) {
	i, err := Open()
	assert.Nil(t, i)
	assert.True(t, os.IsNotExist(err))
	i = &Info{Address: "127.0.0.1:8080"}
	err = i.Create()
	defer Delete()
	assert.NoError(t, err)
	i, err = Open()
	assert.NoError(t, err)
	assert.Equal(t, i.Address, "127.0.0.1:8080")
}
