package view

import (
	"github.com/dtylman/gowd"
	"github.com/dtylman/gowd/bootstrap"
	"github.com/dtylman/pictures/cmd/app/view/darktheme"
	"github.com/dtylman/pictures/model"
	"path/filepath"
)

type thumb struct {
	*gowd.Element
	thumbs *gowd.Element
}

func newThumbView() *thumb {
	t := new(thumb)
	t.Element = gowd.NewElement("div")
	t.thumbs = gowd.NewElement("div")
	t.AddElement(t.thumbs)
	return t
}

func (t *thumb) updateState() {
	if activeSearch == nil {
		var err error
		activeSearch, err = model.NewSearch("")
		if err != nil {
			Root.addAlertError(err)
			return
		}
	}

	/*<div class="col-md-3">
	    <div class="well">
	        <img class="thumbnail img-responsive" alt="Bootstrap template" src="http://www.prepbootstrap.com/Content/images/shared/houses/h9.jpg"
	            onclick="alert('lala')" />
	        <span>
	            jpg.jpg, London,  dog with and Yogev
	        </span>
	    </div>
	</div>*/
	t.thumbs.RemoveElements()
	row := bootstrap.NewElement("div", "row")
	for i, thumb := range activeSearch.Thumbs {
		img := bootstrap.NewElement("img", "thumbnail img-responsive")
		img.SetAttribute("src", "file:///"+thumb.Path)
		img.OnEvent(gowd.OnClick, t.thumbClick)
		img.Object = i
		span := gowd.NewElement("span")
		span.AddElement(gowd.NewText(filepath.Base(thumb.Path)))
		row.AddElement(bootstrap.NewColumn(bootstrap.ColumnMedium, 3, bootstrap.NewElement("div", "well", img, span)))
		if i%4 == 3 {
			t.thumbs.AddElement(row)
			row = bootstrap.NewElement("div", "row")
		}
	}

	//// facets
	//t.facets.RemoveElements()
	//for _, facet := range activeSearch.Facets {
	//	t.facets.AddOption(fmt.Sprintf("%s (%d)", facet.Term, facet.Count), facet.Term)
	//}
}

func (t *thumb) populateToolbar(menu *darktheme.Menu) {
	//build the pagination
	menu.AddButton(menu.TopLeft, "", "fa fa-arrow-left", t.btnPrevClick)
	if activeSearch != nil {
		activePage := activeSearch.Pages.ActivePage()
		for pageOrder, page := range activeSearch.Pages {
			if pageOrder > (activePage-7) && pageOrder < (activePage+7) {
				btn := menu.AddButton(menu.TopLeft, page.Caption, "", t.btnPageClick)
				btn.Object = page
				if page.Active {
					btn.SetClass("active")
				}
			}

		}
		menu.AddButton(menu.TopLeft, "", "fa fa-arrow-right", t.btnNextClick)
	}
}

func (t *thumb) getContent() *gowd.Element {
	return t.Element
}

func (t *thumb) btnPageClick(sender *gowd.Element, event *gowd.EventElement) {
	page := sender.Object.(model.PageItem)
	activeSearch.StartFrom(page.Start)
	t.updateState()
}

func (t *thumb) btnPrevClick(sender *gowd.Element, event *gowd.EventElement) {
	activeSearch.PrevPage()
	t.updateState()
}

func (t *thumb) btnNextClick(sender *gowd.Element, event *gowd.EventElement) {
	activeSearch.NextPage()
	t.updateState()
}

func (t *thumb) thumbClick(sender *gowd.Element, event *gowd.EventElement) {
	hit := sender.Object.(int)
	activeSearch.SetActiveImage(hit)
	Root.setActiveView(newImage())
}
