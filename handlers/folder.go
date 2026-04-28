package handlers

import (
	"fmt"
	"net/http"
	"strconv"

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
		fmt.Fprintf(w, "<tr><td>ID: %d</td><td>NAME:<a href=\"%s\"> %s</a></td><td>PATH: %s</td></tr>", fol.ID, fol.URL(), fol.Name, fol.RealPath)
	}
	fmt.Fprint(w, "</table>")
}

func (fh *FolderHandler) Folder(w http.ResponseWriter, r *http.Request) {
	idString := r.PathValue("id")
	id, err := strconv.ParseInt(idString, 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	folder, err := fh.folderService.Folder(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	children, err := fh.folderService.FoldersByParent(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "<h1>%s</h1>", folder.Name)
	for _, chi := range children {
		fmt.Fprintf(w, "<p>%s</p>", chi.Name)
	}
}
