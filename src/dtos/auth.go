package dtos

import (
	"time"

	"github.com/google/uuid"
)

type PasswordResetUser struct {
	ID 		uuid.UUID
	Email 	string
	Type 	string
}

type PasswordResetToken struct {
	ID 			uuid.UUID
	UserID 		uuid.UUID
	UserType 	string
	Token 		string
	ExpiresAt 	time.Time
	Used 		bool
	UsedAt 		*time.Time
}

type ForgotPasswordInput struct {
	Email string `json:"email" binding:"required,email"`
}

type ResetPasswordInput struct {
	Token 		string `json:"token" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=6"`
}