package models

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
