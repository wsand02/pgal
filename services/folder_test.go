package services

import (
	"database/sql"
	"slices"
	"testing"

	"github.com/wsand02/pgal/models"
	_ "modernc.org/sqlite"
)

func setupDB() *sql.DB {
	db, _ := sql.Open("sqlite", ":memory:")
	db.Exec(models.FolderSchema())
	return db
}

func TestFolderService_AddRootFolder(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		fname     string
		real_path string
		want      int64
	}{
		{
			name:      "dot",
			fname:     ".",
			real_path: ".",
			want:      1,
		},
		{
			name:      "testfolder",
			fname:     "testfolder",
			real_path: "testfolder",
			want:      1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := setupDB()
			s := NewFolderService(db)
			got, err := s.AddRootFolder(tt.fname, tt.real_path)
			if err != nil {
				t.Errorf("AddRootFolder() failed: %v", err)
				return
			}
			if got != tt.want {
				t.Errorf("AddRootFolder() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFolderService_AddChildFolder(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for receiver constructor.
		db *sql.DB
		// Named input parameters for target function.
		fname         string
		real_path     string
		parent_id     int64
		parent_folder string
		create_parent bool
		want          int64
		wantErr       bool
	}{
		{
			name:          "idk",
			fname:         "hellofolder",
			real_path:     "testfolder\\hellofolder",
			parent_id:     1,
			parent_folder: "testfolder",
			create_parent: true,
			want:          2,
			wantErr:       false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := setupDB()
			s := NewFolderService(db)
			if tt.create_parent {
				s.AddRootFolder(tt.parent_folder, tt.parent_folder)
			}

			got, gotErr := s.AddChildFolder(tt.fname, tt.real_path, tt.parent_id)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("AddChildFolder() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("AddChildFolder() succeeded unexpectedly")
			}
			// TODO: update the condition below to compare got with tt.want.
			if got != tt.want {
				t.Errorf("AddChildFolder() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFolderService_ParentId(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		path          string
		want          int64
		create_parent bool
		wantErr       bool
	}{
		{
			name:          "idk",
			path:          "hej",
			want:          1,
			create_parent: true,
			wantErr:       false,
		},
		{
			name:          "idk",
			path:          "hej",
			want:          1,
			create_parent: false,
			wantErr:       true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := setupDB()
			s := NewFolderService(db)
			if tt.create_parent {
				s.AddRootFolder(tt.path, tt.path)
			}
			got, gotErr := s.ParentId(tt.path)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("ParentId() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("ParentId() succeeded unexpectedly")
			}
			// TODO: update the condition below to compare got with tt.want.
			if got != tt.want {
				t.Errorf("ParentId() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFolderService_Folders(t *testing.T) {
	var newOrphanFolder = func(id int64, name, rp string) models.Folder {
		return models.Folder{
			ID: id, Name: name, RealPath: rp,
		}
	}

	var newFolder = func(id int64, name, rp string, pid int64) models.Folder {
		return models.Folder{
			ID: id, Name: name, RealPath: rp, ParentID: &pid,
		}
	}

	tests := []struct {
		name string // description of this test case
		// Named input parameters for receiver constructor.
		want    []models.Folder
		wantErr bool
	}{
		{
			name: "hello",
			want: []models.Folder{
				newOrphanFolder(1, ".", ".")},
		},
		{
			name: "atleastonesub",
			want: []models.Folder{
				newOrphanFolder(1, ".", "."),
				newFolder(2, "sub", "sub", 1),
				newFolder(3, "sub", "sub/hej", 2),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := setupDB()
			s := NewFolderService(db)
			for _, mo := range tt.want {
				if mo.ParentID == nil {
					s.AddRootFolder(mo.Name, mo.RealPath)
				} else {
					s.AddChildFolder(mo.Name, mo.RealPath, *mo.ParentID)
				}
			}
			got, gotErr := s.Folders()
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("Folders() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("Folders() succeeded unexpectedly")
			}
			// TODO: update the condition below to compare got with tt.want.
			if !slices.EqualFunc(got, tt.want, func(a, b models.Folder) bool {
				if a.ID != b.ID || a.Name != b.Name || a.RealPath != b.RealPath {
					return false
				}
				if a.ParentID == nil && b.ParentID == nil {
					return true
				}
				if a.ParentID == nil || b.ParentID == nil {
					return false
				}
				if *a.ParentID == *b.ParentID {
					return true
				}
				return true
			}) {
				t.Errorf("Folders() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFolderService_Folder(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for receiver constructor.
		// Named input parameters for target function.
		id      int64
		want    *models.Folder
		create  bool
		wantErr bool
	}{
		{
			name: "folderyes",
			id:   1,
			want: &models.Folder{
				ID:       1,
				Name:     "hej",
				RealPath: "hej",
			},
			create:  true,
			wantErr: false,
		},
		{
			name: "folderno",
			id:   1,
			want: &models.Folder{
				ID:       1,
				Name:     "hej",
				RealPath: "hej",
			},
			create:  false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := setupDB()
			s := NewFolderService(db)
			if tt.create {
				s.AddRootFolder(tt.want.Name, tt.want.RealPath)
			}
			got, gotErr := s.Folder(tt.id)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("Folder() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("Folder() succeeded unexpectedly")
			}
			// TODO: update the condition below to compare got with tt.want.
			if *got != *tt.want {
				t.Errorf("Folder() = %v, want %v", got, tt.want)
			}
		})
	}
}
