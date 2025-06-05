package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
)

func setupTestDB(t *testing.T) *sql.DB {
	t.Helper()
	db := initDB(":memory:")
	return db
}

func TestCreateGroupEndpoint(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	r := chi.NewRouter()
	r.Post("/groups/create", createGroupEndpoint(db))

	payload := `{"group_name":"Trip to Berlin","default_currency":"EUR","created_by":"+41791234567","participants":["+41791234567","+49123456789"]}`
	req := httptest.NewRequest(http.MethodPost, "/groups/create", bytes.NewBufferString(payload))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}

	var resp createGroupResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	if resp.GroupID == "" {
		t.Fatalf("expected group_id in response")
	}
	if len(resp.InviteLinks) != 2 {
		t.Fatalf("expected 2 invite links, got %d", len(resp.InviteLinks))
	}

	// verify group record
	var name, createdBy, currency string
	err := db.QueryRow("SELECT name, created_by, default_currency FROM groups WHERE id = ?", resp.GroupID).Scan(&name, &createdBy, &currency)
	if err != nil {
		t.Fatalf("query group: %v", err)
	}
	if name != "Trip to Berlin" || createdBy != "+41791234567" || currency != "EUR" {
		t.Fatalf("unexpected group data: %s %s %s", name, createdBy, currency)
	}

	// verify group members and user verification status
	phones := []string{"+41791234567", "+49123456789"}
	for i, phone := range phones {
		var count int
		err = db.QueryRow("SELECT COUNT(*) FROM group_members WHERE group_id = ? AND phone_number = ?", resp.GroupID, phone).Scan(&count)
		if err != nil {
			t.Fatalf("query member: %v", err)
		}
		if count != 1 {
			t.Fatalf("member %d not inserted", i)
		}

		var verified bool
		err = db.QueryRow("SELECT verified FROM users WHERE phone_number = ?", phone).Scan(&verified)
		if err != nil {
			t.Fatalf("query user: %v", err)
		}
		expected := false
		if phone == "+41791234567" {
			expected = true
		}
		if verified != expected {
			t.Fatalf("unexpected verified for %s: %v", phone, verified)
		}
	}
}
