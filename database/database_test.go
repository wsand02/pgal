package database

import (
	_ "modernc.org/sqlite"
	"testing"
)

func Test_SetupDB(t *testing.T) {
	db, err := SetupDB(":memory:")
	if err != nil {
		t.Fatal(err)
	}
	err = db.Ping()
	if err != nil {
		t.Fatal(err)
	}
}
