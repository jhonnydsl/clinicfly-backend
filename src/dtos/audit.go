package dtos

import "github.com/google/uuid"

type AuditLogInput struct {
	ActorID 		*uuid.UUID
	ActorRole 		*string
	Action 			string
	ResourceType 	string
	ResourceID 		*uuid.UUID
	IPAddress 		string
	UserAgent 		string
}