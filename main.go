package main

import (
	"database/sql"
	"flag"
	"fmt"
	"io/fs"
	"log"
	"path/filepath"

	_ "modernc.org/sqlite"
)

var count = 0
var dcount = 0

func addRootFolder(db *sql.DB, name, real_path string) (int64, error) {
	res, err := db.Exec("INSERT INTO folder (name, real_path) VALUES (?, ?)", name, real_path)
	if err != nil {
		return 0, fmt.Errorf("addRootFolder: %v", err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("addRootFolder: %v", err)
	}
	return id, nil
}

func addChildFolder(db *sql.DB, name, real_path string, parent_id int64) (int64, error) {
	res, err := db.Exec("INSERT INTO folder(name, real_path, parent_id) VALUES(?, ?, ?)", name, real_path, parent_id)
	if err != nil {
		return 0, fmt.Errorf("addChildFolder: %v, %v", err, real_path)
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("addChildFolder: %v, %v", err, real_path)
	}
	return id, nil
}

func getParentId(db *sql.DB, path string) (int64, error) {
	var parent_id int64 = 0
	if err := db.QueryRow("SELECT id FROM folder WHERE real_path = ?", path).Scan(&parent_id); err != nil {
		if err == sql.ErrNoRows {
			return 0, fmt.Errorf("getParentId: %v", path)
		}
		return 0, fmt.Errorf("getParentId: %v", err)
	}
	return parent_id, nil
}

func addFile(db *sql.DB, name, real_path string, folder_id int64) (int64, error) {
	res, err := db.Exec("INSERT INTO file(name, real_path, folder_id) VALUES(?, ?, ?)", name, real_path, folder_id)
	if err != nil {
		return 0, fmt.Errorf("addFile: %v", err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("addFile: %v", err)
	}
	return id, nil
}

func main() {
	db, err := sql.Open("sqlite", "file::memory:?cache=shared")
	if err != nil {
		log.Fatal(err)
	}
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	_, err = db.Exec(`
	CREATE TABLE folder (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		real_path TEXT UNIQUE NOT NULL,
		parent_id INTEGER,
		FOREIGN KEY (parent_id) REFERENCES folder(id)
	);
	CREATE TABLE file (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		real_path TEXT UNIQUE NOT NULL,
		folder_id INTEGER NOT NULL,
		FOREIGN KEY (folder_id) REFERENCES folder(id)
	);
	`)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("hello")
	flag.Parse()
	root := flag.Arg(0)
	if len(root) == 0 {
		fmt.Println("nothing")
	}
	// fmt.Println(root)
	clroot := filepath.Clean(root)
	// fmt.Println(clroot)
	// var root_id int64 = 0
	err = filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		// fmt.Println(path)
		dir, _ := filepath.Split(path)
		cldir := filepath.Clean(dir)
		clpath := filepath.Clean(path)
		ba := filepath.Base(path)

		if d.IsDir() && clpath == clroot {
			id, err := addRootFolder(db, d.Name(), clroot)
			if err != nil {
				return err
			}
			fmt.Printf("Root:\n\tBase: %s, Cleaned: %s, ID: %d\n", ba, clroot, id)
		} else {
			parent_id, err := getParentId(db, cldir)
			// if err != nil && root_id == 0 {
			// 	return err
			// }
			if err != nil {
				return err
			}
			// if err != nil && root_id != 0 {
			// 	parent_id = root_id
			// }
			// if root_id == 0 {
			// 	root_id = parent_id
			// }
			fmt.Printf("ParentID: %d Current: %s\n", parent_id, path)
			if parent_id > 0 {
				if d.IsDir() {
					id, err := addChildFolder(db, d.Name(), clpath, parent_id)
					if err != nil {
						return err
					}
					fmt.Printf("Child folder created with ID: %v, path: %v\n", id, clpath)
				} else {
					id, err := addFile(db, d.Name(), clpath, parent_id)
					if err != nil {
						return err
					}
					fmt.Printf("File created with ID: %v, name: %v, parentID: %v\n", id, d.Name(), parent_id)
				}
			}
			if d.IsDir() {
				dcount += 1
				// fmt.Println(d.Name())
			} else {
				count += 1
			}
		}
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%d directories\n", dcount)
	fmt.Printf("%d files\n", count)
}
