package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/wsand02/pgal/database"
	"github.com/wsand02/pgal/handlers"
	"github.com/wsand02/pgal/index"
	"github.com/wsand02/pgal/models"
	"github.com/wsand02/pgal/services"
	_ "modernc.org/sqlite"
)

func main() {
	flag.Parse()
	root := flag.Arg(0)
	if len(root) == 0 {
		root = "."
	}
	db, err := database.SetupDB("file::memory:?cache=shared")
	if err != nil {
		log.Fatal(err.Error())
	}
	db.Exec(models.Schemas())

	folderService := services.NewFolderService(db)
	fileService := services.NewFileService(db)

	walker := index.NewWalker(fileService, folderService)
	walker.Walk(root)

	folderHandler := handlers.NewFolderHandler(folderService)
	fileHandler := handlers.NewFileHandler(fileService)

	http.HandleFunc("GET /folders/", folderHandler.Folders)
	http.HandleFunc("GET /folders/{id}", folderHandler.Folder)
	http.HandleFunc("GET /files/", fileHandler.Files)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
