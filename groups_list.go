package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
)

// listGroupsHandler returns all groups a phone number belongs to.
func listGroupsHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		phone := r.URL.Query().Get("phone")
		if phone == "" {
			http.Error(w, "phone required", http.StatusBadRequest)
			return
		}
		rows, err := db.Query(`SELECT g.id, g.name, g.created_by, g.default_currency, g.created_at
            FROM groups g
            JOIN group_members gm ON g.id = gm.group_id
            WHERE gm.phone_number = ?`, phone)
		if err != nil {
			http.Error(w, "server error", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var groups []Group
		for rows.Next() {
			var g Group
			if err := rows.Scan(&g.ID, &g.Name, &g.CreatedBy, &g.DefaultCurrency, &g.CreatedAt); err != nil {
				http.Error(w, "server error", http.StatusInternalServerError)
				return
			}
			groups = append(groups, g)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(groups)
	}
}
