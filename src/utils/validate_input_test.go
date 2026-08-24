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

func TestValidateFullName(t *testing.T) {
	fullName60 := strings.Repeat("x", 57) + " Ab"
	fullName61 := strings.Repeat("x", 58) + " Ab"
	
	tests := []struct {
		name string
		fullName string
		wantError bool
	}{
		{
			name: "valid name",
			fullName: "Jhonny Lima",
			wantError: false,
		},
		{
			name: "empty name",
			fullName: "",
			wantError: true,
		},
		{
			name: "just spaces",
			fullName: "              ",
			wantError: true,
		},
		{
			name: "4 characters",
			fullName: "abcd",
			wantError: true,
		},
		{
			name: "5 characters",
			fullName: "A Bcd",
			wantError: false,
		},
		{
			name: "60 characters",
			fullName: fullName60,
			wantError: false,
		},
		{
			name: "61 characters",
			fullName: fullName61,
			wantError: true,
		},
		{
			name: "just first name",
			fullName: "Jhonnt",
			wantError: true,
		},
		{
			name: "name with spaces at the end",
			fullName: "   Jhonny Lima   ",
			wantError: false,
		},
		{
			name: "name with multiple spaces",
			fullName: "Jhonny  Lima",
			wantError: false,
		},
		{
			name: "name with accents",
			fullName: "Jhónny Lima",
			wantError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateFullName(tt.fullName)

			hasError := err != nil

			if hasError != tt.wantError {
				t.Errorf("test case %q failed", tt.name)
			}
		})
	}
}

func TestValidatePhone(t *testing.T) {
	phone15 := strings.Repeat("1", 15)
	phone16 := strings.Repeat("1", 16)

	tests := []struct {
		name string
		phone string
		wantError bool
	}{
		{
			name: "valid phone",
			phone: "11987654321",
			wantError: false,
		},
		{
			name: "empty phone",
			phone: "",
			wantError: true,
		},
		{
			name: "just spaces",
			phone: "              ",
			wantError: true,
		},
		{
			name: "8 digits",
			phone: "12345678",
			wantError: true,
		},
		{
			name: "9 digits",
			phone: "123456789",
			wantError: false,
		},
		{
			name: "15 digits",
			phone: phone15,
			wantError: false,
		},
		{
			name: "16 digits",
			phone: phone16,
			wantError: true,
		},
		{
			name: "phone format",
			phone: "(11) 987654321",
			wantError: false,
		},
		{
			name: "with +55",
			phone: "+55 11 987654321",
			wantError: false,
		},
		{
			name: "jumbled letters",
			phone: "11abc98765432166",
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidatePhone(tt.phone)

			hasError := err != nil

			if hasError != tt.wantError {
				t.Errorf("case test %q failed", tt.name)
			}
		})
	}
}

func TestValidatePublicSlug(t *testing.T) {
	slug51 := strings.Repeat("x", 51)
	slug50 := strings.Repeat("x", 50)

	tests := []struct {
		name string
		slug string
		wantError bool
	}{
		{
			name: "valid slug",
			slug: "jhonny-lima",
			wantError: false,
		},
		{
			name: "slug with number",
			slug: "jhonny123",
			wantError: false,
		},
		{
			name: "just numebers",
			slug: "123",
			wantError: false,
		},
		{
			name: "empty slug",
			slug: "",
			wantError: false,
		},
		{
			name: "just spaces",
			slug: "     ",
			wantError: false,
		},
		{
			name: "with 2 characters",
			slug: "ab",
			wantError: true,
		},
		{
			name: "with 3 characters",
			slug: "abc",
			wantError: false,
		},
		{
			name: "with 51 characters",
			slug: slug51,
			wantError: true,
		},
		{
			name: "with 50 characters",
			slug: slug50,
			wantError: false,
		},
		{
			name: "with uppercase letter",
			slug: "Jhonny-lima",
			wantError: true,
		},
		{
			name: "with spaces",
			slug: "jhonny lima",
			wantError: true,
		},
		{
			name: "invalid character",
			slug: "-jhonny",
			wantError: true,
		},
		{
			name: "- at the end",
			slug: "jhonny-",
			wantError: true,
		},
		{
			name: "duble --",
			slug: "jhonny--lima",
			wantError: true,
		},
		{
			name: "with _",
			slug: "jhonny_lima",
			wantError: true,
		},
		{
			name: "with .",
			slug: "jhonny.lima",
			wantError: true,
		},
		{
			name: "complet with number at the end",
			slug: "jhonny-lima-123",
			wantError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidatePublicSlug(tt.slug)

			hasError := err != nil

			if hasError != tt.wantError {
				t.Errorf("case test %q failed", tt.name)
			}
		})
	}
}

