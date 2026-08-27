package services

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jhonnydsl/clinify-backend/src/dtos"
)

type AdminRepository interface {
	CreateAdmin(
		ctx context.Context, 
		admin dtos.AdminInput, 
		birthDate time.Time,
		) (uuid.UUID, error)

	FindAdminIDBySlug(
		ctx context.Context, 
		slug string,
		) (uuid.UUID, error)

	CreateAppointment(
		ctx context.Context, 
		input dtos.AppointmentInput,
		parsedDate,
		start,
		end time.Time,
		clientID uuid.UUID,
	) (uuid.UUID, error)

	GetAllAppointments(
		ctx context.Context,
		adminID uuid.UUID,
		status string,
		page,
		limit int,
	) ([]dtos.AppointmentOutput, int, error)

	GetAppointmentsByDate(
		ctx context.Context,
		adminID uuid.UUID,
		date string,
	) ([]dtos.AppointmentOutput, error)

	GetPatients(
		ctx context.Context,
		adminID uuid.UUID,
		page,
		limit int,
	) ([]dtos.PatientOutput, int, error)

	DeletePatient(
		ctx context.Context,
		patientID uuid.UUID,
	) error

	GetPatientEmailByID(
		ctx context.Context,
		patientID uuid.UUID,
	) (string, error)

	CreateCalendarSlot(
		ctx context.Context,
		input dtos.CalendarSlotsInput,
		start,
		end time.Time,
		adminID uuid.UUID,
	) (uuid.UUID, error)

	GetCalendarSlots(
		ctx context.Context,
		adminID uuid.UUID,
	) ([]dtos.CalendarSlotsOutput, error)

	DeleteCalendarSlot(
		ctx context.Context,
		slotID uuid.UUID,
	) error

	GetCalendarSlotsByWeekday(
		ctx context.Context,
		adminID uuid.UUID,
		weekday int,
	) ([]dtos.CalendarSlotDB, error)

	UpdateCalendarSlot(
		ctx context.Context,
		slotID,
		adminID uuid.UUID,
		input dtos.CalendarSlotsInput,
	) error

	CancelAppointmentByAdmin(
		ctx context.Context,
		appointmentID,
		adminID uuid.UUID,
	) error

	GetAllAppointmentByID(
		ctx context.Context,
		appointmentID,
		adminID uuid.UUID,
	) (dtos.AppointmentDetails, error)

	UpdateAppointment(
		ctx context.Context,
		appointmentID,
		adminID uuid.UUID,
		date,
		startTime,
		endTime time.Time,
	) error

	GetAdminProfile(
		ctx context.Context,
		adminID uuid.UUID,
	) (dtos.AdminProfileOutput, error)

	UpdateAdminProfile(
		ctx context.Context,
		adminID uuid.UUID,
		input dtos.AdminProfileInput,
	) error

	PublicSlugExists(
		ctx context.Context,
		slug string,
		adminID uuid.UUID,
	) (bool, error)

	EmailExists(
		ctx context.Context,
		email string,
		adminID uuid.UUID,
	) (bool, error)
}