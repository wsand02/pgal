package services

import (
	"database/sql"
	"slices"
	"testing"

	"github.com/wsand02/pgal/database"
	"github.com/wsand02/pgal/models"
)

func setupFileTestDB() (*sql.DB, error) {
	db, err := database.SetupDB(":memory:")
	if err != nil {
		return nil, err
	}
	_, err = db.Exec(models.FolderSchema())
	if err != nil {
		return nil, err
	}
	_, err = db.Exec(models.FileSchema())
	if err != nil {
		return nil, err
	}
	return db, nil
}

func newTestFileService(t *testing.T) (*FolderService, *FileService) {
	db, err := setupFileTestDB()
	if err != nil {
		t.Fatal(err)
	}
	return NewFolderService(db), NewFileService(db)
}

func TestFileService_AddFile(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		fname         string
		real_path     string
		folder_id     int64
		create_folder bool
		want          int64
		wantErr       bool
	}{
		{
			name:          "hello",
			fname:         "hello.txt",
			real_path:     "hello.txt",
			folder_id:     1,
			create_folder: true,
			want:          1,
			wantErr:       false,
		},
		{
			name:          "fail pleasehello",
			fname:         "hello.txt",
			real_path:     "hello.txt",
			folder_id:     1,
			create_folder: false,
			want:          1,
			wantErr:       true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.fname, func(t *testing.T) {
			fo, fi := newTestFileService(t)
			if tt.create_folder {
				fo.AddRootFolder(tfoldroot, tfoldroot)
			}
			got, gotErr := fi.AddFile(tt.fname, tt.real_path, tt.folder_id)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("AddFile() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("AddFile() succeeded unexpectedly")
			}
			// TODO: update the condition below to compare got with tt.want.
			if got != tt.want {
				t.Errorf("AddFile() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFileService_Files(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for receiver constructor.
		create_folder bool
		create        bool
		want          []models.File
		wantErr       bool
	}{
		{
			name: "hello",
			want: []models.File{
				models.NewFile(1, "hej.txt", "hej.txt", 1),
				models.NewFile(2, "lmao.go", "lmao.go", 1),
			},
			create_folder: true,
			create:        true,
			wantErr:       false,
		},
		{
			name:          "empty",
			want:          []models.File{},
			create_folder: false,
			create:        false,
			wantErr:       false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fo, fi := newTestFileService(t)
			if tt.create_folder {
				fo.AddRootFolder(tfoldroot, tfoldroot)
			}
			if tt.create {
				for _, tfi := range tt.want {
					fi.AddFile(tfi.Name, tfi.RealPath, tfi.FolderID)
				}
			}
			got, gotErr := fi.Files()
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("Files() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("Files() succeeded unexpectedly")
			}
			// TODO: update the condition below to compare got with tt.want.
			if !slices.EqualFunc(got, tt.want, models.FileEqual) {
				t.Errorf("Files() = %v, want %v", got, tt.want)
			}
		})
	}
}
