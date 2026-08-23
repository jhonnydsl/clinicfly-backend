package utils

import (
	"fmt"
	"net/mail"
	"net/url"
	"regexp"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/jhonnydsl/clinify-backend/src/dtos"
)

func ValidateAdminInput(admin dtos.AdminInput) error {
	fullname := strings.TrimSpace(admin.FullName)

	if utf8.RuneCountInString(fullname) < 5 {
		return fmt.Errorf("name must be at least 5 characters long")
	}

	parts := strings.Fields(fullname)
	if len(parts) < 2 {
		return fmt.Errorf("full name must include first and last name")
	}

	_, err := mail.ParseAddress(admin.Email)
	if err != nil {
		return fmt.Errorf("invalid email format")
	}

	if utf8.RuneCountInString(admin.Password) < 6 {
		return fmt.Errorf("the password must be at least 6 characters long")
	}

	if strings.TrimSpace(admin.BirthDate) == "" {
		return fmt.Errorf("birth date is required")
	}

	parsedDate, err := time.Parse("2006-01-02", admin.BirthDate)
	if err != nil {
		return fmt.Errorf("invalid birth date format, expected YYYY-MM-DD")
	}

	if parsedDate.After(time.Now()) {
		return fmt.Errorf("birth date cannot be in the future")
	}

	age := calculateAge(parsedDate)
	if age < 18 {
		return fmt.Errorf("minimum age is 18 years old")
	}

	phone := strings.TrimSpace(admin.Phone)

	normalized := regexp.MustCompile(`\D`).ReplaceAllString(phone, "")

	if len(normalized) < 10 || len(normalized) > 11 {
		return fmt.Errorf("invalid phone number")
	}

	return nil
}

func calculateAge(birth time.Time) int {
	now := time.Now()

	age := now.Year() - birth.Year()

	if now.Month() < birth.Month() || (now.Month() == birth.Month() && now.Day() < birth.Day()) {
		age--
	}

	return age
}

func ValidatePatientInput(patient dtos.PatientInput) error {
	fullname := strings.TrimSpace(patient.FullName)

	if utf8.RuneCountInString(fullname) < 5 {
		return fmt.Errorf("name must be at least 5 characters long")
	}

	parts := strings.Fields(fullname)
	if len(parts) < 2 {
		return fmt.Errorf("full name must include first and last name")
	}

	_, err := mail.ParseAddress(patient.Email)
	if err != nil {
		return fmt.Errorf("invalid email format")
	}

	if utf8.RuneCountInString(patient.Password) < 6 {
		return fmt.Errorf("the password must be at least 6 characters long")
	}

	if strings.TrimSpace(patient.BirthDate) == "" {
		return fmt.Errorf("birth date is required")
	}

	parsedDate, err := time.Parse("2006-01-02", patient.BirthDate)
	if err != nil {
		return fmt.Errorf("invalid birth date format, expected YYYY-MM-DD")
	}

	if parsedDate.After(time.Now()) {
		return fmt.Errorf("birth date cannot be in the future")
	}

	age := calculateAge(parsedDate)
	if age < 18 {
		return fmt.Errorf("minimum age is 18 years old")
	}

	phone := strings.TrimSpace(patient.Phone)

	normalized := regexp.MustCompile(`\D`).ReplaceAllString(phone, "")

	if len(normalized) < 10 || len(normalized) > 11 {
		return fmt.Errorf("invalid phone number")
	}

	return nil
}

func ValidateCalendarSlotInput(input dtos.CalendarSlotsInput) error {
	if input.Weekday < 0 || input.Weekday > 6 {
		return fmt.Errorf("invalid weekday")
	}

	startTime, err := time.Parse("15:04", input.StartTime)
	if err != nil {
		return fmt.Errorf("invalid start time format, expected HH:MM")
	}

	endTime, err := time.Parse("15:04", input.EndTime)
	if err != nil {
		return fmt.Errorf("invalid end time format, expected HH:MM")
	}

	if !startTime.Before(endTime) {
		return fmt.Errorf("start time must be before end time")
	}

	return nil
}

func ValidateFullName(fullName string) error {
	fullName = strings.TrimSpace(fullName)

	if fullName == "" {
		return fmt.Errorf("full name is required")
	}

	if utf8.RuneCountInString(fullName) < 5 {
		return fmt.Errorf("name must be at least 5 characters long")
	}
	if utf8.RuneCountInString(fullName) > 60 {
		return fmt.Errorf("name must be at most 60 characters long")
	}

	parts := strings.Fields(fullName)

	if len(parts) < 2 {
		return fmt.Errorf("full name must include first and last name")
	}

	return nil
}

func ValidateEmail(email string) error {
	email = strings.TrimSpace(email)

	if email == "" {
		return fmt.Errorf("email is required")
	}

	if utf8.RuneCountInString(email) > 254 {
		return fmt.Errorf("email is too long")
	}

	parsed, err := mail.ParseAddress(email)
	if err != nil {
		return fmt.Errorf("invalid email format")
	}

	if parsed.Address != email {
		return fmt.Errorf("invalid email format")
	}

	return nil
}

