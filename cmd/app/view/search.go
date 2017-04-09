package view

import (
	"fmt"
	"github.com/blevesearch/bleve"
	"github.com/blevesearch/bleve/search/query"
	"github.com/dtylman/gowd"
	"github.com/dtylman/gowd/bootstrap"
	"github.com/dtylman/pictures/model"
)

type search struct {
	*gowd.Element
	inputSearch  *gowd.Element
	btnSearch    *gowd.Element
	facets       *Select
	panelSearch  *gowd.Element
	activeSearch *model.Search
}

func newSearch() *search {
	s := new(search)
	s.Element = bootstrap.NewContainer(true)

	s.panelSearch = bootstrap.NewContainer(true)
	s.AddElement(s.panelSearch)

	s.facets = NewSelect()
	s.facets.Element.SetClass("form-control")
	s.facets.OnEvent(gowd.OnChange, s.facetChanged)

	s.updateState()

	return s
}

func (s *search) facetChanged(sender *gowd.Element, event *gowd.EventElement) {
	term := event.GetValue()
	s.inputSearch.SetAttribute("value", term)
	s.doQuery(term)
}

func (s *search) btnPageClick(sender *gowd.Element, event *gowd.EventElement) {
	page := sender.Object.(model.PageItem)
	err := s.activeSearch.StartFrom(page.From)
	if err != nil {
		Root.addAlertError(err)
		return
	}
	s.updateState()
}

func (s *search) btnPrevClick(sender *gowd.Element, event *gowd.EventElement) {
	s.activeSearch.PrevPage()
	s.updateState()
}

func (s *search) btnNextClick(sender *gowd.Element, event *gowd.EventElement) {
	s.activeSearch.NextPage()
	s.updateState()
}

func (s *search) btnSearchClick(sender *gowd.Element, event *gowd.EventElement) {
	term := s.inputSearch.GetValue()
	s.doQuery(term)
}

func (s *search) doQuery(term string) {
	var query query.Query
	if term == "" {
		query = bleve.NewMatchAllQuery()
	} else {
		query = bleve.NewQueryStringQuery(term)
	}
	var err error
	s.activeSearch, err = model.NewSearch(term, query)
	if err != nil {
		Root.addAlertError(err)
		return
	}
	s.updateState()
}

func (s *search) thumbClick(sender *gowd.Element, event *gowd.EventElement) {
	hit := sender.Object.(int)
	s.activeSearch.SetActiveImage(hit)
	Root.setActiveView(newImage(s.activeSearch))
}

func (s *search) populateToolbar(toolbar *gowd.Element) {
	s.inputSearch = bootstrap.NewInput(bootstrap.InputTypeText)
	s.inputSearch.SetAttribute("placeholder", "Search...")
	s.inputSearch.SetClass("form-control")
	s.inputSearch.OnKeyPressEvent(gowd.OnKeyPress, 13, s.btnSearchClick)
	s.inputSearch.SetAttribute("autofocus", "true")
	s.btnSearch = bootstrap.NewButton(bootstrap.ButtonPrimary, "Search")
	s.btnSearch.OnEvent(gowd.OnClick, s.btnSearchClick)

	toolbar.AddElement(bootstrap.NewInputGroup(s.inputSearch, bootstrap.NewElement("div", "input-group-btn", s.btnSearch), s.facets.Element))

}

func (s *search) updateState() {
	s.panelSearch.RemoveElements()

	well := bootstrap.NewElement("div", "well")
	s.panelSearch.AddElement(well)
	row := bootstrap.NewRow()
	well.AddElement(row)

	if s.activeSearch == nil {
		var err error
		s.activeSearch, err = model.NewSearch("", bleve.NewMatchAllQuery())
		if err != nil {
			Root.addAlertError(err)
			return
		}
	}

	for i, thumb := range s.activeSearch.Thumbs {
		img := bootstrap.NewElement("img", "img-thumbnail")
		img.SetAttribute("src", "file:///"+thumb.Path)

		link := bootstrap.NewLinkButton("")
		link.Object = i
		link.OnEvent(gowd.OnClick, s.thumbClick)
		link.AddElement(img)

		col := bootstrap.NewColumn(bootstrap.ColumnXtraSmall, 3)
		col.AddElement(link)

		row.AddElement(col)
	}

	//build the pagination
	pagination := bootstrap.NewPagination()
	btn := bootstrap.NewLinkButton("<<")
	btn.OnEvent(gowd.OnClick, s.btnPrevClick)
	pagination.Items.AddItem(btn)
	for _, page := range s.activeSearch.Pages {
		btn := bootstrap.NewLinkButton(page.Caption)
		btn.Object = page
		btn.OnEvent(gowd.OnClick, s.btnPageClick)
		item := pagination.Items.AddItem(btn)
		if page.Active {
			item.SetClass("active")
		}

	}
	btn = bootstrap.NewLinkButton(">>")
	btn.OnEvent(gowd.OnClick, s.btnNextClick)
	pagination.Items.AddItem(btn)
	s.panelSearch.AddElement(pagination.Element)

	// facets

	s.facets.RemoveElements()
	for _, facet := range s.activeSearch.Facets {
		s.facets.AddOption(fmt.Sprintf("%s (%d)", facet.Term, facet.Count), fmt.Sprintf("%s: %s", facet.Field, facet.Term))
	}

}

func (s *search) getContent() *gowd.Element {
	return s.Element
}
