package main

import (
	"log"
	"mini-alt/router"
	"mini-alt/storage"
)

func main() {
	store := storage.NewInMemoryStore()
	r := router.SetupRouter(store)

	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
