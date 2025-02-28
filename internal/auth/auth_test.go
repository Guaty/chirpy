package auth

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestCheckPasswordHash(t *testing.T) {
	// First, we need to create some hashed passwords for testing
	password1 := "correctPassword123!"
	password2 := "anotherPassword456!"
	hash1, _ := HashPassword(password1)
	hash2, _ := HashPassword(password2)

	tests := []struct {
		name     string
		password string
		hash     string
		wantErr  bool
	}{
		{
			name:     "Correct password",
			password: password1,
			hash:     hash1,
			wantErr:  false,
		},
		{
			name:     "Incorrect password",
			password: "wrongPassword",
			hash:     hash1,
			wantErr:  true,
		},
		{
			name:     "Password doesn't match different hash",
			password: password1,
			hash:     hash2,
			wantErr:  true,
		},
		{
			name:     "Empty password",
			password: "",
			hash:     hash1,
			wantErr:  true,
		},
		{
			name:     "Invalid hash",
			password: password1,
			hash:     "invalidhash",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := CheckPasswordHash(tt.password, tt.hash)
			if (err != nil) != tt.wantErr {
				t.Errorf("CheckPasswordHash() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestJWT(t *testing.T) {
	// Let's get some initial variables set
	userID := uuid.New()
	tokenSecret := "mySecretKey42"
	validExpirationTime := 10 * time.Minute
	invalidExpirationTime := -2 * time.Minute

	tokenString1, _ := MakeJWT(userID, tokenSecret, validExpirationTime)
	tokenString2, _ := MakeJWT(userID, tokenSecret, invalidExpirationTime)

	tests := []struct {
		name        string
		id          uuid.UUID
		tokenString string
		tokenSecret string
		wantErr     bool
	}{
		{
			name:        "Valid Token",
			id:          userID,
			tokenString: tokenString1,
			tokenSecret: tokenSecret,
			wantErr:     false,
		},
		{
			name:        "Expired Token",
			id:          userID,
			tokenString: tokenString2,
			tokenSecret: tokenSecret,
			wantErr:     true,
		},
		{
			name:        "Wrong secret",
			id:          userID,
			tokenString: tokenString1,
			tokenSecret: "NotMySecret",
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tokenId, err := ValidateJWT(tt.tokenString, tt.tokenSecret)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateJWT() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && tokenId != tt.id {
				t.Errorf("Expected userID %v, got %v", tt.id, tokenId)
			}
		})
	}
}
