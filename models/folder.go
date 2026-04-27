package models

import (
	"database/sql"
	"fmt"

	"github.com/wsand02/pgal/database"
)

type Folder struct {
	ID       int64
	Name     string
	RealPath string
	ParentID *int64
	// Parent   *Folder
}

func folderSchema() string {
	return `CREATE TABLE folder (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		real_path TEXT UNIQUE NOT NULL,
		parent_id INTEGER,
		FOREIGN KEY (parent_id) REFERENCES folder(id)
	);`
}

func AddRootFolder(name, real_path string) (int64, error) {
	res, err := database.DB().Exec("INSERT INTO folder (name, real_path) VALUES (?, ?)", name, real_path)
	if err != nil {
		return 0, fmt.Errorf("addRootFolder: %v", err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("addRootFolder: %v", err)
	}
	return id, nil
}

func AddChildFolder(name, real_path string, parent_id int64) (int64, error) {
	res, err := database.DB().Exec("INSERT INTO folder(name, real_path, parent_id) VALUES(?, ?, ?)", name, real_path, parent_id)
	if err != nil {
		return 0, fmt.Errorf("addChildFolder: %v, %v", err, real_path)
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("addChildFolder: %v, %v", err, real_path)
	}
	return id, nil
}

func ParentId(path string) (int64, error) {
	var parent_id int64 = 0
	if err := database.DB().QueryRow("SELECT id FROM folder WHERE real_path = ?", path).Scan(&parent_id); err != nil {
		if err == sql.ErrNoRows {
			return 0, fmt.Errorf("getParentId: %v", path)
		}
		return 0, fmt.Errorf("getParentId: %v", err)
	}
	return parent_id, nil
}

func Folders() ([]Folder, error) {
	rows, err := database.DB().Query("SELECT * FROM folder")
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
