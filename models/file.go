package models

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
