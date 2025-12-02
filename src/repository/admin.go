package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jhonnydsl/clinify-backend/src/dtos"
	"github.com/jhonnydsl/clinify-backend/src/utils"
)

type AdminRepository struct{}

func (r *AdminRepository) CreateAdmin(ctx context.Context, admin dtos.AdminInput, birthDate time.Time) (uuid.UUID, error) {
	query := `INSERT INTO clients (full_name, email, password_hash, birth_date, crp, bio, profile_image_url, office_address, phone, public_slug)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	RETURNING id`

	var id uuid.UUID

	err := DB.QueryRowContext(
		ctx,
		query, 
		admin.FullName, 
		admin.Email, 
		admin.Password, 
		birthDate,
		admin.Crp,
		admin.Bio,
		admin.ProfileImage,
		admin.OfficeAddress,
		admin.Phone,
		admin.PublicSlug,
	).Scan(&id)
	if err != nil {
		utils.LogError("createAdmin (INSERT clients)", err)
		return uuid.UUID{}, utils.InternalServerError("error creating user admin")
	}

	return id, nil
}

func (r *AdminRepository) FindAdminIDBySlug(ctx context.Context, slug string) (uuid.UUID, error) {
	query := `SELECT id FROM clients WHERE public_slug = $1 LIMIT 1`

	var id uuid.UUID

	err := DB.QueryRowContext(ctx, query, slug).Scan(&id)
	if err != nil {
		utils.LogError("FindAdminBySlug (error in SELECT clients)", err)
		return uuid.UUID{}, utils.InternalServerError("invalid client slug")
	}

	return id, nil
}

func (r *AdminRepository) CreateAppointment(ctx context.Context, input dtos.AppointmentInput, parsedDate, start, end time.Time, clientID uuid.UUID) (uuid.UUID, error) {
	query := `INSERT INTO appointments (client_id, patient_id, date, start_time, end_time, status)
	VALUES ($1, $2, $3, $4, $5, 'scheduled')
	RETURNING id;`

	var id uuid.UUID

	err := DB.QueryRowContext(
		ctx,
		query,
		clientID,
		input.PatientID,
		parsedDate,
		start,
		end,
	).Scan(&id)
	if err != nil {
		utils.LogError("createAppointment repository (error in INSERT)", err)
		return uuid.UUID{}, utils.InternalServerError("error creating appointment")
	}

	return id, nil
}

func (r *AdminRepository) GetPatients(ctx context.Context, adminID uuid.UUID, page, limit int) ([]dtos.PatientOutput, int, error) {
	query := `SELECT id, full_name, email, phone, birth_date FROM patients
	WHERE client_id = $1
	ORDER BY full_name LIMIT $2 OFFSET $3`

	queryCount := `SELECT COUNT(*) FROM patients WHERE client_id = $1 `

	offset := (page - 1) * limit

	var total int
	
	err := DB.QueryRowContext(ctx, queryCount, adminID).Scan(&total)
	if err != nil {
		return nil, 0, utils.InternalServerError("error getting total patients")
	}
	
	rows, err := DB.QueryContext(ctx, query, adminID, limit, offset)
	if err != nil {
		utils.LogError("GetPatients repository (error in SELECT)", err)
		return nil, 0, utils.InternalServerError("error getting patients")
	}
	defer rows.Close()
	
	var patients []dtos.PatientOutput

	for rows.Next() {
		var (
			id uuid.UUID
			fullName string
			email string
			phone string
			birthDate time.Time
		)

		err = rows.Scan(&id, &fullName, &email, &phone, &birthDate)
		if err != nil {
			utils.LogError("getPatients repository (scan error)", err)
			return nil, 0, utils.InternalServerError("error fetching patients")
		}

		patients = append(patients, dtos.PatientOutput{
			ID: id,
			FullName: fullName,
			Email: email,
			Phone: phone,
			BirthDate: birthDate.Format("2006-01-02"),
		})
	}
	
	return patients, total, nil
}

func (r *AdminRepository) DeletePatient(ctx context.Context, patientID uuid.UUID) error {
	query := `DELETE FROM patients WHERE id = $1`

	res, err := DB.ExecContext(ctx, query, patientID)
	if err != nil {
		utils.LogError("deletePatient repository (error deleting patient)", err)
		return utils.InternalServerError("error deleting patient")
	}

	rows, err := res.RowsAffected()
	if err != nil {
		utils.LogError("deletePatient repository (error reading rows affected)", err)
		return utils.InternalServerError("error deleting patient")
	}

	if rows == 0 {
		return utils.NotFoundError("patient not found")
	}

	return nil
}