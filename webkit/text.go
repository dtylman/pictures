package webkit

import (
	"fmt"
	"io"
)

type Text string

func (t Text) Render(writer io.Writer) error {
	_, err := fmt.Fprint(writer, t)
	return err
}
