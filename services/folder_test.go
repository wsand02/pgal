package services

import (
	"database/sql"
	"slices"
	"testing"

	"github.com/wsand02/pgal/models"
	_ "modernc.org/sqlite"
)

const (
	tfoldroot = "."
	tfoldsub1 = "sub"
	tfoldsub2 = "sub/hej"
)

func setupTestDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		return nil, err
	}
	_, err = db.Exec(models.FolderSchema())
	if err != nil {
		return nil, err
	}
	return db, nil
}

func newTestFolderService(t *testing.T) *FolderService {
	db, err := setupTestDB()
	if err != nil {
		t.Fatal(err)
	}
	return NewFolderService(db)
}

func TestFolderService_AddRootFolder(t *testing.T) {
	tests := []struct {
		name     string
		fname    string
		realPath string
		want     int64
	}{
		{"dot", ".", ".", 1},
		{"testfolder", "testfolder", "testfolder", 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := newTestFolderService(t)
			got, err := s.AddRootFolder(tt.fname, tt.realPath)

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
		name         string
		fname        string
		realPath     string
		parentID     int64
		createParent bool
		want         int64
		wantErr      bool
	}{
		{"hellofolder", "hellofolder", tfoldsub1, 1, true, 2, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := newTestFolderService(t)

			if tt.createParent {
				s.AddRootFolder(tfoldroot, tfoldroot)
			}

			got, gotErr := s.AddChildFolder(tt.fname, tt.realPath, tt.parentID)

			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("AddChildFolder() failed: %v", gotErr)
				}
				return
			}

			if tt.wantErr {
				t.Fatal("AddChildFolder() succeeded unexpectedly")
			}

			if got != tt.want {
				t.Errorf("AddChildFolder() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFolderService_Folders(t *testing.T) {
	tests := []struct {
		name string
		want []models.Folder
	}{
		{"hello", []models.Folder{
			models.NewOrphanFolder(1, tfoldroot, tfoldroot),
		}},
		{"atleastonesub", []models.Folder{
			models.NewOrphanFolder(1, tfoldroot, tfoldroot),
			models.NewFolder(2, tfoldsub1, tfoldsub1, 1),
			models.NewFolder(3, tfoldsub2, tfoldsub2, 2),
		}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := newTestFolderService(t)

			// Set up the expected state
			for _, mo := range tt.want {
				if mo.ParentID == nil {
					s.AddRootFolder(mo.Name, mo.RealPath)
				} else {
					s.AddChildFolder(mo.Name, mo.RealPath, *mo.ParentID)
				}
			}

			got, gotErr := s.Folders()

			if gotErr != nil {
				t.Errorf("Folders() failed: %v", gotErr)
				return
			}

			// Compare results
			if !slices.EqualFunc(got, tt.want, models.FoldersEqual) {
				t.Errorf("Folders() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFolderService_Folder(t *testing.T) {
	tests := []struct {
		name    string
		id      int64
		want    *models.Folder
		create  bool
		wantErr bool
	}{
		{"folderyes", 1, &models.Folder{ID: 1, Name: "hej", RealPath: "hej"}, true, false},
		{"folderno", 1, &models.Folder{ID: 1, Name: "hej", RealPath: "hej"}, false, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := newTestFolderService(t)

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

			// Compare results (handle pointer comparison carefully)
			if got == nil || *got != *tt.want {
				t.Errorf("Folder() = %v, want %v", got, tt.want)
			}
		})
	}
}
