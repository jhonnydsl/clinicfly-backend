package dtos

import "github.com/google/uuid"

type CalendarSlotsInput struct {
	Weekday   int    `json:"weekday" binding:"required"`
	StartTime string `json:"start_time" binding:"required"`
	EndTime   string `json:"end_time" binding:"required"`
}

type CalendarSlotsOutput struct {
	ID 			uuid.UUID `json:"id"`
	Weekday 	string 	  `json:"weekday"`
	StartTime 	string 	  `json:"start_time"`
	EndTime 	string 	  `json:"end_time"`
}