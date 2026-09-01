package services_test

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jhonnydsl/clinify-backend/src/dtos"
	"github.com/jhonnydsl/clinify-backend/src/mocks"
	"github.com/jhonnydsl/clinify-backend/src/services"
	"github.com/jhonnydsl/clinify-backend/src/utils"
	"github.com/patrickmn/go-cache"
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

func createValidAppointmentOutput() dtos.AppointmentOutput {
	return dtos.AppointmentOutput{
		ID:        uuid.New(),
		PatientID: uuid.New(),
		FullName:  "Paciente Teste",
		Date:      "2026-08-28",
		StartTime: "13:00",
		EndTime:   "14:00",
		Status:    "scheduled",
	}
}

func createValidPatientOutput() dtos.PatientOutput {
	return dtos.PatientOutput{
		ID:        uuid.New(),
		FullName:  "Paciente Teste",
		Email:     "paciente@gmail.com",
		Phone:     "11999999999",
		BirthDate: "1990-01-01",
	}
}

func createValidInputSlot() dtos.CalendarSlotsInput {
	return dtos.CalendarSlotsInput{
		Weekday: 5,
		StartTime: "13:00",
		EndTime: "14:00",
	}
}

func createValidCalendarSlotsOutput() []dtos.CalendarSlotsOutput {
	return []dtos.CalendarSlotsOutput{
		{
			ID:        uuid.New(),
			Weekday:   "fryday",
			StartTime: "08:00",
			EndTime:   "18:00",
		},
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

func TestGetAppointmentsInvalidStatus(t *testing.T) {
	mockRepo := &mocks.MockAdminRepository{}
	service := createTestService(mockRepo)

	_, _, err := service.GetAppointments(context.Background(), uuid.New(), "invalid", 1, 10)
	if err == nil {
		t.Error("expected error, got nill")
	}
}

func TestGetAppointmentsRepositoryError(t *testing.T) {
	mockRepo := &mocks.MockAdminRepository{
		GetAllAppointmentsError: errors.New("database error"),
	}

	service := createTestService(mockRepo)

	_, _, err := service.GetAppointments(context.Background(), uuid.New(), "scheduled", 1, 10)
	if err == nil {
		t.Error("expected error, got nil")
	}
}

func TestGetAllAppointmentsSuccess(t *testing.T) {
	expectedAppointemnts := []dtos.AppointmentOutput{
		createValidAppointmentOutput(),
	}

	expectedTotal := 1

	mockRepo := &mocks.MockAdminRepository{
		GetAllAppointmentsAppointments: expectedAppointemnts,
		GetAllAppointmentsTotal: expectedTotal,
	}

	service := createTestService(mockRepo)

	appointments, total, err := service.GetAppointments(context.Background(), uuid.New(), "scheduled", 1, 10)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if total != expectedTotal {
		t.Errorf("expected %d appointments, got %d", len(expectedAppointemnts), len(appointments))
	}
}

func TestGetAllAppointmentsInvalidPagination(t *testing.T) {
	mockRepo := &mocks.MockAdminRepository{}
	service := createTestService(mockRepo)

	_, _, err := service.GetAppointments(context.Background(), uuid.New(), "scheduled", 0, 0)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if mockRepo.GetAllAppointmentsPage != 1 {
		t.Errorf("expected page 1, got %d", mockRepo.GetAllAppointmentsPage)
	}

	if mockRepo.GetAllAppointmentsLimit != 10 {
		t.Errorf("expected limit 10, got %d", mockRepo.GetAllAppointmentsLimit)
	}
}

func TestGetAppointmentsFromCache(t *testing.T) {
	adminID := uuid.New()
	status := "scheduled"
	page := 1
	limit := 10

	expectedAppointemnts := []dtos.AppointmentOutput{
		createValidAppointmentOutput(),
	}

	expectedTotal := 1

	chacheKey := fmt.Sprintf("appointments_admin_%s_status_%s_page_%d_limit_%d", adminID.String(), status, page, limit)

	utils.Cache.Set(
		chacheKey,
		&utils.AppointmentsCache{
			Data: expectedAppointemnts,
			Total: expectedTotal,
		},
		cache.DefaultExpiration,
	)

	mockRepo := &mocks.MockAdminRepository{
		GetAllAppointmentsError: errors.New("repository should not be called"),
	}

	service := createTestService(mockRepo)

	appointments, total, err := service.GetAppointments(context.Background(), adminID, status, page, limit)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if total != expectedTotal {
		t.Errorf("expected total %d, got %d", expectedTotal, total)
	}

	if len(appointments) != len(expectedAppointemnts) {
		t.Errorf("expected %d appointments, got %d", len(expectedAppointemnts), len(appointments))
	}
}

func TestGetPatientsSuccess(t *testing.T) {
	expectedPatients := []dtos.PatientOutput{
		createValidPatientOutput(),
	}

	expectedTotal := 1

	mockRepo := &mocks.MockAdminRepository{
		GetPatientsPatients: expectedPatients,
		GetPatientsTotal: expectedTotal,
	}

	service := createTestService(mockRepo)

	patients, total, err := service.GetPatients(context.Background(), uuid.New(), 1, 10)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if total != expectedTotal {
		t.Errorf("expected total %d, got %d", expectedTotal, total)
	}

	if len(patients) != len(expectedPatients) {
		t.Errorf("expected %d patients, got %d", len(expectedPatients), len(patients))
	}
}

func TestGetPatientsRepositoryError(t *testing.T) {
	utils.Cache.Flush()

	mockRepo := &mocks.MockAdminRepository{
		GetPatientsError: errors.New("database error"),
	}

	service := createTestService(mockRepo)

	_, _, err := service.GetPatients(context.Background(), uuid.New(), 1, 10)
	if err == nil {
		t.Error("expected error, got nil")
	}
}

func TestDeletePatientInvalidID(t *testing.T) {
	mockRepo := &mocks.MockAdminRepository{}
	service := createTestService(mockRepo)

	err := service.DeletePatient(context.Background(), uuid.Nil)
	if err == nil {
		t.Error("expected error, got nil")
	}
}

func TestDeletePatientRepositoryError(t *testing.T) {
	mockRepo := &mocks.MockAdminRepository{
		DeletePatientError: errors.New("database error"),
	}

	service := createTestService(mockRepo)

	err := service.DeletePatient(context.Background(), uuid.New())
	if err == nil {
		t.Error("expected error, got nil")
	}
}

func TestCreateCalendarSlotOverlapping(t *testing.T) {
	mockRepo := &mocks.MockAdminRepository{
		GetCalendarSlotsByWeekdaySlots: []dtos.CalendarSlotDB{
			createValidCalendarSlot(),
		},
	}

	service := createTestService(mockRepo)

	input := createValidInputSlot()

	_, err := service.CreateCalendarSlot(context.Background(), input, uuid.New())
	if err == nil {
		t.Error("expected error, got nil")
	}
}

func TestCreateCalendarSlotSuccess(t *testing.T) {
	expectedID := uuid.New()

	mockRepo := &mocks.MockAdminRepository{
		CreateCalendarSlotID: expectedID,
	}

	service := createTestService(mockRepo)

	input := createValidInputSlot()

	id, err := service.CreateCalendarSlot(context.Background(), input, uuid.New())
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if id != expectedID {
		t.Errorf("expected id %v, got %v", expectedID, id)
	}
}

func TestCreateCalendarSlotGetSlotsError(t *testing.T) {
	mockRepo := &mocks.MockAdminRepository{
		GetCalendarSlotsByWeekdayError: errors.New("database error"),
	}

	service := createTestService(mockRepo)

	input := createValidInputSlot()

	_, err := service.CreateCalendarSlot(context.Background(), input, uuid.New())
	if err == nil {
		t.Error("expected error, got nil")
	}
}

func TestCreateCalendarSlotCreateError(t *testing.T) {
	mockRepo := &mocks.MockAdminRepository{
		CreateCalendarSlotID: uuid.New(),
		CreateCalendarSlotError: errors.New("database error"),
	}

	service := createTestService(mockRepo)

	input := createValidInputSlot()

	_, err := service.CreateCalendarSlot(context.Background(), input, uuid.New())
	if err == nil {
		t.Error("expected error, got nil")
	}
}

func TestGetCalendarSlotsSuccess(t *testing.T) {
	expectedSlots := createValidCalendarSlotsOutput()

	mockRepo := &mocks.MockAdminRepository{
		GetCalendarSlotsSlots: expectedSlots,
	}

	service := createTestService(mockRepo)

	slots, err := service.GetCalendarSlots(context.Background(), uuid.New())
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if len(slots) != len(expectedSlots) {
		t.Errorf("expected %d slots, got %d", len(expectedSlots), len(slots))
	}
}

func TestGetCalendarSlotsRepositoryError(t *testing.T) {
	mockRepo := &mocks.MockAdminRepository{
		GetCalendarSlotsError: errors.New("database error"),
	}

	service := createTestService(mockRepo)

	_, err := service.GetCalendarSlots(context.Background(), uuid.New())
	if err == nil {
		t.Error("expected error, got nil")
	}
}

func TestDeleteCalendarSlotInvalidID(t *testing.T) {
	mockRepo := &mocks.MockAdminRepository{}
	service := createTestService(mockRepo)

	err := service.DeleteCalendarSlot(context.Background(), uuid.Nil)
	if err == nil {
		t.Error("expected error, got nil")
	}
}

func TestDeleteCalendarSlotRepositoryError(t *testing.T) {
	mockRepo := mocks.MockAdminRepository{
		DeleteCalendarSlotError: errors.New("database error"),
	}

	service := createTestService(&mockRepo)

	err := service.DeleteCalendarSlot(context.Background(), uuid.New())
	if err == nil {
		t.Error("expected error, got nil")
	}
}

func TestUpdateCalendarSlotOverlapping(t *testing.T) {
	slotID := uuid.New()

	existingSlot := createValidCalendarSlot()
	existingSlot.ID = uuid.New()

	slotToUpdate := createValidCalendarSlot()
	slotToUpdate.ID = slotID

	mockRepo := &mocks.MockAdminRepository{
		GetCalendarSlotsByWeekdaySlots: []dtos.CalendarSlotDB{
			existingSlot,
			slotToUpdate,
		},
	}

	service := createTestService(mockRepo)

	input := createValidInputSlot()
	input.StartTime = "17:00"
	input.EndTime = "19:00"

	_, err := service.UpdateCalendarSlot(context.Background(), slotID, uuid.New(), input)
	if err == nil {
		t.Error("expected error, got nil")
	}
}

func TestUpdateCalendarSlotGetSlotsError(t *testing.T) {
	mockRepo := &mocks.MockAdminRepository{
		GetCalendarSlotsByWeekdayError: errors.New("database error"),
	}

	service := createTestService(mockRepo)
	input := createValidInputSlot()

	_, err := service.UpdateCalendarSlot(context.Background(), uuid.New(), uuid.New(), input)
	if err == nil {
		t.Error("expected error, got nil")
	}
}

func TestUpdateCalendarSlotUpdateError(t *testing.T) {
	mockRepo := &mocks.MockAdminRepository{
		UpdateCalendarSlotError: errors.New("database error"),
	}

	service := createTestService(mockRepo)
	input := createValidInputSlot()

	_, err := service.UpdateCalendarSlot(context.Background(), uuid.New(), uuid.New(), input)
	if err == nil {
		t.Error("expected error, got nil")
	}
}

func TestUpdateCalendarSlotGetUpdateSlotsError(t *testing.T) {
	mockRepo := &mocks.MockAdminRepository{
		GetCalendarSlotsError: errors.New("database error"),
	}

	service := createTestService(mockRepo)
	input := createValidInputSlot()

	_, err := service.UpdateCalendarSlot(context.Background(), uuid.New(), uuid.New(), input)
	if err == nil {
		t.Error("expected error, got nil")
	}
}

func TestGetAvailableSlots(t *testing.T) {
	mockRepo := &mocks.MockAdminRepository{
		GetCalendarSlotsByWeekdaySlots: []dtos.CalendarSlotDB{
			createValidCalendarSlot(),
		},
		GetAppointmentsByDateAppointments: []dtos.AppointmentOutput{
			{
				StartTime: "10:00",
			},
		},
	}

	service := createTestService(mockRepo)

	slots, err := service.GetAvaliableSlots(context.Background(), uuid.New(), "2026-08-28")
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	for _, slot := range slots {
		if slot == "10:00" {
			t.Error("expected 10:00 to be umavailable")
		}
	}
}

func TestGetAvailableSlotsGetCalendarSlotsError(t *testing.T) {
	mockRepo := &mocks.MockAdminRepository{
		GetCalendarSlotsByWeekdayError: errors.New("database error"),
	}

	service := createTestService(mockRepo)

	_, err := service.GetAvaliableSlots(context.Background(), uuid.New(), "2026-08-28")
	if err == nil {
		t.Error("expected error, got nil")
	}
}

func TestGetAvailableSlotsGetAppointmentsError(t *testing.T) {
	mockRepo := &mocks.MockAdminRepository{
		GetAppointmentsByDateError: errors.New("database error"),
	}

	service := createTestService(mockRepo)

	_, err := service.GetAvaliableSlots(context.Background(), uuid.New(), "2028-08-28")
	if err == nil {
		t.Error("expected error, got nil")
	}
}

func TestGetAvailableSlotsInvalidDate(t *testing.T) {
	mockRepo := &mocks.MockAdminRepository{}
	service := createTestService(mockRepo)

	_, err := service.GetAvaliableSlots(context.Background(), uuid.New(), "data-invalida")
	if err == nil {
		t.Error("expected error, got nil")
	}
}

func TestCancelAppointmentByAdminGetAppointmentError(t *testing.T) {
	mockRepo := &mocks.MockAdminRepository{
		GetAllAppointmentByIDError: errors.New("database error"),
	}

	service := createTestService(mockRepo)

	err := service.CancelAppointmentByAdmin(context.Background(), uuid.New(), uuid.New())
	if err == nil {
		t.Error("expected error, got nil")
	}
}

func TestCancelAppointmentByAdminCancelError(t *testing.T) {
	mockRepo := &mocks.MockAdminRepository{
		CancelAppointmentByAdminError: errors.New("database error"),
	}

	service := createTestService(mockRepo)

	err := service.CancelAppointmentByAdmin(context.Background(), uuid.New(), uuid.New())
	if err == nil {
		t.Error("expected error, got nil")
	}
}

func TestUpdateAppointmentInvalidDate(t *testing.T) {
	mockRepo := &mocks.MockAdminRepository{}
	service := createTestService(mockRepo)

	input := dtos.AppointmentUpdateInput{
		Date: "data-invalida",
		StartTime: "10:00",
		EndTime: "11:00",
	}

	err := service.UpdateAppointment(context.Background(), uuid.New(), uuid.New(), input)
	if err == nil {
		t.Error("expected error, got nil")
	}
}

func TestUpdateAppointmentOutsideAvailability(t *testing.T) {
	appointmentID := uuid.New()

	mockRepo := &mocks.MockAdminRepository{
		GetAllAppointmentByIDAppointment: dtos.AppointmentDetails{
			ID: appointmentID,
		},
		GetCalendarSlotsByWeekdaySlots: []dtos.CalendarSlotDB{
			createValidCalendarSlot(),
		},
	}

	service := createTestService(mockRepo)

	input := dtos.AppointmentUpdateInput{
		Date: "2026-08-28",
		StartTime: "19:00",
		EndTime: "20:00",
	}

	err := service.UpdateAppointment(context.Background(), appointmentID, uuid.New(), input)
	if err == nil {
		t.Error("expected error, got nil")
	}
}

func TestUpdateAppointmentOccupiedTime(t *testing.T) {
	appointmentID := uuid.New()
	otherAppointmentID := uuid.New()

	mockRepo := &mocks.MockAdminRepository{
		GetAllAppointmentByIDAppointment: dtos.AppointmentDetails{
			ID:        appointmentID,
			Date:      "2026-08-28",
			StartTime: "08:00",
			EndTime:   "09:00",
		},
		GetCalendarSlotsByWeekdaySlots: []dtos.CalendarSlotDB{
			createValidCalendarSlot(),
		},
		GetAppointmentsByDateAppointments: []dtos.AppointmentOutput{
			{
				ID:        appointmentID,
				StartTime: "08:00",
				EndTime:   "09:00",
			},
			{
				ID:        otherAppointmentID,
				StartTime: "10:00",
				EndTime:   "11:00",
			},
		},
	}

	service := createTestService(mockRepo)

	input := dtos.AppointmentUpdateInput{
		Date:      "2026-08-28",
		StartTime: "10:30",
		EndTime:   "11:30",
	}

	err := service.UpdateAppointment(context.Background(), appointmentID, uuid.New(), input)
	if err == nil {
		t.Error("expected error, got nil")
	}
}

func TestUpdateAppointmentUpdateError(t *testing.T) {
	appointmentID := uuid.New()

	mockRepo := &mocks.MockAdminRepository{
		GetAllAppointmentByIDAppointment: dtos.AppointmentDetails{
			ID: appointmentID,
		},
		GetCalendarSlotsByWeekdaySlots: []dtos.CalendarSlotDB{
			createValidCalendarSlot(),
		},
		UpdateCalendarSlotError: errors.New("database error"),
	}

	service := createTestService(mockRepo)

	input := dtos.AppointmentUpdateInput{
		Date:      "2026-08-28",
		StartTime: "10:00",
		EndTime:   "11:00",
	}

	err := service.UpdateAppointment(context.Background(), appointmentID, uuid.New(), input)
	if err == nil {
		t.Error("expected error, got nil")
	}
}

func TestUpdateAppointmentGetAppointmentsError(t *testing.T){
	appointmentID := uuid.New()

	mockRepo := &mocks.MockAdminRepository{
		GetAllAppointmentByIDAppointment: dtos.AppointmentDetails{
			ID: appointmentID,
		},
		GetCalendarSlotsByWeekdaySlots: []dtos.CalendarSlotDB{
			createValidCalendarSlot(),
		},
		GetAppointmentsByDateError: errors.New("database error"),
	}

	service := createTestService(mockRepo)

	input := dtos.AppointmentUpdateInput{
		Date:      "2026-08-28",
		StartTime: "10:00",
		EndTime:   "11:00",
	}

	err := service.UpdateAppointment(context.Background(), appointmentID, uuid.New(), input)
	if err == nil {
		t.Error("expected error, got nil")
	}
}