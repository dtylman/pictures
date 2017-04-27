package view

import (
	"fmt"
	"github.com/dtylman/gowd"
	"github.com/dtylman/gowd/bootstrap"
	"github.com/dtylman/pictures/model"
)

var activeSearch *model.Search

type search struct {
	*gowd.Element
	inputSearch *gowd.Element
	btnSearch   *gowd.Element
	facets      *Select
}

func newSearchView() *search {
	s := new(search)
	var err error
	s.Element, err = gowd.ParseElement(`<div>
	    	<div class="row">
			<div class="col-lg-4">
		    		<p>
				<input class="form-control" placeholder="Search String">
		    		</p>
			</div>
		<div class="col-lg-1">
		    <p>
			<button type="button" class="btn btn-primary">Search</button>
		    </p>
		</div>
		<div class="col-lg-7">
		    <p>
			<select class="form-control"></select>
		    </p>
		</div>
	    </div></div>`)
	if err != nil {
		panic(err)
	}
	//
	s.inputSearch = bootstrap.NewInput(bootstrap.InputTypeText)
	s.facets = NewSelect()
	s.inputSearch.SetAttribute("placeholder", "Search...")
	s.inputSearch.SetClass("form-control")
	s.inputSearch.OnKeyPressEvent(gowd.OnKeyPress, 13, s.btnSearchClick)
	s.inputSearch.SetAttribute("autofocus", "true")
	s.btnSearch = bootstrap.NewButton(bootstrap.ButtonPrimary, "Search")
	s.btnSearch.OnEvent(gowd.OnClick, s.btnSearchClick)
	s.AddElement(s.btnSearch)
	s.facets.Element.SetClass("form-control")
	s.facets.OnEvent(gowd.OnChange, s.facetChanged)

	return s
}

func (s *search) facetChanged(sender *gowd.Element, event *gowd.EventElement) {
	term := event.GetValue()
	s.inputSearch.SetAttribute("value", term)
	s.doQuery(term)
}

func (s *search) btnPageClick(sender *gowd.Element, event *gowd.EventElement) {
	page := sender.Object.(model.PageItem)
	activeSearch.StartFrom(page.Start)
	s.updateState()
}

func (s *search) btnPrevClick(sender *gowd.Element, event *gowd.EventElement) {
	activeSearch.PrevPage()
	s.updateState()
}

func (s *search) btnNextClick(sender *gowd.Element, event *gowd.EventElement) {
	activeSearch.NextPage()
	s.updateState()
}

func (s *search) btnSearchClick(sender *gowd.Element, event *gowd.EventElement) {
	term := s.inputSearch.GetValue()
	s.doQuery(term)
}

func (s *search) doQuery(term string) {
	var err error
	activeSearch, err = model.NewSearch(term)
	if err != nil {
		Root.addAlertError(err)
		return
	}
	s.updateState()
}

func (s *search) thumbClick(sender *gowd.Element, event *gowd.EventElement) {
	hit := sender.Object.(int)
	activeSearch.SetActiveImage(hit)
	Root.setActiveView(newImage())
}

func (s *search) updateState() {

	well := bootstrap.NewElement("div", "well")
	row := bootstrap.NewRow()
	well.AddElement(row)

	if activeSearch == nil {
		var err error
		activeSearch, err = model.NewSearch("")
		if err != nil {
			Root.addAlertError(err)
			return
		}
	}

	for i, thumb := range activeSearch.Thumbs {
		img := bootstrap.NewElement("img", "img-thumbnail")
		img.SetAttribute("src", "file:///" + thumb.Path)

		link := bootstrap.NewLinkButton("")
		link.Object = i
		link.OnEvent(gowd.OnClick, s.thumbClick)
		link.AddElement(img)

		col := bootstrap.NewColumn(bootstrap.ColumnLarge, 3)
		col.AddElement(link)

		row.AddElement(col)
	}

	//build the pagination
	pagination := bootstrap.NewPagination()
	btn := bootstrap.NewLinkButton("<<")
	btn.OnEvent(gowd.OnClick, s.btnPrevClick)
	pagination.Items.AddItem(btn)
	activePage := activeSearch.Pages.ActivePage()
	for pageOrder, page := range activeSearch.Pages {
		if pageOrder > (activePage - 7) && pageOrder < (activePage + 7) {
			btn := bootstrap.NewLinkButton(page.Caption)
			btn.Object = page
			btn.OnEvent(gowd.OnClick, s.btnPageClick)
			item := pagination.Items.AddItem(btn)
			if page.Active {
				item.SetClass("active")
			}
		}

	}
	btn = bootstrap.NewLinkButton(">>")
	btn.OnEvent(gowd.OnClick, s.btnNextClick)
	pagination.Items.AddItem(btn)

	// facets
	s.facets.RemoveElements()
	for _, facet := range activeSearch.Facets {
		s.facets.AddOption(fmt.Sprintf("%s (%d)", facet.Term, facet.Count), facet.Term)
	}
	s.Element.AddElement(well)
}

func (s *search) getContent() *gowd.Element {
	return s.Element
}
