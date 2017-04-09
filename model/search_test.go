package model

import (
	"github.com/blevesearch/bleve"
	"github.com/dtylman/pictures/indexer/db"
	"testing"
)

func TestNewSearch(t *testing.T) {
	err := db.Open()
	if err != nil {
		t.Fatal(err)
	}
	s, err := NewSearch("", bleve.NewMatchAllQuery())
	if err != nil {
		t.Fatal(err)
	}
	t.Log(s)
	s.buildThumbs()
}
