package models

import (
	"fmt"
	"strconv"
	"strings"
	"database/sql"

	"github.com/wsand02/pgal/database"
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

func AddRootFolder(name, real_path string) (int64, error) {
	res, err := database.GetDB().Exec("INSERT INTO folder (name, real_path) VALUES (?, ?)", name, real_path)
	if err != nil {
		return 0, fmt.Errorf("AddRootFolder: %v", err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("AddRootFolder: %v", err)
	}
	return id, nil
}

func AddChildFolder(name, real_path string, parent_id int64) (int64, error) {
	res, err := database.GetDB().Exec("INSERT INTO folder(name, real_path, parent_id) VALUES(?, ?, ?)", name, real_path, parent_id)
	if err != nil {
		return 0, fmt.Errorf("AddChildFolder: %v, %v", err, real_path)
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("AddChildFolder: %v, %v", err, real_path)
	}
	return id, nil
}

func ParentId(path string) (int64, error) {
	var parent_id int64 = 0
	if err := database.GetDB().QueryRow("SELECT id FROM folder WHERE real_path = ?", path).Scan(&parent_id); err != nil {
		if err == sql.ErrNoRows {
			return 0, fmt.Errorf("ParentId: %v", path)
		}
		return 0, fmt.Errorf("ParentId: %v", err)
	}
	return parent_id, nil
}

func Folders() ([]Folder, error) {
	rows, err := database.GetDB().Query("SELECT * FROM folder")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var folders []Folder
	for rows.Next() {
		var fol Folder
		if err := rows.Scan(&fol.ID, &fol.Name, &fol.RealPath, &fol.ParentID); err != nil {
			return folders, err
		}
		folders = append(folders, fol)
	}
	if err = rows.Err(); err != nil {
		return folders, err
	}
	return folders, nil
}

func GetFolder(id int64) (*Folder, error) {
	var fo Folder
	if err := database.GetDB().QueryRow("SELECT * FROM folder WHERE id = ?", id).Scan(&fo.ID, &fo.Name, &fo.RealPath, &fo.ParentID); err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("Folder: %v", id)
		}
	}
	return &fo, nil
}

func FoldersByParent(id int64) ([]Folder, error) {
	rows, err := database.GetDB().Query("SELECT * FROM folder WHERE parent_id = ?", id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var folders []Folder
	for rows.Next() {
		var fol Folder
		if err := rows.Scan(&fol.ID, &fol.Name, &fol.RealPath, &fol.ParentID); err != nil {
			return folders, err
		}
		folders = append(folders, fol)
	}
	if err = rows.Err(); err != nil {
		return folders, err
	}
	return folders, nil
}
