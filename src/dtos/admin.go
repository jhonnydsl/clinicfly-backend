package dtos

import "github.com/google/uuid"

type AdminInput struct {
	FullName      string `json:"full_name" binding:"required"`
	Email         string `json:"email" binding:"required,email"`
	Password      string `json:"password" binding:"required"`
	BirthDate     string `json:"birth_date"`
	Crp           string `json:"crp"`
	Bio           string `json:"bio"`
	ProfileImage  string `json:"profile_image_url"`
	OfficeAddress string `json:"office_address"`
	Phone         string `json:"phone"`
	PublicSlug    string `json:"public_slug"`
}

type LoginAdmin struct {
	ID            uuid.UUID `json:"id"`
	FullName 	  string 	`json:"full_name"`
	Email         string    `json:"email"`
	PasswordHash  string    `json:"password_hash"`
}

type AdminProfileInput struct {
	FullName 			*string 	`json:"full_name"`
	Email 				*string 	`json:"email"`
	BirthDate 			*string 	`json:" birth_date"`
	CRP 				*string 	`json:"crp"`
	Bio 				*string 	`json:"bio"`
	ProfileImageURL 	*string 	`json:"profile_image_url"`
	OfficeAddress 		*string 	`json:"office_address"`
	Phone 				*string 	`json:"phone"`
	PublicSlug 			*string 	`json:"public_slug"`
}

type AdminProfileOutput struct {
	ID 					uuid.UUID 	`json:"id"`
	FullName 			string 		`json:"full_name"`
	Email 				string 		`json:"email"`
	BirthDate 			string 		`json:"birth_date"`
	CRP 				string 		`json:"crp"`
	Bio 				string 		`json:"bio"`
	ProfileImageURL 	string 		`json:"profile_image_url"`
	OfficeAddress 		string 		`json:"office_address"`
	Phone 				string 		`json:"phone"`
	PublicSlug 			string 		`json:"public_slug"`
}