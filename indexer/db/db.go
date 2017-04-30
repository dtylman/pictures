package db

import (
	"github.com/dtylman/pictures/indexer/picture"
	"database/sql"
	"time"
	"log"
)

const (
	PhaseThumb = "thumb"
	PhaseLocation = "location"
	PhaseObjects = "objects"
	PhaseFaces = "faces"
)

//BatchIndex updates batch of pictures
func BatchIndex(images *picture.Queue) error {
	tx, err := sqldb.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	_, picture := images.Pop()
	for picture != nil {
		err := indexWithTx(tx, picture)
		if err != nil {
			return err
		}
		_, picture = images.Pop()
	}
	return tx.Commit()
}

func indexWithTx(tx *sql.Tx, picture*picture.Index) error {
	_, err := tx.Exec(`INSERT OR REPLACE INTO picture VALUES (?,?,?,?,?,?,?,?,?)`,
		picture.MD5, picture.MimeType, picture.Taken.Unix(), picture.Lat, picture.Long, picture.Location, picture.Objects,
		picture.ExifString(), picture.Faces)
	if err != nil {
		return err
	}
	_, err = tx.Exec(`INSERT OR REPLACE INTO file VALUES (?,?,?,?)`,
		picture.MD5, picture.Path, picture.Album, picture.FileTime.Unix())
	if err != nil {
		return err
	}
	for name, value := range picture.Exif {
		_, err = tx.Exec(`INSERT OR REPLACE INTO exif VALUES (?,?,?)`, picture.MD5, name, value)
		if err != nil {
			return err
		}
	}
	return nil
}

//Index saves one picture into the database
func Index(picture *picture.Index) error {
	log.Println("Index started")
	defer log.Println("Index ended")

	tx, err := sqldb.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	err = indexWithTx(tx, picture)
	if err != nil {
		return err
	}
	return tx.Commit()
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
				objects TEXT,
				exif TEXT,
				faces TEXT) WITHOUT ROWID;`,
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
		`CREATE UNIQUE INDEX idx_processor on processor (md5, name);
		CREATE TABLE location (
			md5 TEXT NOT NULL PRIMARY KEY,
			street TEXT,
			citi TEXT,
			state TEXT,
			postalcode TEXT,
			country TEXT
			) WITHOUT ROWID;
		CREATE TABLE object (
			md5 TEXT NOT NULL PRIMARY KEY,
			name TEXT NOT NULL,
			prob NUMBER NOT NULL
			) WITHOUT ROWID`,
		`CREATE VIEW images_view AS
			SELECT DISTINCT picture.md5, mime_type,	file.path,
			taken, lat, long, location, album, objects, faces
			FROM picture
			INNER JOIN file ON file.md5=picture.md5
			ORDER BY taken, file.time`,
	}
	return execMultiple(schema)
}

func execMultiple(sql []string) error {
	tx, err := sqldb.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	for _, statement := range sql {
		_, err := tx.Exec(statement)
		if err != nil {
			return err
		}
	}
	return tx.Commit()
}

//SetPhase sets phase if does not exists and return true if it was set.
func SetPhase(md5 string, name string) bool {
	tx, err := sqldb.Begin()
	if err != nil {
		log.Println(err)
		return false
	}
	defer tx.Rollback()
	var count int
	err = tx.QueryRow("SELECT count(*) AS count FROM processor WHERE md5=? AND name=?", md5, name).Scan(&count)
	if err != nil {
		log.Println(err)
		return false
	}
	if count > 0 {
		return false
	}
	res, err := tx.Exec(`INSERT INTO processor VALUES (?,?,?) `, md5, name, time.Now().Unix())
	if err != nil {
		log.Println(err)
		return false
	}
	err = tx.Commit()
	if err != nil {
		log.Println(err)
		return false
	}
	ar, err := res.RowsAffected()
	if err != nil {
		log.Println(err)
		return false
	}
	return ar > 0
}

//HasPath returns true if a path exists in the database
func HasPath(path string) (bool, error) {
	var count int
	err := sqldb.QueryRow(`SELECT count(*) AS count FROM file WHERE path=?`, path).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

//HasImage returns tue if image with md5
func HasImage(md5 string) bool {
	var count int
	err := sqldb.QueryRow(`SELECT count(*) AS count FROM picture WHERE md5=?`, md5).Scan(&count)
	if err != nil {
		log.Println(err)
		return false
	}
	return count > 0
}

func removeWithTx(tx*sql.Tx, md5 string) error {
	_, err := tx.Exec(`DELETE FROM picture WHERE md5=?`, md5)
	if err != nil {
		return err
	}
	_, err = tx.Exec(`DELETE FROM file WHERE md5=?`, md5)
	if err != nil {
		return err
	}
	_, err = tx.Exec(`DELETE FROM exif WHERE md5=?`, md5)
	if err != nil {
		return err
	}
	_, err = tx.Exec(`DELETE FROM location WHERE md5=?`, md5)
	if err != nil {
		return err
	}
	_, err = tx.Exec(`DELETE FROM object WHERE md5=?`, md5)
	if err != nil {
		return err
	}
	_, err = tx.Exec(`DELETE FROM processor WHERE md5=?`, md5)
	if err != nil {
		return err
	}
	return nil
}

//Remove removes all items in keys
func Remove(keys []string) error {
	tx, err := sqldb.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	for _, md5 := range keys {
		err = removeWithTx(tx, md5)
		if err != nil {
			return err
		}
	}
	return tx.Commit()
}


//WalkImagesFunc defines a callback to scan alll images in database (use with WalkImages)
type WalkImagesFunc func(key string, image *picture.Index, err error)

//WalkImages executes function for all images in the database
func WalkImages(wf WalkImagesFunc) {
	rows, err := sqldb.Query(`SELECT DISTINCT picture.md5,
			mime_type, file.path, taken, lat, long,	location, album, objects, faces
			FROM picture JOIN file ON file.md5=picture.md5`)
	if err != nil {
		wf("", nil, err)
		return
	}
	defer rows.Close()
	for (rows.Next()) {
		image, err := rows2Image(rows)
		wf(image.MD5, image, err)
	}
}

func rows2Image(rows*sql.Rows) (*picture.Index, error) {
	var image picture.Index
	var taken int64
	err := rows.Scan(&image.MD5, &image.MimeType, &image.Path, &taken, &image.Lat,
		&image.Long, &image.Location, &image.Album, &image.Objects, &image.Faces)
	if err != nil {
		return nil, err
	}
	image.Taken = time.Unix(taken, 0)
	return &image, nil
}