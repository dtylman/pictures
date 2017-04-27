package db

import "github.com/dtylman/pictures/indexer/picture"

const (
	QueryDuplicates = 1
)

type StaticQuery struct {
	queryType int
	res       []*picture.Index
}

func NewStaticQuery(queryType int) *StaticQuery {
	return &StaticQuery{queryType:queryType}
}

func (sq*StaticQuery) Query() error {
	sq.res = make([]*picture.Index, 0)
	if sq.queryType == QueryDuplicates {
		rows, err := sqldb.Query(`SELECT picture.md5, mime_type, file.path, taken, lat, long, location, album, objects, faces
		FROM
		(SELECT DISTINCT count(*) AS count ,md5 FROM file GROUP BY md5 HAVING count>1) AS duplicates
		INNER JOIN file ON file.md5=duplicates.md5
		INNER JOIN picture ON picture.md5=duplicates.md5
		ORDER BY file.md5`)
		if err != nil {
			return err
		}
		for rows.Next() {
			image, err := rows2Image(rows)
			if err != nil {
				return err
			}
			sq.res = append(sq.res, image)
		}
	}
	return nil
}

func (sq*StaticQuery) Results() []*picture.Index {
	return sq.res
}