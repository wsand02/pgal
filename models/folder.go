package models

import (
	"strconv"
	"strings"
)

type Folder struct {
	ID       int64
	Name     string
	RealPath string
	ParentID *int64
	// Parent   *Folder
}

func FolderSchema() string {
	return `CREATE TABLE folder (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		real_path TEXT UNIQUE NOT NULL,
		parent_id INTEGER,
		FOREIGN KEY (parent_id) REFERENCES folder(id)
	);`
}

func NewOrphanFolder(id int64, name, rp string) Folder {
	return Folder{ID: id, Name: name, RealPath: rp}
}

func NewFolder(id int64, name, rp string, pid int64) Folder {
	return Folder{ID: id, Name: name, RealPath: rp, ParentID: &pid}
}

func (f *Folder) URL() string {
	return strings.Join([]string{"/folders/", strconv.Itoa(int(f.ID))}, "")
}

func FoldersEqual(a, b Folder) bool {
	if a.ID != b.ID || a.Name != b.Name || a.RealPath != b.RealPath {
		return false
	}
	if (a.ParentID == nil && b.ParentID == nil) || (a.ParentID != nil && b.ParentID != nil && *a.ParentID == *b.ParentID) {
		return true
	}
	return false
}
