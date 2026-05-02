package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/wsand02/pgal/database"
	"github.com/wsand02/pgal/handlers"
	"github.com/wsand02/pgal/index"
	"github.com/wsand02/pgal/models"
	_ "modernc.org/sqlite"
)

func main() {
	flag.Parse()
	root := flag.Arg(0)
	if len(root) == 0 {
		root = "."
	}
	db := database.GetDB()
	db.Exec(models.Schemas())

	index.Walk(root)

	http.HandleFunc("GET /folders/", handlers.Folders)
	http.HandleFunc("GET /folders/{id}", handlers.Folder)
	http.HandleFunc("GET /files/", handlers.Files)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
