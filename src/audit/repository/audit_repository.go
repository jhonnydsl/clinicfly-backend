package auditrepository

import (
	"context"

	"github.com/jhonnydsl/clinify-backend/src/dtos"
	"github.com/jhonnydsl/clinify-backend/src/repository"
	"github.com/jhonnydsl/clinify-backend/src/utils"
)

type AuditRepository struct{}

func (r *AuditRepository) CreateAuditLog(ctx context.Context, input dtos.AuditLogInput) error {
	query := `
		INSERT INTO audit_logs (
			actor_id,
			actor_role,
			action,
			resource_type,
			resource_id,
			ip_address,
			user_agent
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	_, err := repository.DB.ExecContext(
		ctx,
		query,
		input.ActorID,
		input.ActorRole,
		input.Action,
		input.ResourceType,
		input.ResourceID,
		input.IPAddress,
		input.UserAgent,
	)
	if err != nil {
		utils.LogError("createAuditLog auditrepository (insert error)", err)
		return utils.InternalServerError("error creating audit log")
	}

	return nil
}