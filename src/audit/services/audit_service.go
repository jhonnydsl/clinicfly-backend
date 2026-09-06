package auditservices

import (
	"context"

	auditrepository "github.com/jhonnydsl/clinify-backend/src/audit/repository"
	"github.com/jhonnydsl/clinify-backend/src/dtos"
	"github.com/jhonnydsl/clinify-backend/src/utils"
)

type AuditService struct {
	Repo *auditrepository.AuditRepository
}

func (service *AuditService) CreateAuditLog(ctx context.Context, input dtos.AuditLogInput) error {
	err := service.Repo.CreateAuditLog(ctx, input)
	if err != nil {
		utils.LogError("createAuditLog auditservice (repository error)", err)
		return err
	}

	return nil
}