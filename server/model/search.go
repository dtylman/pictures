package model

import (
	"github.com/blevesearch/bleve"
	"github.com/blevesearch/bleve/search/query"
	"github.com/dtylman/pictures/conf"
	"github.com/dtylman/pictures/indexer/db"
	"path/filepath"
	"strconv"
)

const maxPages = 10

//FacetItem represents facet item in the display
type FacetItem struct {
	Name  string
	Field string
	Term  string
	Count int
}

//PageItem represents a paging item
type PageItem struct {
	From    int
	Active  bool
	Caption string
}

//ImageItem represents one image being now presented
type ImageItem struct {
	Hit      int
	Name     string
	Path     string
	ID       string
	Location string
	Details  map[string]interface{}
}

type Search struct {
	QueryString string
	Result      *bleve.SearchResult
	query       query.Query
	start       int
	ActiveImage ImageItem
	Pages       []PageItem
	Facets      []FacetItem
	Thumbs      []ThumbItem
}

func NewSearch(queryString string, query query.Query) (*Search, error) {
	s := new(Search)
	s.query = query
	s.QueryString = queryString
	s.start = 0
	s.doQuery()
	return s, nil
}

func (s *Search) StartFrom(start int) error {
	s.start = start
	return s.doQuery()
}

func (s *Search) SetActiveHit(hit int) {
	s.ActiveImage.ID = s.Result.Hits[hit].ID
	s.ActiveImage.Hit = hit
	s.ActiveImage.Details = s.Result.Hits[hit].Fields
	s.ActiveImage.Path = s.Result.Hits[hit].Fields["path"].(string)
	s.ActiveImage.Location = s.Result.Hits[hit].Fields["location"].(string)
	s.ActiveImage.Name = filepath.Base(s.ActiveImage.Path)
}

func (s *Search) NextHit() {
	s.SetActiveHit(s.ActiveImage.Hit + 1)
}

func (s *Search) PrevHit() {
	s.SetActiveHit(s.ActiveImage.Hit - 1)
}

func (s *Search) doQuery() error {
	req := bleve.NewSearchRequestOptions(s.query, conf.Options.SearchPageSize, s.start, false)
	req.Fields = []string{"*"}
	req.AddFacet("Location", bleve.NewFacetRequest("location", 6))
	var err error
	s.Result, err = db.Search(req)
	if err != nil {
		return err
	}
	s.buildFacetItems()
	s.buildPages()
	s.buildThumbs()
	return nil
}

func (s *Search) buildPages() {
	if conf.Options.SearchPageSize == 0 {
		s.Pages = make([]PageItem, 0)
		return
	}
	pageCount := int(s.Result.Total) / conf.Options.SearchPageSize
	pages := maxPages
	if pageCount < 10 {
		pages = pageCount
	}
	fromPage := s.start / conf.Options.SearchPageSize
	s.Pages = make([]PageItem, pages)
	for i := 0; i < pages; i++ {
		s.Pages[i].From = i * conf.Options.SearchPageSize
		s.Pages[i].Active = (i == fromPage)
		s.Pages[i].Caption = strconv.Itoa(s.Pages[i].From)
	}
}

func (s *Search) buildFacetItems() {
	s.Facets = make([]FacetItem, 0)
	for fn, fr := range s.Result.Facets {
		for _, term := range fr.Terms {
			s.Facets = append(s.Facets, FacetItem{Name: fn, Field: fr.Field, Term: term.Term, Count: term.Count})
		}
	}
}