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
	"github.com/jhonnydsl/clinify-backend/src/repository"
	"github.com/jhonnydsl/clinify-backend/src/routes"
	"github.com/jhonnydsl/clinify-backend/src/utils"
	"github.com/joho/godotenv"
)

func TestMain(m *testing.M) {
	if err := godotenv.Load("../../.env"); err != nil {
		panic(err)
	}

	if err := repository.Connect(); err != nil {
		panic(err)
	}

	code := m.Run()

	repository.DB.Close()

	os.Exit(code)
}

func setupRouter() *gin.Engine {
	router := gin.Default()

	api := router.Group("/api/v1")

	routes.SetupAdminRoutes(api, nil)

	return router
}

func createTestAdmin(t *testing.T) (uuid.UUID, string) {
	t.Helper()

	repo := &repository.AdminRepository{}

	admin := dtos.AdminInput{
		FullName:      "Admin Teste",
		Email:         fmt.Sprintf("admin-%d@test.com", time.Now().UnixNano()),
		Password:      "123456",
		BirthDate:     "2003-02-03",
		Crp:           "TEST123",
		Bio:           "Admin de teste",
		OfficeAddress: "Rua Teste, 123",
		Phone:         "11999999999",
		PublicSlug:    fmt.Sprintf("admin-teste-%d", time.Now().UnixNano()),
	}

	birthDate, err := utils.ParseDate(admin.BirthDate)
	if err != nil {
		t.Fatalf("error creating birth date: %v", err)
	}

	id, err := repo.CreateAdmin(context.Background(), admin, birthDate)
	if err != nil {
		t.Fatalf("error creating admin test: %v", err)
	}

	return id, admin.Email
}

func createTestAdminToken(t *testing.T) string {
	t.Helper()

	adminID, adminEmail := createTestAdmin(t)

	token, err := utils.GenerateJWT(adminID.String(), "Admin Teste", adminEmail, "admin")
	if err != nil {
		t.Fatalf("error generating token: %v", err)
	}

	return token
}

func postCalendarSlot(t *testing.T, router *gin.Engine, token, body string) *httptest.ResponseRecorder {
	t.Helper()

	req := httptest.NewRequest(http.MethodPost, "/api/v1/admin/calendar-slots", strings.NewReader(body))

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	return rec
}

func TestCreateCalendarSlotIntegration(t *testing.T) {
	token := createTestAdminToken(t)

	router := setupRouter()

	body := `{
		"weekday": 1,
		"start_time": "08:00",
		"end_time": "09:00"
	}`

	rec := postCalendarSlot(t, router, token, body)

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected status 201, got %d: %s", rec.Code, rec.Body.String())
	}

	var response struct {
		ID uuid.UUID `json:"id"`
	}

	if err := json.Unmarshal(rec.Body.Bytes(), &response); err != nil {
		t.Fatalf("error decode response: %v", err)
	}

	t.Logf("slot created: %s", response.ID)

	var weekday int
	var startTime, endTime time.Time

	err := repository.DB.QueryRow(`
		SELECT weekday, start_time, end_time
		FROM calendar_slots
		WHERE id = $1
	`, response.ID).Scan(&weekday, &startTime, &endTime)

	if err != nil {
		t.Fatalf("failed to query created slot: %v", err)
	}

	if weekday != 1 {
		t.Fatalf("expected weekday 1, got %d", weekday)
	}

	 if startTime.Format("15:04") != "08:00" {
		t.Fatalf("expected start time 08:00, got %s", startTime.Format("15:04"))
	 }

	 if endTime.Format("15:04") != "09:00" {
		t.Fatalf("expected end time 09:00, got %s", endTime.Format("15:04"))
	 }

	 t.Logf("slot confirmed in database: weekday=%d, start=%s, end=%s", weekday, startTime.Format("15:04"), endTime.Format("15:04"))
}

func TestCreateCalendarSlotOverlapIntegration(t *testing.T) {
	token := createTestAdminToken(t)

	router := setupRouter()

	firstBody := `{
		"weekday": 1,
		"start_time": "08:00",
		"end_time": "09:00"
	}`

	rec := postCalendarSlot(t, router, token, firstBody)

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected first slot creation to return 201, got %d: %s", rec.Code, rec.Body.String())
	}

	secondeBody := `{
		"weekday": 1,
		"start_time": "08:30",
		"end_time": "09:30"
	}`

	rec = postCalendarSlot(t, router, token, secondeBody)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected overlapping slot creation to return 400, got %d: %s", rec.Code, rec.Body.String())
	}
}