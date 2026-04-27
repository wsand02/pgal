package database

import (
	"database/sql"
	"log"
	"sync"
)

var (
	once sync.Once
	db   *sql.DB
)

func DB() *sql.DB {
	once.Do(func() {
		var err error
		db, err = sql.Open("sqlite", "file::memory:?cache=shared")
		if err != nil {
			log.Fatal(err)
		}
	})
	return db
}
