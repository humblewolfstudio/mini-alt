package main

import (
	_ "github.com/mattn/go-sqlite3"
	"log"
	"mini-alt/router"
	"mini-alt/storage"
)

func main() {
	store, err := storage.NewSQLiteStore("./mini-alt.sqlite")
	if err != nil {
		log.Fatal(err)
		return
	}

	r := router.SetupRouter(store)

	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
