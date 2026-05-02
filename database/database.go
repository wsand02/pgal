package database

import (
	"database/sql"
	"log"
	"sync"

	_ "modernc.org/sqlite"
)

var db *sql.DB
var once sync.Once

func GetDB() *sql.DB {
	var err error
	once.Do(func() {
		db, err = sql.Open("sqlite", "file::memory:?cache=shared&_foreign_keys=on")
		if err != nil {
			log.Fatal(err)
		}
	})
	if err != nil {
		log.Fatal(err)
	}
	return db
}
