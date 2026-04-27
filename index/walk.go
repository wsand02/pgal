package index

import (
	"fmt"
	"io/fs"
	"path/filepath"

	"github.com/wsand02/pgal/models"
)

func Walk(root string) error {
	clroot := filepath.Clean(root)
	var count = 0
	var dcount = 0
	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		// fmt.Println(path)
		dir, _ := filepath.Split(path)
		cldir := filepath.Clean(dir)
		clpath := filepath.Clean(path)
		ba := filepath.Base(path)

		if d.IsDir() && clpath == clroot {
			id, err := models.AddRootFolder(d.Name(), clroot)
			if err != nil {
				return err
			}
			fmt.Printf("Root:\n\tBase: %s, Cleaned: %s, ID: %d\n", ba, clroot, id)
		} else {
			parent_id, err := models.ParentId(cldir)
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
					id, err := models.AddChildFolder(d.Name(), clpath, parent_id)
					if err != nil {
						return err
					}
					fmt.Printf("Child folder created with ID: %v, path: %v\n", id, clpath)
				} else {
					id, err := models.AddFile(d.Name(), clpath, parent_id)
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
	fmt.Printf("%d directories\n", dcount)
	fmt.Printf("%d files\n", count)
	if err != nil {
		return err
	}
	return nil
}
