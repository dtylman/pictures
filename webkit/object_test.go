package webkit

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRender(t *testing.T) {
	c := NewObject("div", "lala", "123")
	var output bytes.Buffer
	err := c.Render(&output)
	assert.NoError(t, err)
	assert.Equal(t, output.String(), `<div class="lala" id="123"></div>`)
	c.AddChild(NewObject("span", "nana", "123"))
	output.Reset()
	err = c.Render(&output)
	assert.NoError(t, err)
	assert.Equal(t, output.String(), `<div class="lala" id="123"><span class="nana" id="123"></span></div>`)
	c.AddChild(Text("hoho"))
	output.Reset()
	c.Render(&output)
	t.Log(output.String())
}
