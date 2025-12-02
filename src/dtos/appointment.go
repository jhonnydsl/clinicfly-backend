package dtos

import "github.com/google/uuid"

type AppointmentInput struct {
	PatientID string `json:"patient_id" binding:"required"`
	Date      string `json:"date" binding:"required"`
	StartTime string `json:"start_time" binding:"required"`
	EndTime   string `json:"end_time" binding:"required"`
}

type AppointmentOutput struct {
	ID 			uuid.UUID `json:"id"`
	PatientID 	uuid.UUID `json:"patient_id"`
	FullName    string    `json:"full_name"`
	Date 		string    `json:"date"`
	StartTime 	string    `json:"start_time"`
	EndTime 	string    `json:"end_time"`
	Status 		string    `json:"status"`
}