package index

import (
	"fmt"
	"io/fs"
	"path/filepath"

	"github.com/wsand02/pgal/services"
)

type Walker struct {
	fileService   *services.FileService
	folderService *services.FolderService
}

func NewWalker(fis *services.FileService, fos *services.FolderService) *Walker {
	return &Walker{fileService: fis, folderService: fos}
}

func (w *Walker) Walk(root string) error {
	clroot := filepath.Clean(root)
	var count = 0
	var dcount = 0
	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		dir, _ := filepath.Split(path)
		cldir := filepath.Clean(dir)
		clpath := filepath.Clean(path)
		ba := filepath.Base(path)

		if d.IsDir() && clpath == clroot {
			id, err := w.folderService.AddRootFolder(d.Name(), clroot)
			if err != nil {
				return err
			}
			fmt.Printf("Root:\n\tBase: %s, Cleaned: %s, ID: %d\n", ba, clroot, id)
		} else {
			parent_id, err := w.folderService.ParentId(cldir)
			if err != nil {
				return err
			}
			fmt.Printf("ParentID: %d Current: %s\n", parent_id, path)
			if parent_id > 0 {
				if d.IsDir() {
					id, err := w.folderService.AddChildFolder(d.Name(), clpath, parent_id)
					if err != nil {
						return err
					}
					fmt.Printf("Child folder created with ID: %v, path: %v\n", id, clpath)
				} else {
					id, err := w.fileService.AddFile(d.Name(), clpath, parent_id)
					if err != nil {
						return err
					}
					fmt.Printf("File created with ID: %v, name: %v, parentID: %v\n", id, d.Name(), parent_id)
				}
			}
			if d.IsDir() {
				dcount += 1
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
