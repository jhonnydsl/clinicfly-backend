package utils

import (
	"testing"

	"golang.org/x/crypto/bcrypt"
)

func TestCheckPassword(t *testing.T) {
	hashed, err := bcrypt.GenerateFromPassword(
		[]byte("senha123"),
		bcrypt.DefaultCost,
	)
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name string
		hash string
		password string
		wantError bool
	}{
		{
			name: "correct password",
			hash: string(hashed),
			password: "senha123",
			wantError: false,
		},
		{
			name: "invalid password",
			hash: string(hashed),
			password: "senhaerrada",
			wantError: true,
		},
		{
			name: "invalid hash",
			hash: "hash-invalido",
			password: "senha123",
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := CheckPassword(tt.hash, tt.password)

			hasError := err != nil

			if hasError != tt.wantError {
				t.Errorf("case test %q failed", tt.name)
			}
		})
	}
}