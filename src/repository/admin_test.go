package repository

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/jhonnydsl/clinify-backend/src/dtos"
)

func setupMockDB(t *testing.T) sqlmock.Sqlmock {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("erro ao criar sqlmock: %v", err)
	}

	DB = db

	t.Cleanup(func() {
		DB.Close()
	})

	return mock
}

func createAppointmentRows(appointmentID, patientID uuid.UUID, fullName, status string) *sqlmock.Rows {
	return sqlmock.NewRows([]string{
		"id",
		"patient_id",
		"full_name",
		"date",
		"start_time",
		"end_time",
		"status",
	}).AddRow(
		appointmentID,
		patientID,
		fullName,
		time.Date(2026, 8, 28, 0, 0, 0, 0, time.UTC),
		time.Date(2026, 8, 28, 8, 0, 0, 0, time.UTC),
		time.Date(2026, 8, 28, 9, 0, 0, 0, time.UTC),
		status,
	)
}

func createValidAdminProfileInput() dtos.AdminProfileInput {
	fullName := "Jhonny Lima"
	email := "jhonny@email.com"
	birthDate := "2003-02-03"
	crp := "123456"
	bio := "Psicólogo clínico"
	profileImageURL := "https://example.com/profile.jpg"
	officeAddress := "Rua Exemplo, 123"
	phone := "11999999999"
	publicSlug := "jhonny-lima"

	return dtos.AdminProfileInput{
		FullName:        &fullName,
		Email:           &email,
		BirthDate:       &birthDate,
		CRP:             &crp,
		Bio:             &bio,
		ProfileImageURL: &profileImageURL,
		OfficeAddress:   &officeAddress,
		Phone:           &phone,
		PublicSlug:      &publicSlug,
	}
}

func TestGetAllAppointmentsSuccess(t *testing.T) {
	mock := setupMockDB(t)

	adminID := uuid.New()
	patientID := uuid.New()
	appointmentID := uuid.New()

	CountRows := sqlmock.NewRows([]string{"count"}).AddRow(2)

	mock.ExpectQuery(`SELECT COUNT\(\*\) FROM appointments WHERE client_id = \$1`).WithArgs(adminID).WillReturnRows(CountRows)

	rows := createAppointmentRows(appointmentID, patientID, "Jhonny Lima", "scheduled")

	mock.ExpectQuery(`SELECT a\.id, a\.patient_id, p\.full_name, a\.date, a\.start_time, a\.end_time, a\.status`).WithArgs(adminID, 10, 0).WillReturnRows(rows)

	repo := &AdminRepository{}

	result, total, err := repo.GetAllAppointments(context.Background(), adminID, "", 1, 10)
	if err != nil {
		t.Fatalf("espected error, got: %v", err)
	}

	if total != 2 {
		t.Fatalf("expected total 2, got %d", total)
	}

	if len(result) != 1 {
		t.Fatalf("expected 1 appointment, got %d", len(result))
	}

	if result[0].FullName != "Jhonny Lima" {
		t.Fatalf("expected Jhonny Lima, got %s", result[0].FullName)
	}

	if result[0].StartTime != "08:00" {
		t.Fatalf("eexpected 08:00, got %s", result[0].StartTime)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations")
	}
}

func TestGetAllAppointmentsWithStatus(t *testing.T) {
	mock := setupMockDB(t)

	adminID := uuid.New()
	patientID := uuid.New()
	appointmentID := uuid.New()
	status := "scheduled"

	countRows := sqlmock.NewRows([]string{"count"}).AddRow(1)

	mock.ExpectQuery(`SELECT COUNT\(\*\) FROM appointments WHERE client_id = \$1 AND status = \$2`).WithArgs(adminID, status).WillReturnRows(countRows)

	rows := createAppointmentRows(appointmentID, patientID, "Jhonny Lima", "scheduled")

	mock.ExpectQuery(`SELECT a\.id, a\.patient_id, p\.full_name, a\.date, a\.start_time, a\.end_time, a\.status`).WithArgs(adminID, status, 10, 0).WillReturnRows(rows)

	repo := &AdminRepository{}

	result, total, err := repo.GetAllAppointments(context.Background(), adminID, status, 1, 10)

	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	if total != 1 {
		t.Fatalf("expected total 1, got %d", total)
	}

	if len(result) != 1 {
		t.Fatalf("expected 1 appointment, got %d", len(result))
	}

	if result[0].Status != status {
		t.Fatalf("expected status %s, got %s", status, result[0].Status)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestUpdateAdminProfileSuccess(t *testing.T) {
	mock := setupMockDB(t)

	adminID := uuid.New()
	input := createValidAdminProfileInput()

	mock.ExpectExec(`UPDATE clients`).WithArgs(
		input.FullName,
		input.Email,
		input.BirthDate,
		input.CRP,
		input.Bio,
		input.ProfileImageURL,
		input.OfficeAddress,
		input.Phone,
		input.PublicSlug,
		adminID,
	).WillReturnResult(sqlmock.NewResult(1, 1))

	repo := &AdminRepository{}

	err := repo.UpdateAdminProfile(context.Background(), adminID, input)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestUpdateAdminProfileExecError(t *testing.T) {
	mock := setupMockDB(t)

	adminID := uuid.New()
	input := createValidAdminProfileInput()

	mock.ExpectExec(`UPDATE clients`).WithArgs(
		input.FullName,
		input.Email,
		input.BirthDate,
		input.CRP,
		input.Bio,
		input.ProfileImageURL,
		input.OfficeAddress,
		input.Phone,
		input.PublicSlug,
		adminID,
	).WillReturnError(errors.New("database error"))

	repo := &AdminRepository{}

	err := repo.UpdateAdminProfile(context.Background(), adminID, input)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if err.Error() != "internal server error" {
		t.Fatalf("expected 'internal server error', got: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}