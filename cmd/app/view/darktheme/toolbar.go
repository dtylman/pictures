package darktheme

import (
	"github.com/dtylman/gowd"
	"github.com/dtylman/gowd/bootstrap"
)

type Toolbar struct {
	*gowd.Element
}

/*<div class="row">
                <div class="col-lg-4">
                    <p>
                        <input class="form-control" placeholder="Enter text">
                        </p>
                    </div>
                 <div class="col-lg-8">

              <p>

                <button type="button" class="btn btn-primary">Primary</button>
                <button type="button" class="btn btn-success">Success</button>
                <button type="button" class="btn btn-info">Info</button>
                <button type="button" class="btn btn-warning">Warning</button>
                <button type="button" class="btn btn-danger">Danger</button>
                <button type="button" class="btn btn-link">Link</button>
              </p>
                 </div>
            <div>*/

func NewToolbar() *Toolbar {
	t := new(Toolbar)
	t.Element = bootstrap.NewRow()
	return t;
}