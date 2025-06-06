package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/twilio/twilio-go"
	twilioApi "github.com/twilio/twilio-go/rest/api/v2010"
)

// sendSMS sends a text message using the Twilio API.
func sendSMS(client *twilio.RestClient, to, body string) error {
	from := os.Getenv("TWILIO_FROM")
	params := &twilioApi.CreateMessageParams{}
	params.SetTo(to)
	params.SetFrom(from)
	params.SetBody(body)
	_, err := client.Api.CreateMessage(params)
	return err
}

// registerHandler initiates login by sending a magic link via SMS.
func registerHandler(db *sql.DB, client *twilio.RestClient, baseURL string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			Phone string `json:"phone"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Phone == "" {
			http.Error(w, "invalid request", http.StatusBadRequest)
			return
		}

		if _, err := db.Exec(`INSERT OR IGNORE INTO users(phone_number) VALUES(?)`, req.Phone); err != nil {
			http.Error(w, "server error", http.StatusInternalServerError)
			return
		}

		token := uuid.New().String()
		expires := time.Now().Add(15 * time.Minute)
		if _, err := db.Exec(`INSERT INTO login_tokens(token, phone_number, expires_at) VALUES(?,?,?)`, token, req.Phone, expires); err != nil {
			http.Error(w, "server error", http.StatusInternalServerError)
			return
		}

		link := fmt.Sprintf("%s/verify?token=%s", baseURL, token)
		if err := sendSMS(client, req.Phone, "Your login link: "+link); err != nil {
			log.Printf("send sms: %v", err)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"link": link})
	}
}

// verifyHandler validates a magic link token and marks the user as verified.
func verifyHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := r.URL.Query().Get("token")
		if token == "" {
			http.Error(w, "token required", http.StatusBadRequest)
			return
		}

		var phone string
		var expires time.Time
		err := db.QueryRow(`SELECT phone_number, expires_at FROM login_tokens WHERE token = ?`, token).Scan(&phone, &expires)
		if err == sql.ErrNoRows || time.Now().After(expires) {
			http.Error(w, "invalid or expired token", http.StatusBadRequest)
			return
		} else if err != nil {
			http.Error(w, "server error", http.StatusInternalServerError)
			return
		}

		tx, err := db.Begin()
		if err != nil {
			http.Error(w, "server error", http.StatusInternalServerError)
			return
		}
		defer tx.Rollback()

		if _, err := tx.Exec(`DELETE FROM login_tokens WHERE token = ?`, token); err != nil {
			http.Error(w, "server error", http.StatusInternalServerError)
			return
		}
		if _, err := tx.Exec(`UPDATE users SET verified = 1 WHERE phone_number = ?`, phone); err != nil {
			http.Error(w, "server error", http.StatusInternalServerError)
			return
		}
		if err := tx.Commit(); err != nil {
			http.Error(w, "server error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"phone": phone})
	}
}
