package main

import (
	"embed"
	"io/fs"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/twilio/twilio-go"
)

//go:embed frontend
var embeddedFrontend embed.FS

func main() {
	db := initDB("data.db")
	defer db.Close()

	twilioClient := twilio.NewRestClient()
	baseURL := os.Getenv("BASE_URL")
	if baseURL == "" {
		baseURL = "http://localhost:8080"
	}

	r := chi.NewRouter()

	// Use embedded frontend files
	frontendFS, err := fs.Sub(embeddedFrontend, "frontend")
	if err != nil {
		log.Fatalf("sub fs: %v", err)
	}

	// Serve the frontend index.html when the root path is requested
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		data, err := fs.ReadFile(frontendFS, "index.html")
		if err != nil {
			http.Error(w, "index file not found", http.StatusInternalServerError)
			return
		}
		w.Write(data)
	})

	// Serve static frontend assets such as app.js and component files
	fileServer := http.FileServer(http.FS(frontendFS))
	r.Handle("/*", http.StripPrefix("/", fileServer))

	// Auth endpoints
	r.Post("/register", registerHandler(db, twilioClient, baseURL))
	r.Get("/verify", verifyHandler(db))

	// Deprecated: the /api/groups endpoint used an older schema and has
	// been replaced by /groups/create.
	r.Post("/groups/create", createGroupEndpoint(db))

	log.Fatal(http.ListenAndServe(":8080", r))
}
