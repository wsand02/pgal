package services

import (
	"database/sql"
	"fmt"

	"github.com/wsand02/pgal/models"
)

type FolderService struct {
	db *sql.DB
}

func NewFolderService(db *sql.DB) *FolderService {
	return &FolderService{db: db}
}

func (s *FolderService) AddRootFolder(name, real_path string) (int64, error) {
	res, err := s.db.Exec("INSERT INTO folder (name, real_path) VALUES (?, ?)", name, real_path)
	if err != nil {
		return 0, fmt.Errorf("AddRootFolder: %v", err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("AddRootFolder: %v", err)
	}
	return id, nil
}

func (s *FolderService) AddChildFolder(name, real_path string, parent_id int64) (int64, error) {
	res, err := s.db.Exec("INSERT INTO folder(name, real_path, parent_id) VALUES(?, ?, ?)", name, real_path, parent_id)
	if err != nil {
		return 0, fmt.Errorf("AddChildFolder: %v, %v", err, real_path)
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("AddChildFolder: %v, %v", err, real_path)
	}
	return id, nil
}

func (s *FolderService) ParentId(path string) (int64, error) {
	var parent_id int64 = 0
	if err := s.db.QueryRow("SELECT id FROM folder WHERE real_path = ?", path).Scan(&parent_id); err != nil {
		if err == sql.ErrNoRows {
			return 0, fmt.Errorf("ParentId: %v", path)
		}
		return 0, fmt.Errorf("ParentId: %v", err)
	}
	return parent_id, nil
}

func (s *FolderService) Folders() ([]models.Folder, error) {
	rows, err := s.db.Query("SELECT * FROM folder")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var folders []models.Folder
	for rows.Next() {
		var fol models.Folder
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

func (s *FolderService) Folder(id int64) (*models.Folder, error) {
	var fo models.Folder
	if err := s.db.QueryRow("SELECT * FROM folder WHERE id = ?", id).Scan(&fo.ID, &fo.Name, &fo.RealPath, &fo.ParentID); err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("Folder: %v", id)
		}
	}
	return &fo, nil
}
