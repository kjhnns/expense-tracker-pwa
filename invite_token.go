package main

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// InviteClaims represents the payload stored in an invite token.
type InviteClaims struct {
	Phone   string `json:"phone"`
	GroupID string `json:"group_id"`
	jwt.RegisteredClaims
}

// GenerateInviteToken creates a signed token granting a phone access to a group.
// The token expires after 48 hours.
func GenerateInviteToken(phone string, groupID string) (string, error) {
	secret := os.Getenv("INVITE_TOKEN_SECRET")
	if secret == "" {
		return "", errors.New("INVITE_TOKEN_SECRET not set")
	}
	claims := InviteClaims{
		Phone:   phone,
		GroupID: groupID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(48 * time.Hour)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

// ValidateInviteToken parses and validates a previously generated invite token.
// It returns the claims if the token is valid and not expired.
func ValidateInviteToken(tokenStr string) (*InviteClaims, error) {
	secret := os.Getenv("INVITE_TOKEN_SECRET")
	if secret == "" {
		return nil, errors.New("INVITE_TOKEN_SECRET not set")
	}

	parsed, err := jwt.ParseWithClaims(tokenStr, &InviteClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}
	claims, ok := parsed.Claims.(*InviteClaims)
	if !ok || !parsed.Valid {
		return nil, errors.New("invalid token")
	}
	return claims, nil
}
