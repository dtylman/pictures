package webkit

import (
	"fmt"
	"io"
)

const (
	AttrClass = "class"
	AttrID    = "id"
)

type ObjectAttributes map[string]string

func NewAttributes() ObjectAttributes {
	return make(ObjectAttributes)
}

func (oa ObjectAttributes) Render(writer io.Writer) error {
	for name, value := range oa {
		_, err := fmt.Fprintf(writer, ` %s="%s"`, name, value)
		if err != nil {
			return err
		}
	}
	return nil
}
