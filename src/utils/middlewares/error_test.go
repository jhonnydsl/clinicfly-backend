package middlewares

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/jhonnydsl/clinify-backend/src/dtos"
)

func executeErrorMiddlewareTest(t *testing.T, err error) (int, string) {
	t.Helper()

	router := gin.New()
	router.Use(ErrorMiddlewareHandle())

	router.GET("/test", func(c *gin.Context) {
		c.Error(err)
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	return w.Code, w.Body.String()
}

func TestErrorMiddlewareHandleNoError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.Use(ErrorMiddlewareHandle())

	called := false

	router.GET("/test", func(c *gin.Context) {
		called = true
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	if !called {
		t.Error("expected handler to be called")
	}
}

func TestErrorMiddlewareHandleAPIError(t *testing.T) {
	status, body := executeErrorMiddlewareTest(t, &dtos.APIError{
		StatusCode: http.StatusNotFound,
		Message: "patient not found",
	})

	if status != http.StatusNotFound {
		t.Errorf("expected status 404, got %d", status)
	}

	if body != `{"error":"patient not found"}` {
		t.Errorf("unexpected response body: %s", body)
	}
}

func TestErrorMiddlewareHandleInternalError(t *testing.T) {
	status, body := executeErrorMiddlewareTest(t, errors.New("database connection failed"))
	if status != http.StatusInternalServerError {
		t.Errorf("expected status 500, got %d", status)
	}

	if body != `{"error":"error internal server"}` {
		t.Errorf("unexpected response body: %s", body)
	}
}