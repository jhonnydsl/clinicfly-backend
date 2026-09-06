package services

import (
	"context"

	"github.com/google/uuid"
	auditservices "github.com/jhonnydsl/clinify-backend/src/audit/services"
	"github.com/jhonnydsl/clinify-backend/src/dtos"
	"github.com/jhonnydsl/clinify-backend/src/repository"
	"github.com/jhonnydsl/clinify-backend/src/utils"
)

type LoginService struct {
	Repo *repository.LoginRepository
	AuditService *auditservices.AuditService
}

func (service *LoginService) AuditLogin(ctx context.Context, actorID *uuid.UUID, actorRole *string, action, ipAddress, userAgent string) {
    err := service.AuditService.CreateAuditLog(ctx, dtos.AuditLogInput{
        ActorID:      actorID,
        ActorRole:    actorRole,
        Action:       action,
        ResourceType: "auth",
        ResourceID:   nil,
        IPAddress:    ipAddress,
        UserAgent:    userAgent,
    })

    if err != nil {
        utils.LogError("LoginUser service (audit log error)", err)
    }
}

func (service *LoginService) LoginUser(ctx context.Context, email, password, ipAddress, userAgent string) (dtos.LoginOutput, error) {
	admin, err := service.Repo.GetAdminByEmail(ctx, email)
	if err == nil {
		if err := utils.CheckPassword(admin.PasswordHash, password); err != nil {
			role := "admin"
			service.AuditLogin(ctx, &admin.ID, &role, "login.failed", ipAddress, userAgent)
			
			return dtos.LoginOutput{}, utils.BadRequestError("email or password incorrect")
		}

		token, err := utils.GenerateJWT(admin.ID.String(), admin.FullName, admin.Email, "admin")
		if err != nil {
			utils.LogError("LoginUser service (error generating token)", err)
			return dtos.LoginOutput{}, utils.InternalServerError("failed authentication")
		}

		role := "admin"
		service.AuditLogin(ctx, &admin.ID, &role, "login.success", ipAddress, userAgent)

		return dtos.LoginOutput{
			ID: admin.ID.String(),
			FullName: admin.FullName,
			Email: admin.Email,
			Role: "admin",
			Token: token,
		}, nil
	}

	patient, err := service.Repo.GetPatientByEmail(ctx, email)
	if err == nil {
		if err := utils.CheckPassword(patient.PasswordHash, password); err != nil {
			role := "patient"
			service.AuditLogin(ctx, &patient.ID, &role, "login.failed", ipAddress, userAgent)

			return dtos.LoginOutput{}, utils.BadRequestError("email or password incorrect")
		}

		token, err := utils.GenerateJWT(patient.ID.String(), patient.FullName, patient.Email, "patient")
		if err != nil {
			utils.LogError("loginUser service (error generating token)", err)
			return dtos.LoginOutput{}, utils.InternalServerError("failed authentication")
		}

		role := "patient"
		service.AuditLogin(ctx, &patient.ID, &role, "login.success", ipAddress, userAgent)

		return dtos.LoginOutput{
			ID: patient.ID.String(),
			FullName: patient.FullName,
			Email: patient.Email,
			Role: "patient",
			Token: token,
		}, nil
	}

	service.AuditLogin(ctx, nil, nil, "login.failed", ipAddress, userAgent)

	return dtos.LoginOutput{}, utils.BadRequestError("email or password incorrect")
}