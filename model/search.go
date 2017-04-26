package model

import (
	"github.com/dtylman/pictures/indexer/db"
	"github.com/dtylman/pictures/indexer/picture"
	"github.com/dtylman/pictures/indexer/thumbs"
	"log"
	"github.com/dtylman/pictures/conf"
	"strconv"
	"github.com/blevesearch/bleve"
	_ "github.com/blevesearch/bleve/analysis/analyzer/keyword"
)

//FacetItem represents facet item in the display
type FacetItem struct {
	Name  string
	Field string
	Term  string
	Count int
}

//PageItem represents a paging item
type PageItem struct {
	Start   int
	Active  bool
	Caption string
}

type ThumbItem struct {
	Path string
	MD5  string
}

type Search struct {
	term        string
	hit         int
	start       int
	Results     []*picture.Index
	ActiveImage *picture.Index
	Pages       []PageItem
	Facets      []FacetItem
	Thumbs      []ThumbItem
}

func NewSearch(term string) (*Search, error) {
	s := new(Search)
	s.term = term
	err := s.doQuery()
	if err != nil {
		return nil, err
	}
	return s, nil
}

func (s *Search) SetActiveImage(hit int) {
	s.hit = hit
	s.ActiveImage = s.Results[s.start + s.hit]
}

func (s *Search) NextImage() {
	nextHit := s.hit + 1
	if nextHit + s.start >= s.Total() {
		//nowhere to go
		return
	}
	if nextHit >= conf.Options.SearchPageSize {
		s.NextPage()
		return
	}
	s.SetActiveImage(nextHit)
}

func (s *Search) PrevImage() {
	prevHit := s.hit - 1
	if prevHit + s.start < 0 {
		s.PrevPage()
		return
	}
	s.SetActiveImage(prevHit)
}

func (s *Search) NextPage() {
	from := s.start + conf.Options.SearchPageSize
	if from >= s.Total() {
		//no where to go
		return
	}
	s.start = from
	s.SetActiveImage(0)
	s.buildPages()
	s.buildThumbs()
}

func (s *Search) PrevPage() {
	from := s.start - conf.Options.SearchPageSize
	if from < 0 {
		//no where to go
		return
	}
	s.start = from
	s.SetActiveImage(0)
	s.buildPages()
	s.buildThumbs()
}

func (s *Search) doQuery() error {
	tq := db.NewTermQuery(s.term, false, db.NOLIMIT)
	err := tq.Query()
	if err != nil {
		return err
	}
	s.Results = tq.Result
	s.start = 0
	s.hit = 0
	err = s.buildFacetItems()
	if err != nil {
		return err
	}
	s.buildPages()
	s.buildThumbs()
	return nil
}

func (s *Search) buildPages() {
	if conf.Options.SearchPageSize == 0 {
		s.Pages = make([]PageItem, 0)
		return
	}
	pageCount := s.Total() / conf.Options.SearchPageSize
	fromPage := s.start / conf.Options.SearchPageSize
	s.Pages = make([]PageItem, pageCount)
	for i := 0; i < pageCount; i++ {
		s.Pages[i].Start = i * conf.Options.SearchPageSize
		s.Pages[i].Active = (i == fromPage)
		s.Pages[i].Caption = strconv.Itoa(s.Pages[i].Start)
	}
}

func (s *Search) buildFacetItems() error {
	keywordIndex := bleve.NewTextFieldMapping()
	keywordIndex.Store = false
	keywordIndex.IncludeInAll = false
	keywordIndex.IncludeTermVectors = false
	keywordIndex.Analyzer = "keyword"
	mapping := bleve.NewDocumentMapping()
	mapping.AddFieldMappingsAt("location", keywordIndex)
	mapping.AddFieldMappingsAt("album", keywordIndex)
	disabledSection := bleve.NewDocumentDisabledMapping()
	mapping.AddSubDocumentMapping("_all", disabledSection)

	indexMapping := bleve.NewIndexMapping()
	indexMapping.DefaultMapping = mapping
	indexMapping.DefaultAnalyzer = "standard"

	index, err := bleve.NewMemOnly(indexMapping)
	if err != nil {
		return err
	}
	defer index.Close()
	s.Facets = make([]FacetItem, 0)
	b := index.NewBatch()
	for _, image := range s.Results {
		err = b.Index(image.MD5, image)
		if err != nil {
			return err
		}
	}
	err = index.Batch(b)
	if err != nil {
		return err
	}
	req := bleve.NewSearchRequest(bleve.NewMatchAllQuery())
	req.AddFacet("Location", bleve.NewFacetRequest("location", 6))
	req.AddFacet("Album", bleve.NewFacetRequest("album", 4))
	sr, err := index.Search(req)
	if err != nil {
		return err
	}
	for fn, fr := range sr.Facets {
		for _, term := range fr.Terms {
			s.Facets = append(s.Facets, FacetItem{Name: fn, Field: fr.Field, Term: term.Term, Count: term.Count})
		}
	}
	return nil
}

func (s *Search) buildThumbs() {
	thumbsCount := s.Total()
	if thumbsCount > conf.Options.SearchPageSize {
		thumbsCount = conf.Options.SearchPageSize
	}
	s.Thumbs = make([]ThumbItem, thumbsCount)
	for i := 0; i < thumbsCount; i++ {
		s.Thumbs[i].MD5 = s.Results[s.start + i].MD5
		var err error
		s.Thumbs[i].Path, err = thumbs.MakeThumb(s.Results[s.start + i].Path, s.Thumbs[i].MD5, false)
		if err != nil {
			log.Println(err)
		}
	}

}

func (s*Search) StartFrom(start int) {
	if start <= s.Total() {
		s.start = start
		s.SetActiveImage(0)
		s.buildPages()
		s.buildThumbs()
	}
}

func (s*Search) Total() int {
	if s.Results != nil {
		return len(s.Results)
	}
	return 0
}