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
