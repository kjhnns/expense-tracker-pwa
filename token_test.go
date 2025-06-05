package main

import (
	"os"
	"testing"
	"time"
)

func TestGenerateAndValidateInviteToken(t *testing.T) {
	os.Setenv("INVITE_TOKEN_SECRET", "secret123")
	defer os.Unsetenv("INVITE_TOKEN_SECRET")

	token, err := GenerateInviteToken("+41791234567", "abc123")
	if err != nil {
		t.Fatalf("generate token: %v", err)
	}

	claims, err := ValidateInviteToken(token)
	if err != nil {
		t.Fatalf("validate token: %v", err)
	}

	if claims.Phone != "+41791234567" || claims.GroupID != "abc123" {
		t.Fatalf("unexpected claims: %#v", claims)
	}

	if time.Until(claims.ExpiresAt.Time) < 47*time.Hour {
		t.Fatalf("expiry too short: %v", claims.ExpiresAt.Time)
	}
}
