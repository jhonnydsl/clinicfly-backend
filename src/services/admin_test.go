package services_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jhonnydsl/clinify-backend/src/dtos"
	"github.com/jhonnydsl/clinify-backend/src/mocks"
	"github.com/jhonnydsl/clinify-backend/src/services"
	"golang.org/x/crypto/bcrypt"
)

func createTestService(mockRepo *mocks.MockAdminRepository) services.AdminService {
	mockMailer := &mocks.MockMailer{}
	
	return services.AdminService{
		Repo: mockRepo,
		Mailer: mockMailer,
	}
}

func createValidAppointmentInput() dtos.AppointmentInput {
	return dtos.AppointmentInput{
		PatientID: uuid.New().String(),
		Date:      "2026-08-28",
		StartTime: "13:00",
		EndTime:   "14:00",
	}
}

func createValidCalendarSlot() dtos.CalendarSlotDB {
	return dtos.CalendarSlotDB{
		StartTime: time.Date(2026, 8, 28, 8, 0, 0, 0, time.Local),
		EndTime:   time.Date(2026, 8, 28, 18, 0, 0, 0, time.Local),
	}
}

func createValidAdminInput() dtos.AdminInput {
	return dtos.AdminInput{
		FullName:  "Jhonny Lima",
		Email:     "jhonny@gmail.com",
		Password:  "123456",
		BirthDate: "1990-01-01",
		Phone:     "11999999999",
	}
}

func TestCreateAdminInvalidInput(t *testing.T) {
	mockRepo := &mocks.MockAdminRepository{}

	service := createTestService(mockRepo)

	admin := dtos.AdminInput{}

	_, err := service.CreateAdmin(context.Background(), admin)
	if err == nil {
		t.Error("expected error, got nil")
	}
}

func TestCreateAdminRepositoryError(t *testing.T) {
	mockRepo := &mocks.MockAdminRepository{
		CreateAdminError: errors.New("database error"),
	}

	service := createTestService(mockRepo)

	admin := createValidAdminInput()

	_, err := service.CreateAdmin(context.Background(), admin)
	if err == nil {
		t.Error("expected error, got nil")
	}
}

func TestCreateAdminSuccess(t *testing.T) {
	expectedID := uuid.New()

	mockRepo := &mocks.MockAdminRepository{
		CreateAdminID: expectedID,
	}

	service := createTestService(mockRepo)

	admin := createValidAdminInput()

	id, err := service.CreateAdmin(context.Background(), admin)

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	
	if err := bcrypt.CompareHashAndPassword(
		[]byte(mockRepo.ReceivedAdmin.Password),
		[]byte(admin.Password),
	); err != nil {
		t.Errorf("expected password to be a valid hash of the original password: %v", err)
	}


	if id != expectedID {
		t.Errorf("expected id %v, got %v", expectedID, id)
	}
}

func TestCreateAppointmentInvalidTimeRange(t *testing.T) {
	mockRepo := &mocks.MockAdminRepository{}
	service := createTestService(mockRepo)

	input := createValidAppointmentInput()
	input.StartTime = "14:00"
	input.EndTime = "13:00"

	_, err := service.CreateAppointment(context.Background(), input, uuid.New())
	if err == nil {
		t.Error("expected error, got nil")
	}
}

func TestCreateAppointmentGetCalendarSlotsError(t *testing.T) {
	mockRepo := &mocks.MockAdminRepository{
		GetCalendarSlotsByWeekdayError: errors.New("database error"),
	}

	service := createTestService(mockRepo)

	input := createValidAppointmentInput()

	_, err := service.CreateAppointment(context.Background(), input, uuid.New())
	if err == nil {
		t.Error("expected error, got nil")
	}
}

func TestCreateAppointmentOutsideAvailability(t *testing.T) {
	mockRepo := &mocks.MockAdminRepository{
		GetCalendarSlotsByWeekdaySlots: []dtos.CalendarSlotDB{
			createValidCalendarSlot(),
		},
	}
	service := createTestService(mockRepo)

	input := createValidAppointmentInput()
	input.StartTime = "19:00"
	input.EndTime = "20:00"

	_, err := service.CreateAppointment(context.Background(), input, uuid.New())

	if err == nil {
		t.Error("expected error, got nil")
	}
}

