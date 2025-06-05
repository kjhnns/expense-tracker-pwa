package main

import (
	"io/fs"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
)

func TestRootServesIndex(t *testing.T) {
	frontendFS, err := fs.Sub(embeddedFrontend, "frontend")
	if err != nil {
		t.Fatalf("sub fs: %v", err)
	}

	r := chi.NewRouter()
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		data, err := fs.ReadFile(frontendFS, "index.html")
		if err != nil {
			http.Error(w, "index file not found", http.StatusInternalServerError)
			return
		}
		w.Write(data)
	})
	fileServer := http.FileServer(http.FS(frontendFS))
	r.Handle("/*", http.StripPrefix("/", fileServer))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}
	if !strings.Contains(w.Body.String(), "Welcome to the Zero-Registration Expense Tracker") {
		t.Fatalf("unexpected body: %s", w.Body.String())
	}
}
