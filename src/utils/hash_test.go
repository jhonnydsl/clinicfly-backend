package utils

import (
	"testing"

	"golang.org/x/crypto/bcrypt"
)

func TestHashedPassword(t *testing.T) {
	tests := []struct {
		name string
		password string
	}{
		{
			name: "valid password",
			password: "senha123",
		},
		{
			name: "empty password",
			password: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hashed, err := HashPassword(tt.password)

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if hashed == "" {
				t.Errorf("expected hash, got empty string")
			}

			err = bcrypt.CompareHashAndPassword(
				[]byte(hashed),
				[]byte(tt.password),
			)

			if err != nil {
				t.Errorf("generated hash does not match password")
			}
		})
	}
}