func TestCreateAppointmentsAlreadyBooked(t *testing.T) {
	mockRepo := &mocks.MockAdminRepository{
		GetCalendarSlotsByWeekdaySlots: []dtos.CalendarSlotDB{
				createValidCalendarSlot(),
		},
		GetAppointmentsByDateAppointments: []dtos.AppointmentOutput{
			{
				ID: 			uuid.New(),
				PatientID: 		uuid.New(),
				FullName: 		"Outro paciente",
				Date: 			"2026-08-28",
				StartTime: 		"13:00",
				EndTime: 		"14:00",
				Status: 		"scheduled",
			},
		},
	}

	service := createTestService(mockRepo)

	input := createValidAppointmentInput()
	input.StartTime = "13:00"
	input.EndTime = "14:00"

	_, err := service.CreateAppointment(context.Background(), input, uuid.New())
	if err == nil {
		t.Error("expected error, got nil")
	}
}

func TestCreateAppointmentsGetAppointmentsByDateError(t *testing.T) {
	mockRepo := &mocks.MockAdminRepository{
		GetCalendarSlotsByWeekdaySlots: []dtos.CalendarSlotDB{
			createValidCalendarSlot(),
		},
		GetCalendarSlotsByWeekdayError: errors.New("database error"),
	}

	service := createTestService(mockRepo)

	input := createValidAppointmentInput()

	_, err := service.CreateAppointment(context.Background(), input, uuid.New())
	if err == nil {
		t.Error("expected error, got nil")
	}
}

func TestCreateAppointmentRepositoryError(t *testing.T) {
	mockRepo := &mocks.MockAdminRepository{
		GetCalendarSlotsByWeekdaySlots: []dtos.CalendarSlotDB{
			createValidCalendarSlot(),
		},
		CreateAppointmentError: errors.New("database error"),
	}

	service := createTestService(mockRepo)

	input := createValidAppointmentInput()

	_, err := service.CreateAppointment(context.Background(), input, uuid.New())
	if err == nil {
		t.Error("expected error, got nil")
	}
}

func TestCreateAppointmentInvalidPatientID(t *testing.T) {
	mockRepo := &mocks.MockAdminRepository{
		GetCalendarSlotsByWeekdaySlots: []dtos.CalendarSlotDB{
			createValidCalendarSlot(),
		},
	}

	service := createTestService(mockRepo)

	input := createValidAppointmentInput()
	input.PatientID = "paciente-invalido"

	_, err := service.CreateAppointment(context.Background(), input, uuid.New())
	if err == nil {
		t.Error("expected error, got nil")
	}
}

func TestCreateAppointmentGetPatientEmailError(t *testing.T) {
	mockRepo := &mocks.MockAdminRepository{
		GetCalendarSlotsByWeekdaySlots: []dtos.CalendarSlotDB{
			createValidCalendarSlot(),
		},
		CreateAppointmentID: uuid.New(),
		GetPatientEmailByIDError: errors.New("database error"),
	}

	service := createTestService(mockRepo)

	input := createValidAppointmentInput()

	_, err := service.CreateAppointment(context.Background(), input, uuid.New())
	if err == nil {
		t.Error("expected error, got nil")
	}
}

func TestCreateAppointmentSuccess(t *testing.T) {
	expectedID := uuid.New()

	mockRepo := &mocks.MockAdminRepository{
		GetCalendarSlotsByWeekdaySlots: []dtos.CalendarSlotDB{
			createValidCalendarSlot(),
		},
		CreateAppointmentID: expectedID,
		GetPatientEmailByIDEmail: "paciente@gmail.com",
	}

	service := createTestService(mockRepo)
	input := createValidAppointmentInput()

	id, err := service.CreateAppointment(context.Background(), input, uuid.New())
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if id != expectedID {
		t.Errorf("expected id %v, got %v", expectedID, id)
	}
}