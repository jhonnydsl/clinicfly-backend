package utils

import (
	"fmt"
	"net/mail"
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

	if now.YearDay() < birth.YearDay() {
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