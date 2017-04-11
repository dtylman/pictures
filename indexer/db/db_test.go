package db

import (
	"testing"
	"github.com/dtylman/pictures/indexer/picture"
	"os"
	"github.com/dtylman/pictures/indexer/darknet"
)

func TestIndex(t *testing.T) {
	err := Open()
	if err != nil {
		t.Fatal(err)
	}
	defer Close()

	fileInfo, err := os.Stat("/tmp/lala.jpg")
	if err != nil {
		t.Fatal(err)
	}
	i, err := picture.NewIndex("/tmp/lala.jpg", fileInfo)
	if err != nil {
		t.Fatal(err)
	}
	err = Index(i)
	if err != nil {
		t.Fatal(err)
	}
	i.Objects = make([]darknet.Object, 2)
	i.Objects[0].Name = "horse"
	i.Objects[1].Name = "dog"
	i.Objects[1].Count = 2
	i.SetPhase("data")
	err = Index(i)
	if err != nil {
		t.Fatal(err)
	}
}