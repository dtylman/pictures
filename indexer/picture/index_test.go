package picture

import (
	"github.com/stretchr/testify/assert"
	"os"
	"path/filepath"
	"testing"
)

func TestNewIndex(t *testing.T) {
	fileName, err := filepath.Abs(filepath.Join("_testdata", "legs.jpg"))
	if err != nil {
		t.Fatal(err)
	}
	fileInfo, err := os.Stat(fileName)
	if err != nil {
		t.Fatal(err)
	}
	i, err := NewIndex(fileName, fileInfo)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, "_testdata", i.Album)
}
