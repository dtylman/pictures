package db

import "github.com/dtylman/pictures/indexer/picture"

type Query interface {
	Query() error
	Results() []*picture.Index
}
