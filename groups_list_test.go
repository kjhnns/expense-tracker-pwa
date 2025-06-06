package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
)

func TestListGroupsHandler(t *testing.T) {
	db := initDB(":memory:")
	defer db.Close()

	// insert sample data
	_, err := db.Exec(`INSERT INTO users(phone_number, verified) VALUES('+111', 1),('+222', 1)`)
	if err != nil {
		t.Fatalf("insert users: %v", err)
	}
	groupID := "g1"
	_, err = db.Exec(`INSERT INTO groups(id, name, created_by, default_currency) VALUES(?,?,?,?)`, groupID, "Test", "+111", "EUR")
	if err != nil {
		t.Fatalf("insert group: %v", err)
	}
	_, err = db.Exec(`INSERT INTO group_members(group_id, phone_number) VALUES(?, ?), (?, ?)`, groupID, "+111", groupID, "+222")
	if err != nil {
		t.Fatalf("insert members: %v", err)
	}

	r := chi.NewRouter()
	r.Get("/groups", listGroupsHandler(db))

	req := httptest.NewRequest(http.MethodGet, "/groups?phone=%2B111", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 got %d", w.Code)
	}

	var groups []Group
	if err := json.NewDecoder(w.Body).Decode(&groups); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if len(groups) != 1 || groups[0].ID != groupID {
		t.Fatalf("unexpected groups %+v", groups)
	}
}
