package utils

import (
	"strings"
	"testing"
	"time"
)

func TestValidateBio(t *testing.T) {
	bio501 := strings.Repeat("x", 501)
	bio500 := strings.Repeat("x", 500)
	bio499 := strings.Repeat("x", 499)

	tests := []struct {
		name      string
		bio       string
		wantError bool
	}{
		{
			name:      "valid bio",
			bio:       "Pediatra a 10 anos",
			wantError: false,
		},
		{
			name:      "empty bio",
			bio:       "",
			wantError: false,
		},
		{
			name:      "short bio",
			bio:       "abc",
			wantError: true,
		},
		{
			name:      "bio with 501 characters",
			bio:       bio501,
			wantError: true,
		},
		{
			name:      "bio with 500 characters",
			bio:       bio500,
			wantError: false,
		},
		{
			name:      "bio with 499 characters",
			bio:       bio499,
			wantError: false,
		},
		{
			name:      "bio just with space",
			bio:       "           ",
			wantError: false,
		},
		{
			name:      "valid bio with spaces at the ends",
			bio:       "   Pediatra a 10 anos   ",
			wantError: false,
		},
		{
			name:      "bio with 9 characters",
			bio:       "123456789",
			wantError: true,
		},
		{
			name:      "bio with 10 characters",
			bio:       "1234567890",
			wantError: false,
		},
		{
			name:      "unicode",
			bio:       "Olá mundo!",
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
		name      string
		email     string
		wantError bool
	}{
		{
			name:      "valid email",
			email:     "jhonny@email.com",
			wantError: false,
		},
		{
			name:      "void email",
			email:     "",
			wantError: true,
		},
		{
			name:      "just spaces",
			email:     "              ",
			wantError: true,
		},
		{
			name:      "email with spaces at the end",
			email:     "   jhonny@gmail.com   ",
			wantError: false,
		},
		{
			name:      "invalid email",
			email:     "1234567890",
			wantError: true,
		},
		{
			name:      "email without @",
			email:     "jhonnygmail.com",
			wantError: true,
		},
		{
			name:      "email with 255 characters",
			email:     email255,
			wantError: true,
		},
		{
			name:      "email with 254 characters",
			email:     email254,
			wantError: false,
		},
		{
			name:      "display name",
			email:     "Jhonny <jhonny@gmail.com>",
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

func TestValidateBirthDate(t *testing.T) {
	now := time.Now()

	exactly18 := time.Date(
		now.Year()-18,
		now.Month(),
		now.Day(),
		0, 0, 0, 0,
		time.Local,
	).Format("2006-01-02")

	tests := []struct {
		name      string
		date      string
		wantError bool
	}{
		{
			name:      "valid date",
			date:      "2000-06-10",
			wantError: false,
		},
		{
			name:      "empty date",
			date:      "",
			wantError: true,
		},
		{
			name:      "just spaces",
			date:      "           ",
			wantError: true,
		},
		{
			name:      "valid date with spaces",
			date:      "   2003-02-03    ",
			wantError: false,
		},
		{
			name:      "invalid format",
			date:      "03-02-2005",
			wantError: true,
		},
		{
			name:      "impossible date",
			date:      "2000-99-99",
			wantError: true,
		},
		{
			name:      "future date",
			date:      "2099-01-01",
			wantError: true,
		},
		{
			name:      "under 18 years of age",
			date:      "2010-01-01",
			wantError: true,
		},
		{
			name:      "exactly 18 years",
			date:      exactly18,
			wantError: false,
		},
		{
			name:      "date 18+",
			date:      "1960-01-01",
			wantError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateBirthDate(tt.date)

			hasError := err != nil

			if hasError != tt.wantError {
				t.Errorf("test case %q failed: got error %v", tt.name, err)
			}
		})
	}
}

func TestCalculateAge(t *testing.T) {
	now := time.Now()

	exactly18 := time.Date(
		now.Year()-18,
		now.Month(),
		now.Day(),
		0, 0, 0, 0,
		time.Local,
	)

	over18 := time.Date(
		now.Year()-20,
		now.Month(),
		now.Day(),
		0, 0, 0, 0,
		time.Local,
	)

	beforeBirthDay := time.Date(
		now.Year()-20,
		now.Month(),
		now.Day()+1,
		0, 0, 0, 0,
		time.Local,
	)

	afterBirthDay := time.Date(
		now.Year()-20,
		now.Month(),
		now.Day()-1,
		0, 0, 0, 0,
		time.Local,
	)

	tests := []struct {
		name    string
		date    time.Time
		wantAge int
	}{
		{
			name:    "exactly 18 years",
			date:    exactly18,
			wantAge: 18,
		},
		{
			name:    "20 years old",
			date:    over18,
			wantAge: 20,
		},
		{
			name:    "before birth day",
			date:    beforeBirthDay,
			wantAge: 19,
		},
		{
			name:    "after birth day",
			date:    afterBirthDay,
			wantAge: 20,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := calculateAge(tt.date)

			if got != tt.wantAge {
				t.Errorf("expected age %d, got %d", tt.wantAge, got)
			}
		})
	}
}
