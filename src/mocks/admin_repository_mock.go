package mocks

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jhonnydsl/clinify-backend/src/dtos"
	"github.com/jhonnydsl/clinify-backend/src/services"
)

type MockMailer struct {
	SendError error
}

type MockAdminRepository struct{
	CreateAdminID 		uuid.UUID
	CreateAdminError 	error

	ReceivedAdmin 		dtos.AdminInput

	GetCalendarSlotsByWeekdaySlots []dtos.CalendarSlotDB
	GetCalendarSlotsByWeekdayError error

	GetAppointmentsByDateAppointments []dtos.AppointmentOutput
	GetAppointmentsByDateError error

	CreateAppointmentID uuid.UUID
	CreateAppointmentError error

	GetPatientEmailByIDEmail string
	GetPatientEmailByIDError error

	GetAllAppointmentsAppointments []dtos.AppointmentOutput
	GetAllAppointmentsTotal int
	GetAllAppointmentsError error
	GetAllAppointmentsPage int
	GetAllAppointmentsLimit int

	GetPatientsPatients []dtos.PatientOutput
	GetPatientsTotal int
	GetPatientsError error

	DeletePatientError error

	CreateCalendarSlotID uuid.UUID
	CreateCalendarSlotError error

	GetCalendarSlotsSlots []dtos.CalendarSlotsOutput
	GetCalendarSlotsError error

	DeleteCalendarSlotError error

	UpdateCalendarSlotError error

	GetAllAppointmentByIDAppointment dtos.AppointmentDetails
	GetAllAppointmentByIDError error

	CancelAppointmentByAdminError error
}

var _ services.AdminRepository = (*MockAdminRepository)(nil)

func (m *MockMailer) Send(to, object, body string) error {
	return m.SendError
}

func (m *MockAdminRepository) CreateAdmin(ctx context.Context, admin dtos.AdminInput, birthDate time.Time) (uuid.UUID, error) {
	m.ReceivedAdmin = admin
	
	return m.CreateAdminID, m.CreateAdminError
}

func(m *MockAdminRepository) CancelAppointmentByAdmin(ctx context.Context, appointmentID, adminID uuid.UUID) error {
	return m.CancelAppointmentByAdminError
}


func (m *MockAdminRepository) CreateAppointment(ctx context.Context, input dtos.AppointmentInput, parsedDate, start, end time.Time, clientID uuid.UUID) (uuid.UUID, error) {
	return m.CreateAppointmentID, m.CreateAppointmentError
}

func (m *MockAdminRepository) CreateCalendarSlot(ctx context.Context, input dtos.CalendarSlotsInput, start, end time.Time, adminID uuid.UUID) (uuid.UUID, error) {
	return m.CreateCalendarSlotID, m.CreateCalendarSlotError
}

func (m *MockAdminRepository) DeleteCalendarSlot(ctx context.Context, slotID uuid.UUID) error {
	return m.DeleteCalendarSlotError
}

func (m *MockAdminRepository) DeletePatient(ctx context.Context, patientID uuid.UUID) error {
	return m.DeletePatientError
}

func (m *MockAdminRepository) EmailExists(
	ctx context.Context,
	email string,
	adminID uuid.UUID,
) (bool, error) {
	return false, nil
}

func (m *MockAdminRepository) FindAdminIDBySlug(
	ctx context.Context,
	slug string,
) (uuid.UUID, error) {
	return uuid.Nil, nil
}

func (m *MockAdminRepository) GetAdminProfile(
	ctx context.Context,
	adminID uuid.UUID,
) (dtos.AdminProfileOutput, error) {
	return dtos.AdminProfileOutput{}, nil
}

func (m *MockAdminRepository) GetAllAppointmentByID(ctx context.Context, appointementID, adminID uuid.UUID) (dtos.AppointmentDetails, error) {
	return m.GetAllAppointmentByIDAppointment, m.GetAllAppointmentByIDError
}

func (m *MockAdminRepository) GetAllAppointments(ctx context.Context, adminID uuid.UUID, status string, page, limit int) ([]dtos.AppointmentOutput, int, error) {
	m.GetAllAppointmentsPage = page
	m.GetAllAppointmentsLimit = limit

	return m.GetAllAppointmentsAppointments, m.GetAllAppointmentsTotal, m.GetAllAppointmentsError
}

func (m *MockAdminRepository) GetAppointmentsByDate(ctx context.Context, adminID uuid.UUID, date string) ([]dtos.AppointmentOutput, error) {
	return m.GetAppointmentsByDateAppointments, m.GetAppointmentsByDateError
}

func (m *MockAdminRepository) GetCalendarSlots(ctx context.Context, adminID uuid.UUID) ([]dtos.CalendarSlotsOutput, error) {
	return m.GetCalendarSlotsSlots, m.GetCalendarSlotsError
}

func (m *MockAdminRepository) GetCalendarSlotsByWeekday(ctx context.Context, adminID uuid.UUID, weekday int) ([]dtos.CalendarSlotDB, error) {
	return m.GetCalendarSlotsByWeekdaySlots, m.GetCalendarSlotsByWeekdayError
}

func (m *MockAdminRepository) GetPatientEmailByID(ctx context.Context, patientID uuid.UUID) (string, error) {
	return m.GetPatientEmailByIDEmail, m.GetPatientEmailByIDError
}

func (m *MockAdminRepository) GetPatients(ctx context.Context, adminID uuid.UUID, page, limit int) ([]dtos.PatientOutput, int, error) {
	return m.GetPatientsPatients, m.GetPatientsTotal, m.GetPatientsError
}

func (m *MockAdminRepository) PublicSlugExists(
	ctx context.Context,
	slug string,
	adminID uuid.UUID,
) (bool, error) {
	return false, nil
}

func (m *MockAdminRepository) UpdateAdminProfile(
	ctx context.Context,
	adminID uuid.UUID,
	input dtos.AdminProfileInput,
) error {
	return nil
}

func (m *MockAdminRepository) UpdateAppointment(
	ctx context.Context,
	appointmentID,
	adminID uuid.UUID,
	date,
	startTime,
	endTime time.Time,
) error {
	return nil
}

func (m *MockAdminRepository) UpdateCalendarSlot(ctx context.Context, slotID, adminID uuid.UUID, input dtos.CalendarSlotsInput) error {
	return m.UpdateCalendarSlotError
}