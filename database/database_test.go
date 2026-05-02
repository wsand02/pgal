package database

import (
	_ "modernc.org/sqlite"
	"testing"
)

func Test_SetupDB(t *testing.T) {
	err := GetDB().Ping()
	if err != nil {
		t.Fatal(err)
	}
}
