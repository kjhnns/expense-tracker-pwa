package main

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
)

// createGroupRequestV2 is the expected payload for /groups/create
// It mirrors the product requirements document.
type createGroupRequestV2 struct {
	GroupName       string   `json:"group_name"`
	DefaultCurrency string   `json:"default_currency"`
	CreatedBy       string   `json:"created_by"`
	Participants    []string `json:"participants"`
}

type createGroupResponse struct {
	GroupID     string            `json:"group_id"`
	InviteLinks map[string]string `json:"invite_links"`
}

// createGroupEndpoint creates a new group and adds participants.
func createGroupEndpoint(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req createGroupRequestV2
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid request", http.StatusBadRequest)
			return
		}

		tx, err := db.Begin()
		if err != nil {
			http.Error(w, "server error", http.StatusInternalServerError)
			return
		}
		defer tx.Rollback()

		groupID := uuid.New().String()
		if _, err := tx.Exec(`INSERT INTO groups(id, name, created_by, default_currency) VALUES(?,?,?,?)`,
			groupID, req.GroupName, req.CreatedBy, req.DefaultCurrency); err != nil {
			http.Error(w, "server error", http.StatusInternalServerError)
			return
		}

		inviteLinks := make(map[string]string)

		for _, phone := range req.Participants {
			verified := 0
			if phone == req.CreatedBy {
				verified = 1
			}
			if _, err := tx.Exec(`INSERT OR IGNORE INTO users(phone_number, verified, notify_by_sms, notify_by_email) VALUES(?, ?, 1, 1)`,
				phone, verified); err != nil {
				http.Error(w, "server error", http.StatusInternalServerError)
				return
			}

			if phone == req.CreatedBy {
				// ensure verified flag true if user already existed
				if _, err := tx.Exec(`UPDATE users SET verified = 1 WHERE phone_number = ?`, phone); err != nil {
					http.Error(w, "server error", http.StatusInternalServerError)
					return
				}
			}

			if _, err := tx.Exec(`INSERT OR IGNORE INTO group_members(group_id, phone_number) VALUES(?, ?)`,
				groupID, phone); err != nil {
				http.Error(w, "server error", http.StatusInternalServerError)
				return
			}
			token, err := GenerateInviteToken(phone, groupID)
			if err != nil {
				http.Error(w, "server error", http.StatusInternalServerError)
				return
			}
			inviteLinks[phone] = "https://app.com/invite?token=" + token
		}

		if err := tx.Commit(); err != nil {
			http.Error(w, "server error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(createGroupResponse{
			GroupID:     groupID,
			InviteLinks: inviteLinks,
		})
	}
}
