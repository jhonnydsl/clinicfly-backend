package utils

import (
	"strings"
	"testing"
)

func TestValidateBio(t *testing.T) {
	bio501 := strings.Repeat("x", 501)
	bio500 := strings.Repeat("x", 500)
	bio499 := strings.Repeat("x", 499)

	tests := []struct {
		name string
		bio string
		wantError bool
	}{
		{
			name: "valid bio",
			bio: "Pediatra a 10 anos",
			wantError: false,
		},
		{
			name: "empty bio",
			bio: "",
			wantError: false,
		},
		{
			name: "short bio",
			bio: "abc",
			wantError: true,
		},
		{
			name: "bio with 501 characters",
			bio: bio501,
			wantError: true,
		},
		{
			name: "bio with 500 characters",
			bio: bio500,
			wantError: false,
		},
		{
			name: "bio with 499 characters",
			bio: bio499,
			wantError: false,
		},
		{
			name: "bio just with space",
			bio: "           ",
			wantError: false,
		},
		{
			name: "valid bio with spaces at the ends",
			bio: "   Pediatra a 10 anos   ",
			wantError: false,
		},
		{
			name: "bio with 9 characters",
			bio: "123456789",
			wantError: true,
		},
		{
			name: "bio with 10 characters",
			bio: "1234567890",
			wantError: false,
		},
		{
			name: "unicode",
			bio: "Olá mundo!",
			wantError: false,
		},

	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateBio(tt.bio)

			hasError := err != nil

			if hasError != tt.wantError {
				t.Errorf("test case %q failed", tt.name)
			}
		})
	}
}

func TestValidateEmail(t *testing.T) {
	email255 := strings.Repeat("a", 245) + "@gmail.com"
	email254 := strings.Repeat("a", 244) + "@gmail.com"

	tests := []struct {
		name string
		email string
		wantError bool
	}{
		{
			name: "valid email",
			email: "jhonny@email.com",
			wantError: false,
		},
		{
			name: "void email",
			email: "",
			wantError: true,
		},
		{
			name: "just spaces",
			email: "              ",
			wantError: true,
		},
		{
			name: "email with spaces at the end",
			email: "   jhonny@gmail.com   ",
			wantError: false,
		},
		{
			name: "invalid email",
			email: "1234567890",
			wantError: true,
		},
		{
			name: "email without @",
			email: "jhonnygmail.com",
			wantError: true,
		},
		{
			name: "email with 255 characters",
			email: email255,
			wantError: true,
		},
		{
			name: "email with 254 characters",
			email: email254,
			wantError: false,
		},
		{
			name: "display name",
			email: "Jhonny <jhonny@gmail.com>",
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateEmail(tt.email)

			hasError := err != nil

			if hasError != tt.wantError {
				t.Errorf("test case %q failed", tt.name)
			}
		})
	}
}