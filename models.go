package main

import (
	"database/sql"
	"time"
)

// User represents an application user identified by phone number.
type User struct {
	PhoneNumber    string         `db:"phone_number" json:"phone_number"`
	Verified       bool           `db:"verified" json:"verified"`
	DisplayName    sql.NullString `db:"display_name" json:"display_name"`
	Email          sql.NullString `db:"email" json:"email"`
	NotifyBySMS    bool           `db:"notify_by_sms" json:"notify_by_sms"`
	NotifyByEmail  bool           `db:"notify_by_email" json:"notify_by_email"`
	PaymentMethods sql.NullString `db:"payment_methods" json:"payment_methods"`
}

// Group represents a collection of users who share expenses.
type Group struct {
	ID              string    `db:"id" json:"id"`
	Name            string    `db:"name" json:"name"`
	CreatedBy       string    `db:"created_by" json:"created_by"`
	DefaultCurrency string    `db:"default_currency" json:"default_currency"`
	CreatedAt       time.Time `db:"created_at" json:"created_at"`
}

// GroupMember links a user to a group with optional display name override.
type GroupMember struct {
	GroupID             string         `db:"group_id" json:"group_id"`
	PhoneNumber         string         `db:"phone_number" json:"phone_number"`
	DisplayNameOverride sql.NullString `db:"display_name_override" json:"display_name_override"`
}
