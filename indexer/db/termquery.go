package db

import (
	"github.com/dtylman/pictures/indexer/picture"
	"fmt"
)

type TermQuery struct {
	term   string
	exact  bool
	limit  int
	Result []*picture.Index
}

const NOLIMIT = 0

func NewTermQuery(term string, exact bool, limit int) *TermQuery {
	return &TermQuery{term:term, exact:exact, limit: limit}
}

func (tq*TermQuery) Query() error {
	tq.Result = make([]*picture.Index, 0)
	var operator string
	if tq.exact {
		operator = fmt.Sprintf(` = '%s' `, tq.term)
	} else {
		operator = fmt.Sprintf(` LIKE '%%%s%%' `, tq.term)
	}
	sql := `SELECT DISTINCT picture.md5,
	mime_type, file.path, taken, lat, long,	location, album, objects, faces
	FROM picture JOIN file ON file.md5=picture.md5
	WHERE
	picture.mime_type ` + operator + ` OR
	file.path ` + operator + ` OR
	location ` + operator + ` OR
	album ` + operator + ` OR
	objects ` + operator + ` OR
	faces` + operator + `
	ORDER BY picture.taken, file.time `
	if tq.limit > 0 {
		sql += fmt.Sprintf("LIMIT %d", tq.limit)
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
		tq.Result = append(tq.Result, image)
	}
	return nil
}

