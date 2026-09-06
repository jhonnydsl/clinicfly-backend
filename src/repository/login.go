package repository

import (
	"context"
	"database/sql"

	"github.com/jhonnydsl/clinify-backend/src/dtos"
	"github.com/jhonnydsl/clinify-backend/src/utils"
)

type LoginRepository struct{}

func (r *LoginRepository) GetAdminByEmail(ctx context.Context, email string) (dtos.LoginAdmin, error) {
	query := `SELECT id, full_name, email, password_hash FROM clients WHERE email = $1`

	var admin dtos.LoginAdmin

	err := DB.QueryRowContext(ctx, query, email).Scan(
		&admin.ID,
		&admin.FullName,
		&admin.Email,
		&admin.PasswordHash,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return dtos.LoginAdmin{}, err
		}

		utils.LogError("getAdminByEmail repository (error select data in db)", err)
		return dtos.LoginAdmin{}, utils.InternalServerError("error logging in")
	}

	return admin, nil
}

func (r *LoginRepository) GetPatientByEmail(ctx context.Context, email string) (dtos.LoginPatient, error) {
	query := `SELECT id, full_name, email, password_hash FROM patients WHERE email = $1`

	var patient dtos.LoginPatient

	err := DB.QueryRowContext(ctx, query, email).Scan(
		&patient.ID,
		&patient.FullName,
		&patient.Email,
		&patient.PasswordHash,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return dtos.LoginPatient{}, err
		}
		utils.LogError("getPatientByEmail repository (error select data in db)", err)
		return dtos.LoginPatient{}, utils.InternalServerError("error logging in")
	}

	return patient, nil
}