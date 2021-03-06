package view

import (
	"fmt"

	"github.com/dtylman/gowd"
	"github.com/dtylman/gowd/bootstrap"
	"github.com/dtylman/pictures/cmd/app/view/darktheme"
	"github.com/dtylman/pictures/indexer/db"
	"github.com/dtylman/pictures/model"
	"github.com/dtylman/pictures/cmd/app/view/querybuilder"
)

type search struct {
	*gowd.Element
	inputSearch *gowd.Element
	txtStats    *gowd.Element
}

func newSearchView() *search {
	s := new(search)
	var err error
	em:=gowd.NewElementMap()
	s.Element, err = gowd.ParseElement(`<div>
                <div class="row">
                    <div class="col-lg-12 text-center v-center">
                        <h1>Text Search</h1>
                        <p class="lead" id="pnlSubtitle"></p>
                        <br>
                        <br>
                        <div class="input-group" style="width: 340px; text-align: center; margin: 0 auto;" id="pnlSearch">
                        </div>
                    </div>
                </div>
                <br>
                <div id="pnlAdvanced" />
                <div class="text-center">
                    <h3>Or try one of these:</h3>
                </div>
                <div class="row">
                    <div class="col-lg-12 text-center v-center" style="font-size: 39pt;" id="pnlButtons">

                    </div>
                </div>
            </div>`,em)
	if err != nil {
		panic(err)
	}
	pnlSubtitle := s.Find("pnlSubtitle")
	s.txtStats = gowd.NewText("")
	pnlSubtitle.AddElement(s.txtStats)
	pnlSearch := s.Find("pnlSearch")
	pnlSubtitle := em["pnlSubtitle"]
	stats, err := db.Stats()
	if err != nil {
		panic(err)
	} else {
		pnlSubtitle.AddElement(gowd.NewText(stats))
	}
	pnlSearch := em["pnlSearch"]
	s.inputSearch = bootstrap.NewInput(bootstrap.InputTypeText)
	s.inputSearch.SetClass("form-control input-lg")
	s.inputSearch.SetAttribute("placeholder", "Search anything")
	s.inputSearch.SetAttribute("autofocus", "true")

	btnSearch := bootstrap.NewButton(bootstrap.ButtonPrimary, "Search")
	btnSearch.SetClass("btn btn-lg btn-primary")
	btnSearch.OnEvent(gowd.OnClick, s.btnSearchClick)
	pnlSearch.AddElement(s.inputSearch)
	pnlSearch.AddElement(bootstrap.NewElement("span", "input-group-btn", btnSearch))

	pnlAdvanced := em["pnlAdvanced"]
	pnlAdvanced.AddElement(querybuilder.NewBuilder().Element)

	btnDuplicates := bootstrap.NewButton(bootstrap.ButtonDefault, "Duplicates")
	pnlButtons := em["pnlButtons"]
	pnlButtons.AddElement(btnDuplicates)

	return s
}

func (s *search) updateState() {
	stats, err := db.Stats()
	if err != nil {
		Root.addAlertError(err)
	} else {
		s.txtStats.SetText(fmt.Sprintf("%d images, %d files", stats.ImagesCount, stats.FilesCount))
	}
}

func (s *search) getContent() *gowd.Element {
	return s.Element
}

func (s *search) populateToolbar(menu *darktheme.Menu) {

}

func (s *search) btnDuplicatesClick(sender *gowd.Element, event *gowd.EventElement) {

}

func (s *search) btnSearchClick(sender *gowd.Element, event *gowd.EventElement) {
	var err error
	term := s.inputSearch.GetValue()
	activeSearch, err = model.NewSearch(term)
	if err != nil {
		Root.addAlertError(err)
		return
	}
	Root.setActiveView(Root.thumb)
}
