package picture

import (
	"github.com/stretchr/testify/assert"
	"os"
	"path/filepath"
	"testing"
)

func TestQueue_Items(t *testing.T) {
	fileName, err := filepath.Abs(filepath.Join("_testdata", "legs.jpg"))
	if err != nil {
		t.Fatal(err)
	}
	fileInfo, err := os.Stat(fileName)
	if err != nil {
		t.Fatal(err)
	}
	pic, err := NewIndex(fileName, fileInfo)
	if err != nil {
		t.Fatal(err)
	}

	q := NewQueue()
	assert.Equal(t, 0, q.Length())
	q.PushBack(pic)
	assert.Equal(t, 1, q.Length())
	keys := q.Keys()
	assert.Equal(t, keys[0], pic.MD5)
	vals := q.Items()
	assert.Equal(t, vals[0].MD5, pic.MD5)
	len, pic1 := q.Pop()
	assert.Equal(t, 0, q.Length())
	assert.Equal(t, 1, len)
	assert.EqualValues(t, pic, pic1)

}