func ValidateBirthDate(birthDate string) error {
	birthDate = strings.TrimSpace(birthDate)

	if birthDate == "" {
		return fmt.Errorf("birth date is required")
	}

	parsedDate, err := time.Parse("2006-01-02", birthDate)
	if err != nil {
		return fmt.Errorf("invalid birth date format, expected YYYY-MM-DD")
	}

	now := time.Now()

	if parsedDate.After(now) {
		return fmt.Errorf("birth date cannot be in the future")
	}

	if calculateAge(parsedDate) < 18 {
		return fmt.Errorf("minimum age is 18 years old")
	}

	return nil
}

var phoneRegex = regexp.MustCompile(`\D`)

func ValidatePhone(phone string) error {
	phone = strings.TrimSpace(phone)

	if phone == "" {
		return fmt.Errorf("phone is required")
	}

	normalized := phoneRegex.ReplaceAllString(phone, "")

	if len(normalized) < 9 || len(normalized) > 15 {
		return fmt.Errorf("invalid phone number")
	}

	return nil
}

func ValidateBio(bio string) error {
	bio = strings.TrimSpace(bio)

	if bio == "" {
		return nil
	}

	length := utf8.RuneCountInString(bio)

	if length < 10 {
		return fmt.Errorf("bio must be at least 10 characters long")
	}

	if length > 500 {
		return fmt.Errorf("bio must be at most 500 characters long")
	}

	return nil
}

func ValidateProfileImageURL(imageURL string) error {
	imageURL = strings.TrimSpace(imageURL)

	if imageURL == "" {
		return nil
	}

	if utf8.RuneCountInString(imageURL) > 2048 {
		return fmt.Errorf("profile image URL is too long")
	}

	parsedURL, err := url.Parse(imageURL)
	if err != nil {
		return fmt.Errorf("invalid profile image URL")
	}

	if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
		return fmt.Errorf("profile image URL must use http or https")
	}

	if parsedURL.Host == "" {
		return fmt.Errorf("invalid profile image URL")
	}

	return nil
}

func ValidateOfficeAddress(address string) error {
	address = strings.TrimSpace(address)

	if address == "" {
		return nil
	}

	length := utf8.RuneCountInString(address)

	if length < 5 {
		return fmt.Errorf("office address must be at least 5 characters long")
	}

	if length > 255 {
		return fmt.Errorf("office address must be at most 255 characters long")
	}

	return nil
}

var publicSlugRegex = regexp.MustCompile(`^[a-z0-9]+(?:-[a-z0-9]+)*$`)

func ValidatePublicSlug(slug string) error {
	slug = strings.TrimSpace(slug)

	if slug == "" {
		return nil
	}

	length := utf8.RuneCountInString(slug)

	if length < 3 {
		return fmt.Errorf("public slug must be at least 3 characters long")
	}

	if length > 50 {
		return fmt.Errorf("public slug must be at most 50 characters long")
	}

	if !publicSlugRegex.MatchString(slug) {
		return fmt.Errorf("public slug can only contain lowercase letters, numbers and hyphens")
	}

	return nil
}

func ValidateCRP(crp string) error {
	crp = strings.TrimSpace(crp)

	if crp == "" {
		return nil
	}

	if utf8.RuneCountInString(crp) > 50 {
		return fmt.Errorf("crp is too long")
	}

	if strings.ContainsAny(crp, "\r\n\t") {
		return fmt.Errorf("invalid crp")
	}

	return nil
}

func ValidateAdminProfileInput(input dtos.AdminProfileInput) error {
	if input.FullName != nil {
		if err := ValidateFullName(*input.FullName); err != nil {
			return err
		}
	}

	if input.Email != nil {
		if err := ValidateEmail(*input.Email); err != nil {
			return err
		}
	}

	if input.BirthDate != nil {
		if err := ValidateBirthDate(*input.BirthDate); err != nil {
			return err
		}
	}

	if input.CRP != nil {
		if err := ValidateCRP(*input.CRP); err != nil {
			return err
		}
	}

	if input.Bio != nil {
		if err := ValidateBio(*input.Bio); err != nil {
			return err
		}
	}

	if input.ProfileImageURL != nil {
		if err := ValidateProfileImageURL(*input.ProfileImageURL); err != nil {
			return err
		}
	}

	if input.OfficeAddress != nil {
		if err := ValidateOfficeAddress(*input.OfficeAddress); err != nil {
			return err
		}
	}

	if input.Phone != nil {
		if err := ValidatePhone(*input.Phone); err != nil {
			return err
		}
	}

	if input.PublicSlug != nil {
		if err := ValidatePublicSlug(*input.PublicSlug); err != nil {
			return err
		}
	}

	return nil
}
