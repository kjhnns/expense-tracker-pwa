package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func main() {
	r := chi.NewRouter()

	// Serve the frontend index.html when the root path is requested
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./frontend/index.html")
	})

	// Serve static frontend assets such as app.js and component files
	fileServer := http.FileServer(http.Dir("./frontend"))
	r.Handle("/*", http.StripPrefix("/", fileServer))

	http.ListenAndServe(":8080", r)
}
