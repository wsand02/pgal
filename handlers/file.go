package handlers

import (
	"fmt"
	"net/http"

	"github.com/wsand02/pgal/models"
)

func Files(w http.ResponseWriter, r *http.Request) {
	files, err := models.Files()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprint(w, "<table border='1'><tr><th>ID</th><th>NAME</th><th>PATH</th></tr>")
	for _, fi := range files {
		fmt.Fprintf(w, "<tr><td>ID: %d</td><td>NAME: %s</td><td>PATH: %s</td></tr>", fi.ID, fi.Name, fi.RealPath)
	}
	fmt.Fprint(w, "</table>")
}
