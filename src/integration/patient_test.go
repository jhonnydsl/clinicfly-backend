package integration

import (
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
	"github.com/jhonnydsl/clinify-backend/src/mailer"
	"github.com/jhonnydsl/clinify-backend/src/repository"
	"github.com/jhonnydsl/clinify-backend/src/routes"
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