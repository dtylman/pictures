package model

import (
	"github.com/dtylman/pictures/indexer/db"
	"github.com/dtylman/pictures/indexer/picture"
	"github.com/dtylman/pictures/indexer/thumbs"
	"log"
	"github.com/dtylman/pictures/conf"
	"strconv"
	"strings"
	"sort"
)

type ThumbItem struct {
	Path string
	MD5  string
}

type Search struct {
	query       string
	hit         int
	start       int
	Results     []*picture.Index
	ActiveImage *picture.Index
	Pages       Pages
	Facets      Facets
	Thumbs      []ThumbItem
}

func NewSearch(query string) (*Search, error) {
	s := new(Search)
	s.query = query
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
	var q db.Query
	if s.query == "duplicates" {
		q = db.NewStaticQuery(db.QueryDuplicates)
	} else {
		q = db.NewTermQuery(s.query, false, db.NOLIMIT)
	}
	err := q.Query()
	if err != nil {
		return err
	}
	s.Results = q.Results()
	s.start = 0
	s.hit = 0
	s.buildFacetItems()
	s.buildPages()
	s.buildThumbs()
	return nil
}

func (s *Search) buildPages() {
	if conf.Options.SearchPageSize == 0 {
		s.Pages = make(Pages, 0)
		return
	}
	pageCount := s.Total() / conf.Options.SearchPageSize
	fromPage := s.start / conf.Options.SearchPageSize
	s.Pages = make(Pages, pageCount)
	for i := 0; i < pageCount; i++ {
		s.Pages[i].Start = i * conf.Options.SearchPageSize
		s.Pages[i].Active = (i == fromPage)
		s.Pages[i].Caption = strconv.Itoa(s.Pages[i].Start)
	}
}

func (s *Search) buildFacetItems() {
	facetMap := make(map[string]int)
	for _, image := range s.Results {
		for _, term := range strings.Split(image.Album + " " + image.Location, " ") {
			if term != "" {
				facetMap[term]++
			}
		}
	}
	s.Facets = make(Facets, 0)

	for term, count := range facetMap {
		s.Facets = append(s.Facets, FacetItem{Term: term, Count: count})
	}
	sort.Sort(s.Facets)
	if s.Facets.Len() > 25 {
		s.Facets = s.Facets[0:25]
	}
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