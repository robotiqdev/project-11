package main

import (
	"log"
	"net/http"

	"github.com/workspace/repo/internal/router"
	"github.com/workspace/repo/internal/storage"
)

func main() {
	store := storage.NewInMemoryTaskStore()
	handler := router.New(store)
	log.Println("listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", handler))
}