func TestValidateProfileImageURL(t *testing.T) {
	prefix := "http://example.com/"

	url2048 :=  prefix + strings.Repeat("x", 2048-len(prefix))
	url2049 :=  prefix + strings.Repeat("x", 2049-len(prefix))
	tests := []struct {
		name string
		url string
		wantError bool
	}{
		{
			name: "valid http url",
			url: "http://example.com/image.jpg",
			wantError: false,
		},
		{
			name: "valid https url",
			url: "https://example.com/image.jpg",
			wantError: false,
		},
		{
			name: "empty url",
			url: "",
			wantError: false,
		},
		{
			name: "just spaces",
			url: "           ",
			wantError: false,
		},
		{
			name: "spaces at the end",
			url: "   https://example.com/image.jpg   ",
			wantError: false,
		},
		{
			name: "url with 2049 chars",
			url: url2049,
			wantError: true,
		},
		{
			name: "with 2048 chars",
			url: url2048,
			wantError: false,
		},
		{
			name: "invali url",
			url: "not-a-url",
			wantError: true,
		},
		{
			name: "without protocol",
			url: "example.com/image.jpg",
			wantError: true,
		},
		{
			name: "FTP",
			url: "ftp://example.com/image.jpg",
			wantError: true,
		},
		{
			name: "without host",
			url: "https:///image.jpg",
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateProfileImageURL(tt.url)

			hasError := err != nil

			if hasError != tt.wantError {
				t.Errorf("case test %q failed", tt.name)
			}
		})
	}
}

func TestValidateOfficeAddress(t *testing.T) {
	address255 := strings.Repeat("x", 255)
	address256 := strings.Repeat("x", 256)

	tests := []struct {
		name string
		address string
		wantError bool
	}{
		{
			name: "valid address",
			address: "Rua das Flores 123",
			wantError: false,
		},
		{
			name: "empty address",
			address: "",
			wantError: false,
		},
		{
			name: "just spaces",
			address: "          ",
			wantError: false,
		},
		{
			name: "4 characters",
			address: "abcd",
			wantError: true,
		},
		{
			name: "5 characters",
			address: "abcde",
			wantError: false,
		},
		{
			name: "255 char",
			address: address255,
			wantError: false,
		},
		{
			name: "256 char",
			address: address256,
			wantError: true,
		},
		{
			name: "with spaces",
			address: "   Rua das Flores 123   ",
			wantError: false,
		},
		{
			name: "with accents",
			address: "São Paulo",
			wantError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateOfficeAddress(tt.address)

			hasError := err != nil

			if hasError != tt.wantError {
				t.Errorf("case test %q failed", tt.name)
			}
		})
	}
}

func TestValidateCRP(t *testing.T) {
	crp50 := strings.Repeat("x", 50)
	crp51 := strings.Repeat("x", 51)

	tests := []struct {
		name string
		crp string
		wantError bool
	}{
		{
			name: "valid crp",
			crp: "CRP 06/123456",
			wantError: false,
		},
		{
			name: "empty crp",
			crp: "",
			wantError: false,
		},
		{
			name: "just spaces",
			crp: "           ",
			wantError: false,
		},
		{
			name: "50 characters",
			crp: crp50,
			wantError: false,
		},
		{
			name: "51 characters",
			crp: crp51,
			wantError: true,
		},
		{
			name: "with newline",
			crp: "CRP\n06/123456",
			wantError: true,
		},
		{
			name: "with tab",
			crp: "CRP\t06/123456",
			wantError: true,
		},
		{
			name: "with carriage return",
			crp: "CRP\r06/123456",
			wantError: true,
		},
		{
			name: "spaces at the end",
			crp: "   CRP 06/123456   ",
			wantError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateCRP(tt.crp)

			hasError := err != nil

			if hasError != tt.wantError {
				t.Errorf("case test %q failed", tt.name)
			}
		})
	}
}