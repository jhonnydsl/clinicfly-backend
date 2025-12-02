package dtos

import "github.com/google/uuid"

type PatientInput struct {
	FullName   string `json:"full_name" binding:"required"`
	Email      string `json:"email" binding:"required,email"`
	Password   string `json:"password_hash" binding:"required"`
	Phone      string `json:"phone" binding:"required"`
	BirthDate  string `json:"birth_date" binding:"required"`
	PublicSlug string `json:"public_slug" binding:"required"`
}

type LoginPatient struct {
	ID 				uuid.UUID 	`json:"id"`
	FullName 		string `json:"full_name"`
	Email 			string 		`json:"email"`
	PasswordHash 	string 		`json:"password_hash"`
}

type PatientOutput struct {
	ID 			uuid.UUID `json:"id"`
	FullName 	string 	  `json:"full_name"`
	Email 		string    `json:"email"`
	Phone 		string    `json:"phone"`
	BirthDate 	string    `json:"birth_date"`
}