package main

import (
	"log"
	"net/http"
	"time"

	"github.com/workspace/repo/internal/router"
	"github.com/workspace/repo/internal/storage"
)

func main() {
	store := storage.NewInMemoryTaskStore()
	server := &http.Server{
		Addr:         ":8080",
		Handler:      router.New(store),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}
	log.Println("listening on :8080")
	log.Fatal(server.ListenAndServe())
}
