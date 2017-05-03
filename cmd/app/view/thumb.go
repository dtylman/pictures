package view

import (
	"fmt"
	"github.com/dtylman/gowd"
	"github.com/dtylman/gowd/bootstrap"
	"github.com/dtylman/pictures/cmd/app/view/darktheme"
	"github.com/dtylman/pictures/indexer/thumbs"
	"github.com/dtylman/pictures/model"
	"log"
)

type thumb struct {
	*gowd.Element
	thumbs     *gowd.Element
	pagination *bootstrap.Pagination
}

func newThumbView() *thumb {
	t := new(thumb)
	t.Element = gowd.NewElement("div")
	t.thumbs = gowd.NewElement("div")
	t.pagination = bootstrap.NewPagination()
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
	for i, image := range activeSearch.Page() {
		thumbPath, err := thumbs.MakeThumb(image.Path, image.MD5, false)
		if err != nil {
			log.Println(err)
		}
		img := bootstrap.NewElement("img", "thumbnail img-responsive")
		img.SetAttribute("src", "file:///"+thumbPath)
		img.OnEvent(gowd.OnClick, t.thumbClick)
		img.Object = i
		span := gowd.NewElement("span")
		span.AddElement(gowd.NewText(fmt.Sprintf("%s (%s)", image.Name(), image.Album)))
		row.AddElement(bootstrap.NewColumn(bootstrap.ColumnMedium, 3, bootstrap.NewElement("div", "well", img, span)))
		if i%4 == 3 {
			t.thumbs.AddElement(row)
			row = bootstrap.NewElement("div", "row")
		}
	}

	//build the pagination
	t.pagination.Items.RemoveElements()

	btn := t.pagination.AddItem("", false, t.btnPrevClick)
	icon := gowd.NewElement("i")
	icon.SetClass("fa fa-arrow-left")
	btn.AddElement(icon)

	pages := activeSearch.Pages()
	activePage := pages.ActivePage()
	for pageOrder, page := range pages {
		if pageOrder > (activePage-7) && pageOrder < (activePage+7) {
			item := t.pagination.AddItem(page.Caption, page.Active, t.btnPageClick)
			item.Object = page
		}

	}

	btn = t.pagination.AddItem("", false, t.btnNextClick)
	icon = gowd.NewElement("i")
	icon.SetClass("fa fa-arrow-right")
	btn.AddElement(icon)

	//// facets
	//t.facets.RemoveElements()
	//for _, facet := range activeSearch.Facets {
	//	t.facets.AddOption(fmt.Sprintf("%s (%d)", facet.Term, facet.Count), facet.Term)
	//}
}

func (t *thumb) populateToolbar(menu *darktheme.Menu) {
	if activeSearch != nil {
		menu.TopLeft.AddElement(t.pagination.Element)
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
