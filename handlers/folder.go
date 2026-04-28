package handlers

import (
	"fmt"
	"net/http"

	"github.com/wsand02/pgal/services"
)

type FolderHandler struct {
	folderService *services.FolderService
}

func NewFolderHandler(fs *services.FolderService) *FolderHandler {
	return &FolderHandler{
		folderService: fs,
	}
}

func (fh *FolderHandler) Folders(w http.ResponseWriter, r *http.Request) {
	folders, err := fh.folderService.Folders()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprint(w, "<table border='1'><tr><th>ID</th><th>NAME</th><th>PATH</th></tr>")
	for _, fol := range folders {
		fmt.Fprintf(w, "<tr><td>ID: %d</td><td>NAME: %s</td><td>PATH: %s</td></tr>", fol.ID, fol.Name, fol.RealPath)
	}
	fmt.Fprint(w, "</table>")
}
