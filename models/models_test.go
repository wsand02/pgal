package models

import (
	"testing"

	"github.com/wsand02/pgal/database"
)

func Test_Models(t *testing.T) {
	database.GetDB().Exec(Schemas())
	var want int64
	rID, err := AddRootFolder(".", ".")
	if err != nil {
		t.Fatal(err)
	}
	want = 1
	if rID != want {
		t.Fatalf("got %v want %v", rID, want)
	}
	cID, err := AddChildFolder("test", "test", 1)
	if err != nil {
		t.Fatal(err)
	}
	want = 2
	if cID != want {
		t.Fatalf("got %v want %v", cID, 1)
	}
	tID, err := AddFile("test.txt", "test/test.txt", 2)
	if err != nil {
		t.Fatal(err)
	}
	want = 1
	if tID != want {
		t.Fatalf("got %v want %v", tID, want)
	}
	fls, err := Folders()
	if err != nil {
		t.Fatal(err)
	}
	want = 2
	var got int64
	got = int64(len(fls))
	if got != want {
		t.Fatalf("got %v want %v", len(fls), want)
	}
	fis, err := Files()
	if err != nil {
		t.Fatal(err)
	}
	want = 1
	got = int64(len(fis))
	if got != want {
		t.Fatalf("got %v want %v", len(fis), want)
	}
}
