package services

import (
	"database/sql"
	"fmt"

	"github.com/wsand02/pgal/models"
)

type FileService struct {
	db *sql.DB
}

func NewFileService(db *sql.DB) *FileService {
	return &FileService{db}
}

func (s *FileService) AddFile(name, real_path string, folder_id int64) (int64, error) {
	res, err := s.db.Exec("INSERT INTO file(name, real_path, folder_id) VALUES(?, ?, ?)", name, real_path, folder_id)
	if err != nil {
		return 0, fmt.Errorf("addFile: %v", err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("addFile: %v", err)
	}
	return id, nil
}

func (s *FileService) Files() ([]models.File, error) {
	rows, err := s.db.Query("SELECT * FROM file")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var files []models.File
	for rows.Next() {
		var fi models.File
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
