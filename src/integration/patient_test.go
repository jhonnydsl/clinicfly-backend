package integration

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jhonnydsl/clinify-backend/src/dtos"
	"github.com/jhonnydsl/clinify-backend/src/mailer"
	"github.com/jhonnydsl/clinify-backend/src/repository"
	"github.com/jhonnydsl/clinify-backend/src/routes"
	"github.com/jhonnydsl/clinify-backend/src/utils"
)

func setupPatientRouter() *gin.Engine {
	router := gin.Default()

	api := router.Group("/api/v1")

	testMailer := mailer.NewMailer(
		os.Getenv("SMTP_EMAIL"),
		os.Getenv("SMTP_PASSWORD"),
	)

	routes.SetupPatientRoutes(api, testMailer)

	return router
}

func createTestCalendarSlot(t *testing.T, adminID uuid.UUID, weekday int, start, end string) {
	t.Helper()

	repo := &repository.AdminRepository{}

	input := dtos.CalendarSlotsInput{
		Weekday:   weekday,
		StartTime: start,
		EndTime:   end,
	}

	startTime, err := utils.ParseTime(input.StartTime)
	if err != nil {
		t.Fatalf("failed to parse start time: %v", err)
	}

	endTime, err := utils.ParseTime(input.EndTime)
	if err != nil {
		t.Fatalf("failed to parse end time: %v", err)
	}

	_, err = repo.CreateCalendarSlot(
		context.Background(),
		input,
		startTime,
		endTime,
		adminID,
	)
	if err != nil {
		t.Fatalf("failed to create test calendar slot: %v", err)
	}
}

func TestCreatePatientIntegration(t *testing.T) {
	adminID, _, adminSlug := createTestAdmin(t)

	router := setupPatientRouter()

	body := fmt.Sprintf(`{
		"full_name": "Patient Integration Test",
		"email": "patient-integration-%d@test.com",
		"password_hash": "123456",
		"phone": "11988888888",
		"birth_date": "2000-01-01",
		"public_slug": "%s"
	}`, time.Now().UnixNano(), adminSlug)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/patient", strings.NewReader(body))

	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected status 201, got %d: %s", rec.Code, rec.Body.String())
	}

	var response struct {
		ID uuid.UUID `json:"id"`
	}

	if err := json.Unmarshal(rec.Body.Bytes(), &response); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	t.Logf("patient created: %s", response.ID)

	var dbClientID uuid.UUID
	var passwordHash string

	err := repository.DB.QueryRow(`
		SELECT client_id, password_hash
		FROM patients
		WHERE id = $1
	`, response.ID).Scan(&dbClientID, &passwordHash)

	if err != nil {
		t.Fatalf("failed to query created patient: %v", err)
	}

	if dbClientID != adminID {
		t.Fatalf("expected client_id %s, got %s", adminID, dbClientID)
	}

	if passwordHash == "123456" {
		t.Fatal("password was stored in plain text")
	}

	t.Logf("patient correctly stored with hashed password")
}

func TestGetDoctorAvailableSlotsIntegration(t *testing.T) {
	adminID, _, _ := createTestAdmin(t)

	createTestCalendarSlot(t, adminID, 1, "08:00", "10:00")

	router := setupPatientRouter()

	req := httptest.NewRequest(
		http.MethodGet,
		fmt.Sprintf("/api/v1/patient/doctors/%s/available-slots?date=2026-09-07", adminID),
		nil,
	)
	
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", rec.Code, rec.Body.String())
	}

	var slots []string

	if err := json.Unmarshal(rec.Body.Bytes(), &slots); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	expected := []string{"08:00", "09:00"}

	if len(slots) != len(expected) {
		t.Fatalf("expected %d available slots, got %d: %v", len(expected), len(slots), slots)
	}

	for i := range expected {
		if slots[i] != expected[i] {
			t.Fatalf("expected slot %s at position %d, got %s", expected[i], i, slots[i])
		}
	}

	t.Logf("available slots correctly returned: %v", slots)
}

func TestGetDoctorAvailableSlotsWithAppointmentIntegration(t *testing.T) {
	adminID, _, adminSlug := createTestAdmin(t)

	patientID := createTestPatient(t, adminSlug)

	createTestCalendarSlot(t, adminID, 1, "08:00", "10:00")

	repo := &repository.AdminRepository{}

	appointmentInput := dtos.AppointmentInput{
		PatientID: patientID.String(),
		Date:      "2026-09-07",
		StartTime: "08:00",
		EndTime:   "09:00",
	}

	parsedDate, err := utils.ParseDate(appointmentInput.Date)
	if err != nil {
		t.Fatalf("failed to parse appointment date: %v", err)
	}

	appointmentStart, err := utils.ParseDateTime(
		appointmentInput.Date,
		appointmentInput.StartTime,
	)
	if err != nil {
		t.Fatalf("failed to parse appointment start: %v", err)
	}

	appointmentEnd, err := utils.ParseDateTime(
		appointmentInput.Date,
		appointmentInput.EndTime,
	)
	if err != nil {
		t.Fatalf("failed to parse appointment end: %v", err)
	}

	_, err = repo.CreateAppointment(
		context.Background(),
		appointmentInput,
		parsedDate,
		appointmentStart,
		appointmentEnd,
		adminID,
	)
	if err != nil {
		t.Fatalf("failed to create test appointment: %v", err)
	}

	router := setupPatientRouter()

	req := httptest.NewRequest(
		http.MethodGet,
		fmt.Sprintf("/api/v1/patient/doctors/%s/available-slots?date=2026-09-07", adminID),
		nil,
	)

	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", rec.Code, rec.Body.String())
	}

	var slots []string

	if err := json.Unmarshal(rec.Body.Bytes(), &slots); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	expected := []string{"09:00"}

	if len(slots) != len(expected) {
		t.Fatalf("expected %d available slots, got %d: %v", len(expected), len(slots), slots)
	}

	if slots[0] != expected[0] {
		t.Fatalf("expected slot %s, got %s", expected[0], slots[0])
	}

	t.Logf("occupied slot correctly removed: %v", slots)
}