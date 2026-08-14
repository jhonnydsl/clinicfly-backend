package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/jhonnydsl/clinify-backend/src/dtos"
	"github.com/jhonnydsl/clinify-backend/src/utils"
)

type AuthRepository struct{}

func (r *AuthRepository) GetUserByEmail(ctx context.Context, email string) (dtos.PasswordResetUser, error) {
	query := `
		SELECT id, email, 'admin' AS user_type
		FROM clients
		WHERE email = $1

		UNION ALL

		SELECT id, email, 'patient' AS user_type
		FROM patients
		WHERE email = $1
	`

	var user dtos.PasswordResetUser

	err := DB.QueryRowContext(ctx, query, email).Scan(
		&user.ID,
		&user.Email,
		&user.Type,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return dtos.PasswordResetUser{}, utils.NotFoundError("user not found")
		}

		utils.LogError("getUserByEmail repository (query error)", err)
		return dtos.PasswordResetUser{}, utils.InternalServerError("error getting user")
	}

	return user, nil
}

func (r *AuthRepository) CreatePasswordResetToken(ctx context.Context, userID uuid.UUID, userType, token string, expiresAt time.Time) error {
	query := `
		INSERT INTO password_reset_tokens (
			user_id,
			user_type,
			token,
			expires_at
		)
		VALUES ($1, $2, $3, $4)
	`

	_, err := DB.ExecContext(ctx, query, userID, userType, token, expiresAt)
	if err != nil {
		utils.LogError("createPasswordResetToken repository (insert error)", err)
		return utils.InternalServerError("error creating password reset token")
	}

	return nil
}

func (r *AuthRepository) GetValidPasswordResetToken(ctx context.Context, token string) (dtos.PasswordResetToken, error) {
	query := `
		SELECT
			id,
			user_id,
			user_type,
			token,
			expires_at,
			used,
			used_at
		FROM password_reset_tokens
		WHERE token = $1
			AND used = FALSE
			AND expires_at > CURRENT_TIMESTAMP
	`

	var resetToken dtos.PasswordResetToken

	err := DB.QueryRowContext(ctx, query, token).Scan(
		&resetToken.ID,
		&resetToken.UserID,
		&resetToken.UserType,
		&resetToken.Token,
		&resetToken.ExpiresAt,
		&resetToken.Used,
		&resetToken.UsedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return dtos.PasswordResetToken{}, utils.NotFoundError("invalid or expired password reset token")
		}

		utils.LogError("getValidPasswordResetToken repository (query error)", err)
		return dtos.PasswordResetToken{}, utils.InternalServerError("error getting password reset token")
	}

	return resetToken, nil
}

// Unico método que redefine o password e invalida o token de uma única vez
func (r *AuthRepository) ResetPassword(ctx context.Context, userID, tokenID uuid.UUID, userType, passwordHash string) error {
	tx, err := DB.BeginTx(ctx, nil)
	if err != nil {
		utils.LogError("resetPassword repository (begin transaction error)", err)
		return utils.InternalServerError("internal server error")
	}

	defer tx.Rollback()

	var query string

	switch userType {
	case "admin":
		query = `
			UPDATE clients
			SET password_hash = $1
			WHERE id = $2
		`

	case "patient":
		query = `
			UPDATE patients
			SET password_hash = $1
			WHERE id = $2
		`

	default:
		return utils.BadRequestError("invalid user type")
	}

	result, err := tx.ExecContext(ctx, query, passwordHash, userID)
	if err != nil {
		utils.LogError("resetPassword repository (password update error)", err)
		return utils.InternalServerError("internal server error")
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		utils.LogError("resetPassword repository (password rows affected error)", err)
		return utils.InternalServerError("internal server error")
	}

	if rowsAffected == 0 {
		return utils.NotFoundError("user not found")
	}

	result, err = tx.ExecContext(ctx, `
		UPDATE password_reset_tokens
		SET
			used = TRUE,
			used_at = CURRENT_TIMESTAMP
		WHERE id = $1
			AND used = FALSE
	`, tokenID)

	if err != nil {
		utils.LogError("resetPassword repository (token update error)", err)
		return utils.InternalServerError("internal server error")
	}

	rowsAffected, err = result.RowsAffected()
	if err != nil {
		utils.LogError("resetPassword repository (token rows affected error)", err)
		return utils.InternalServerError("internal server error")
	}

	if rowsAffected == 0 {
		return utils.NotFoundError("password reset token not found")
	}

	if err := tx.Commit(); err != nil {
		utils.LogError("resetPassword repository (commit error)", err)
		return utils.InternalServerError("internal server error")
	}

	return nil
}