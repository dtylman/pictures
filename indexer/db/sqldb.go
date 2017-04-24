package db

import (
	"github.com/dtylman/pictures/indexer/picture"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

/*
MD5      string    `json:"md5"`
	MimeType string    `json:"mime_type"`
	Path     string    `json:"path"`
	FileTime time.Time `json:"file_time"`
	Taken    time.Time `json:"taken"`
	Exif     string    `json:"exif"`
	Lat      float64   `json:"lat"`
	Long     float64   `json:"long"`
	Location string    `json:"location"`
	Album    string    `json:"album"`
	Objects  string `json:"objects"`
	Phases   map[string]time.Time `json:"phases"`
*/

//BatchIndex updates batch of pictures
func BatchIndex(images *picture.Queue) {
	_, picture := images.Pop()
	for picture != nil {
		err := Index(picture)
		if err != nil {
			log.Println(err)
		}
		_, picture = images.Pop()
	}
}

//Index saves one picture into the database
func Index(picture *picture.Index) error {
	tx, err := sqldb.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	_, err = tx.Exec(`INSERT OR REPLACE INTO picture VALUES (?,?,?,?,?,?,?)`,
		picture.MD5, picture.MimeType, picture.Taken.Unix(), picture.Lat, picture.Long, picture.Location, picture.Objects)
	if err != nil {
		return err
	}
	_, err = tx.Exec(`INSERT OR REPLACE INTO file VALUES (?,?,?,?)`,
		picture.MD5, picture.Path, picture.Album, picture.FileTime.Unix())
	if err != nil {
		return err
	}
	return tx.Commit()
	//	sqlite.Exec("insert into files values(?,?)",picture.MD5,picture.Path, picture.Album, picture.FileTime)
	//	sqlite.Exec("insert into exif values(?,?,?)",picture.MD5,picture.Exif)
	//	sqlite.Exec("insert into phases values(?,?,?",picture.MD5//,name, time)
	//
	//	//query:
	//	//select where location like '' or exif like '' or path like '' or
	//
}

func createSchema() error {
	schema := []string{
		`CREATE TABLE picture (
				md5 TEXT PRIMARY KEY,
				mime_type TEXT NOT NULL,
				taken INT,
				lat REAL,
				long REAL,
				location TEXT,
				objects TEXT) WITHOUT ROWID;`,
		`CREATE UNIQUE INDEX idx_picture on picture (md5 ASC);`,
		`CREATE TABLE file (
			    md5 TEXT NOT NULL,
			    path TEXT NOT NULL PRIMARY KEY,
			    album TEXT NOT NULL,
			    time INTEGER NOT NULL
			) WITHOUT ROWID;`,
		`CREATE UNIQUE INDEX idx_file on file (path ASC);`,
		`CREATE TABLE exif (
			    md5 TEXT NOT NULL PRIMARY KEY,
			    name TEXT NOT NULL,
			    value TEXT
			) WITHOUT ROWID;
			CREATE UNIQUE INDEX idx_exif on exif (md5 ASC);`,
		`CREATE TABLE processor (
			    md5 TEXT NOT NULL PRIMARY KEY,
			    name TEXT NOT NULL,
			    time INTEGER
			    ) WITHOUT ROWID;`,
		`CREATE UNIQUE INDEX idx_processor on processor (md5, name);`,
	}
	tx, err := sqldb.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	for _, statement := range schema {
		_, err := tx.Exec(statement)
		if err != nil {
			return err
		}
	}
	return tx.Commit()
}


/*
some samples:
 > get albums facet
--select album,  count(*) as total  from file group by album order by total desc limit 5

 > find duplicates
--select md5, path ,count(*) as count from file group by md5,path having count > 1
 */