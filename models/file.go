package models

import (
	"fmt"

	"github.com/wsand02/pgal/database"
)

type File struct {
	ID       int64
	Name     string
	RealPath string
	FolderID int64
}

func FileSchema() string {
	return `CREATE TABLE file (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		real_path TEXT UNIQUE NOT NULL,
		folder_id INTEGER NOT NULL,
		FOREIGN KEY (folder_id) REFERENCES folder(id)
	);`
}

func FileEqual(a, b File) bool {
	if a.ID != b.ID || a.Name != b.Name || a.RealPath != b.RealPath || a.FolderID != b.FolderID {
		return false
	}
	return true
}

func NewFile(id int64, name, rp string, fid int64) File {
	return File{ID: id, Name: name, RealPath: rp, FolderID: fid}
}

func AddFile(name, real_path string, folder_id int64) (int64, error) {
	res, err := database.GetDB().Exec("INSERT INTO file(name, real_path, folder_id) VALUES(?, ?, ?)", name, real_path, folder_id)
	if err != nil {
		return 0, fmt.Errorf("addFile: %v", err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("addFile: %v", err)
	}
	return id, nil
}

func Files() ([]File, error) {
	rows, err := database.GetDB().Query("SELECT * FROM file")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var files []File
	for rows.Next() {
		var fi File
		if err := rows.Scan(&fi.ID, &fi.Name, &fi.RealPath, &fi.FolderID); err != nil {
			return files, err
		}
		files = append(files, fi)
	}
	if err = rows.Err(); err != nil {
		return files, err
	}
	return files, nil
}
