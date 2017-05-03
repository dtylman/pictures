package db

import (
	"fmt"

	"github.com/dtylman/pictures/indexer/picture"
)

type TermQuery struct {
	term  string
	exact bool
	limit int
	res   []*picture.Index
}

const NOLIMIT = 0

func NewTermQuery(term string, exact bool, limit int) *TermQuery {
	return &TermQuery{term: term, exact: exact, limit: limit}
}

func (tq *TermQuery) Query() error {
	tq.res = make([]*picture.Index, 0)
	var operator string
	if tq.exact {
		operator = fmt.Sprintf(` = '%s' `, tq.term)
	} else {
		operator = fmt.Sprintf(` LIKE '%%%s%%' `, tq.term)
	}
	sql := `SELECT DISTINCT 
			md5, mime_type,	path,
			taken, lat, long, location, album, objects, faces 
			FROM images_view 
			WHERE mime_type ` + operator +
		` OR path ` + operator +
		` OR location ` + operator +
		` OR album ` + operator +
		` OR objects ` + operator +
		` OR faces` + operator
	if tq.limit > 0 {
		sql += fmt.Sprintf(" LIMIT %d", tq.limit)
	}
	if !tq.exact {
		sql += ` COLLATE NOCASE`
	}
	rows, err := sqldb.Query(sql)
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		image, err := rows2Image(rows)
		if err != nil {
			return err
		}
		tq.res = append(tq.res, image)
	}
	return nil
}

func (tq *TermQuery) Results() []*picture.Index {
	return tq.res
}
