package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/jhonnydsl/clinify-backend/src/dtos"
	"github.com/jhonnydsl/clinify-backend/src/utils"
)

type PatientRepository struct{}

func (r *PatientRepository) CreatePatient(ctx context.Context, patient dtos.PatientInput, birthDate time.Time, clientID uuid.UUID) (uuid.UUID, error) {
	query := `INSERT INTO patients (full_name, email, password_hash, phone, birth_date, client_id)
	VALUES ($1, $2, $3, $4, $5, $6)
	RETURNING id`

	var id uuid.UUID

	err := DB.QueryRowContext(
		ctx,
		query,
		patient.FullName,
		patient.Email,
		patient.Password,
		patient.Phone,
		birthDate,
		clientID,
	).Scan(&id)
	if err != nil {
		utils.LogError("PatientRepository (erro ao criar user paciente)", err)
		return uuid.UUID{}, utils.InternalServerError("error creating user patient")
	}

	return id, nil
}

func (r *PatientRepository) CancelAppointmentByPatient(ctx context.Context, appointmentID, patientID uuid.UUID) error {
	query := `
		UPDATE appointments
		SET status = 'cancelled'
		WHERE id = $1
			AND patient_id = $2
			AND status = 'scheduled'
	`

	result, err := DB.ExecContext(ctx, query, appointmentID, patientID)
	if err != nil {
		utils.LogError("cancelAppointmentByPatient repository (update error)", err)
		return utils.InternalServerError("error cancelling appointment")
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		utils.LogError("cancelAppointmentByPatient repository (rows affected error)", err)
		return utils.InternalServerError("error cancelling appointment")
	}

	if rowsAffected == 0 {
		return utils.NotFoundError("appointment not found")
	}

	return nil
}

func (r *PatientRepository) GetAppointmentByID(ctx context.Context, appointmentID, patientID uuid.UUID) (dtos.AppointmentDetails, error) {
	query := `
		SELECT
			a.id,
			a.date,
			a.start_time,
			a.end_time,
			p.email,
			c.email
		FROM appointments a
		JOIN patients p ON p.id = a.patient_id
		JOIN clients c ON c.id = a.client_id
		WHERE a.id = $1
			AND a.patient_id = $2
			AND a.status = 'scheduled'
	`

	var (
		appointment dtos.AppointmentDetails
		dateDB time.Time
		startTime time.Time
		endTime time.Time
	)

	err := DB.QueryRowContext(ctx, query, appointmentID, patientID).Scan(
		&appointment.ID,
		&dateDB,
		&startTime,
		&endTime,
		&appointment.PatientEmail,
		&appointment.AdminEmail,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return dtos.AppointmentDetails{}, utils.NotFoundError("appointment not found")
		}

		utils.LogError("getAppointmentByID repository (query error)", err)
		return dtos.AppointmentDetails{}, utils.InternalServerError("error getting appointment")
	}

	appointment.Date = dateDB.Format("2006-01-02")
	appointment.StartTime = startTime.Format("15:04")
	appointment.EndTime = endTime.Format("15:04")

	return appointment, nil
}