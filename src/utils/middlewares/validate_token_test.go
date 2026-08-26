package middlewares

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jhonnydsl/clinify-backend/src/utils"
)

func executeAuthMiddlewareTest(t *testing.T, claims jwt.MapClaims) (int, bool) {
	t.Helper()

	t.Setenv("JWT_SECRET", "secret-de-teste")

	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		claims,
	)

	tokenString, err := token.SignedString([]byte("secret-de-teste"))
	if err != nil {
		t.Fatalf("failed to generate token: %v", err)
	}

	router := gin.New()
	router.Use(AuthMiddleware())

	called := false

	router.GET("/test", func(c *gin.Context) {
		called = true
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Authorization", "Bearer "+tokenString)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	return w.Code, called
}

func TestAuthMiddlewareMissingAuthorization(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.Use(AuthMiddleware())

	called := false

	router.GET("/test", func(c *gin.Context) {
		called = true
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected status 401, got %d", w.Code)
	}

	if called {
		t.Error("expected handler not to be called")
	}
}

func TestAuthMiddlewareOptions(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.Use(AuthMiddleware())

	called := false

	router.OPTIONS("/test", func(c *gin.Context) {
		called = true
	})

	req := httptest.NewRequest(http.MethodOptions, "/test", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusNoContent {
		t.Errorf("expected status 204, got %d", w.Code)
	}

	if called {
		t.Error("expected handler not to be called")
	}
}

func TestAuthMiddlewareValidToken(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Setenv("JWT_SECRET", "secret-de-teste")

	token, err := utils.GenerateJWT(
		"123",
		"Jhonny Lima",
		"jhonny@gmail.com",
		"admin",
	)
	if err != nil {
		t.Fatalf("failed to generate token: %v", err)
	}

	router := gin.New()
	router.Use(AuthMiddleware())

	router.GET("/test", func(c *gin.Context) {
		if c.GetString("id") != "123" {
			t.Errorf("expected id 123, got %q", c.GetString("id"))
		}

		if c.GetString("role") != "admin" {
			t.Errorf("expected role admin, got %q", c.GetString("role"))
		}
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}
}

func TestAuthMiddlewareInvalidToken(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Setenv("JWT_SECRET", "secret-de-teste")

	router := gin.New()
	router.Use(AuthMiddleware())

	called := false

	router.GET("/test", func(c *gin.Context) {
		called = true
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Authorization", "Bearer token-invalido")

	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected status 401, got %d", w.Code)
	}

	if called {
		t.Error("expected handler not be called")
	}
}

func TestAuthMiddlewareUnexpectedSigningMethod(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Setenv("JWT_SECRET", "secret-de-teste")

	token := jwt.NewWithClaims(
		jwt.SigningMethodNone,
		jwt.MapClaims{
			"id": "123",
			"role": "admin",
		},
	)

	tokenString, err := token.SignedString(jwt.UnsafeAllowNoneSignatureType)
	if err != nil {
		t.Fatalf("failed to generate token: %v", err)
	}

	router := gin.New()
	router.Use(AuthMiddleware())

	called := false

	router.GET("/test", func(c *gin.Context) {
		called = true
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Authorization", "Bearer "+tokenString)

	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected status 401, got %d", w.Code)
	}

	if called {
		t.Error("expected handler not be called")
	}
}

func TestAuthMiddlewareInvalidID(t *testing.T) {
	status, called := executeAuthMiddlewareTest(t, jwt.MapClaims{
		"role": "admin",
	})

	if status != http.StatusUnauthorized {
		t.Errorf("expected status 401, got %d", status)
	}

	if called {
		t.Error("expected handler not to be called")
	}
}

func TestAuthMiddlewareInvalidIDType(t *testing.T) {
	status, called := executeAuthMiddlewareTest(t, jwt.MapClaims{
		"id": 123,
		"role": "admin",
	})

	if status != http.StatusUnauthorized {
		t.Errorf("expected status 401, got %d", status)
	}

	if called {
		t.Error("expected handler not to be called")
	}
}