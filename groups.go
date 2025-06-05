package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"
)

// LegacyGroup represents the previous group model with inline phone numbers.
//
// Deprecated: this type will be removed in favour of `Group` which stores
// members in a separate table.
type LegacyGroup struct {
	ID              int64     `json:"id"`
	Phones          []string  `json:"phones"`
	Name            string    `json:"name"`
	CreatedBy       string    `json:"created_by"`
	CreatedAt       time.Time `json:"created_at"`
	DefaultCurrency string    `json:"default_currency"`
}

type createGroupRequest struct {
	Phones          []string `json:"phones"`
	Name            string   `json:"name"`
	CreatedBy       string   `json:"created_by"`
	DefaultCurrency string   `json:"default_currency"`
}

// createGroupHandler saves a new group with phone numbers to the database.
//
// Deprecated: this handler operates on the old `groups` schema that stored all
// member phone numbers in a JSON column. New code should use
// `createGroupEndpoint` instead.
func createGroupHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req createGroupRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid request", http.StatusBadRequest)
			return
		}

		phonesJSON, err := json.Marshal(req.Phones)
		if err != nil {
			http.Error(w, "server error", http.StatusInternalServerError)
			return
		}

		stmt, err := db.Prepare(`INSERT INTO groups(phone_numbers, name, created_by, default_currency) VALUES(?,?,?,?)`)
		if err != nil {
			http.Error(w, "server error", http.StatusInternalServerError)
			return
		}
		defer stmt.Close()

		res, err := stmt.Exec(string(phonesJSON), req.Name, req.CreatedBy, req.DefaultCurrency)
		if err != nil {
			http.Error(w, "server error", http.StatusInternalServerError)
			return
		}
		id, err := res.LastInsertId()
		if err != nil {
			http.Error(w, "server error", http.StatusInternalServerError)
			return
		}

		var createdAt time.Time
		err = db.QueryRow(`SELECT created_at FROM groups WHERE id = ?`, id).Scan(&createdAt)
		if err != nil {
			http.Error(w, "server error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(LegacyGroup{
			ID:              id,
			Phones:          req.Phones,
			Name:            req.Name,
			CreatedBy:       req.CreatedBy,
			CreatedAt:       createdAt,
			DefaultCurrency: req.DefaultCurrency,
		})
	}
}
