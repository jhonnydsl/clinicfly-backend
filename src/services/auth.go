package services

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"time"

	"github.com/jhonnydsl/clinify-backend/src/dtos"
	"github.com/jhonnydsl/clinify-backend/src/mailer"
	"github.com/jhonnydsl/clinify-backend/src/repository"
	"github.com/jhonnydsl/clinify-backend/src/utils"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	Repo *repository.AuthRepository
	Mailer *mailer.Mailer
}

func (service *AuthService) ForgotPassword(ctx context.Context, email string) error {
	user, err := service.Repo.GetUserByEmail(ctx, email)
	if err != nil {
		if utils.GetStatusCode(err) == http.StatusNotFound {
			return nil
		}

		utils.LogError("forgotPassword service (error getting user)", err)
		return err
	}

	tokenBytes := make([]byte, 32)

	_, err = rand.Read(tokenBytes)
	if err != nil {
		utils.LogError("forgotPassword service (error generating token)", err)
		return utils.InternalServerError("internal server error")
	}

	token := hex.EncodeToString(tokenBytes)

	expiresAt := time.Now().Add(30 * time.Minute)

	err = service.Repo.CreatePasswordResetToken(ctx, user.ID, user.Type, token, expiresAt)
	if err != nil {
		utils.LogError("forgotPassword service (error creating token)", err)
		return err
	}

	body := utils.BuildPasswordResetEmailBody(token)

	go func() {
		if err := service.Mailer.Send(
			user.Email,
			"Recuperação de Senha",
			body,
		); err != nil {
			utils.LogError("forgotPassword service (error sending email)", err)
		}
	}()

	return nil
}

func (service *AuthService) ResetPassword(ctx context.Context, input dtos.ResetPasswordInput) error {
	resetToken, err := service.Repo.GetValidPasswordResetToken(ctx, input.Token)
	if err != nil {
		utils.LogError("resetPassword service (error getting reset token)", err)
		return err
	}

	passwordHash, err := bcrypt.GenerateFromPassword(
		[]byte(input.NewPassword),
		bcrypt.DefaultCost,
	)
	if err != nil {
		utils.LogError("resetPassword service (error generating password hash)", err)
		return utils.InternalServerError("internal server error")
	}

	err = service.Repo.ResetPassword(ctx, resetToken.UserID, resetToken.ID, resetToken.UserType, string(passwordHash))
	if err != nil {
		utils.LogError("resetPassword service (error resetting password)", err)
		return err
	}

	return nil
}