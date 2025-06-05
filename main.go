package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func main() {
	db := initDB("data.db")
	defer db.Close()

	r := chi.NewRouter()

	// Serve the frontend index.html when the root path is requested
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./frontend/index.html")
	})

	// Serve static frontend assets such as app.js and component files
	fileServer := http.FileServer(http.Dir("./frontend"))
	r.Handle("/*", http.StripPrefix("/", fileServer))

	r.Post("/api/groups", createGroupHandler(db))

	log.Fatal(http.ListenAndServe(":8080", r))
}
