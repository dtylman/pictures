package view

import (
	"github.com/blevesearch/bleve"
	"github.com/dtylman/gowd"
	"github.com/dtylman/gowd/bootstrap"
	"github.com/dtylman/pictures/server/model"
)

type search struct {
	*gowd.Element
	inputSearch  *bootstrap.Input
	btnSearch    *gowd.Element
	albums       *gowd.Element
	pagination   *bootstrap.Pagination
	panelSearch  *gowd.Element
	activeSearch *model.Search
}

func newSearch() *search {
	s := new(search)
	s.Element = bootstrap.NewContainer(true)

	s.inputSearch = bootstrap.NewInput(bootstrap.InputTypeText, "Search")
	s.inputSearch.SetAttribute("placeholder", "Search...")
	s.AddElement(s.inputSearch.Element)

	s.btnSearch = bootstrap.NewButton(bootstrap.ButtonPrimary, "Search")
	s.AddElement(s.btnSearch)

	s.panelSearch = bootstrap.NewRow()
	s.AddElement(s.panelSearch)

	s.populateSearch()

	return s
}

func (s *search) populateSearch() {
	s.panelSearch.RemoveElements()

	if s.activeSearch == nil {
		var err error
		s.activeSearch, err = model.NewSearch("", bleve.NewMatchAllQuery())
		if err != nil {
			Root.addAlertError(err)
			return
		}
	}

	for _, thumb := range s.activeSearch.Thumbs {
		img := bootstrap.NewElement("img", "img-thumbnail")
		img.SetAttribute("src", "file:///"+thumb.Path)
		col := bootstrap.NewColumn(bootstrap.ColumnXtraSmall, 3)
		col.AddElement(img)
		s.panelSearch.AddElement(col)
	}

	//build the pagination
	s.pagination = bootstrap.NewPagination()
	btn := bootstrap.NewLinkButton("<<")
	btn.OnEvent(gowd.OnClick, s.btnPrevClick)
	s.pagination.Items.AddItem(btn)
	for _, page := range s.activeSearch.Pages {
		btn := bootstrap.NewLinkButton(page.Caption)
		btn.Object = page
		s.pagination.Items.AddItem(btn)

	}
	btn = bootstrap.NewLinkButton(">>")
	btn.OnEvent(gowd.OnClick, s.btnNextClick)
	s.pagination.Items.AddItem(btn)
	s.panelSearch.AddElement(s.pagination.Element)
}

func (s *search) btnPrevClick(sender *gowd.Element, event *gowd.EventElement) {
	s.activeSearch.PrevPage()
	s.populateSearch()
}

func (s *search) btnNextClick(sender *gowd.Element,
	event *gowd.EventElement) {
	s.activeSearch.NextPage()
	s.populateSearch()
}
