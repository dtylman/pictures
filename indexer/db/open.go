package db

import (
	"os"

	"github.com/dtylman/pictures/conf"

	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

var (
	sqldb *sql.DB
)

func openSQlite() error {
	path, err := conf.SqlitePath()
	if err != nil {
		return err
	}

	sqldb, err = sql.Open("sqlite3", path+"?_busy_timeout=10000")
	if err != nil {
		return err
	}
	_, err = os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			err = createSchema()
			if err != nil {
				return nil
			}
		}
	}
	//setup session
	_, err = sqldb.Exec(`PRAGMA synchronous=OFF`)
	if err != nil {
		return err
	}
	_, err = sqldb.Exec(`PRAGMA count_changes=OFF`)
	if err != nil {
		return err
	}
	_, err = sqldb.Exec(`PRAGMA journal_mode=MEMORY`)
	if err != nil {
		return err
	}
	_, err = sqldb.Exec(`PRAGMA temp_store=MEMORY`)
	if err != nil {
		return err
	}
	return nil
}

//Open opens or creates the local database
func Open() error {
	return openSQlite()
}

//Close closes the db
func Close() {
	err := sqldb.Close()
	if err != nil {
		log.Println(err)
	}
}

//DeleteDatabase removes the database and all data
func DeleteDatabase() error {
	path, err := conf.FilesPath()
	if err != nil {
		return err
	}
	err = os.RemoveAll(path)
	if err != nil {
		return err
	}
	Close()
	path, err = conf.SqlitePath()
	if err != nil {
		return err
	}
	return os.Remove(path)
}
