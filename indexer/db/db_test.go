package db

import (
	"github.com/blevesearch/bleve"
	"testing"
)

func TestSearch(t *testing.T) {
	req := bleve.NewSearchRequest(bleve.NewDocIDQuery([]string{"8d68af1e1b8408e00889bb63b32cc4ad"}))
	req.Fields = []string{"*"}
	sr, err := Search(req)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(sr.Hits[0].Fields)
}
