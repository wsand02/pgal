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
	database.DB().Exec(models.Schemas())
	index.Walk(root)

	http.HandleFunc("/folders/", handlers.Folders)
	http.HandleFunc("/files/", handlers.Files)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
