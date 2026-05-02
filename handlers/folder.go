package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/wsand02/pgal/models"
)

func Folders(w http.ResponseWriter, r *http.Request) {
	folders, err := models.Folders()
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

func Folder(w http.ResponseWriter, r *http.Request) {
	idString := r.PathValue("id")
	id, err := strconv.ParseInt(idString, 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	folder, err := models.GetFolder(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	children, err := models.FoldersByParent(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "<h1>%s</h1>", folder.Name)
	for _, chi := range children {
		fmt.Fprintf(w, "<p>%s</p>", chi.Name)
	}
}
