package utils

import (
	"os"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func TestGenerateToken(t *testing.T) {
	os.Setenv("JWT_SECRET", "secret-de-teste")

	token, err := GenerateJWT(
		"123",
		"Jhonny Lima",
		"jhonny@gmail.com",
		"admin",
	)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if token == "" {
		t.Error("expected token, got empty string")
	}

	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte("secret-de-teste"), nil
	})
	if err != nil {
		t.Fatalf("failed to parse token: %v", err)
	}

	if !parsedToken.Valid {
		t.Fatal("expected token to be valid")
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		t.Fatal("expected jwt.MapClaims")
	}

	if claims["id"] != "123" {
		t.Errorf("expected id 123, got %v", claims["id"])
	}

	if claims["full_name"] != "Jhonny Lima" {
		t.Errorf("expected full_name Jhonny Lima, got %v", claims["full_name"])
	}

	if claims["email"] != "jhonny@gmail.com" {
		t.Errorf("expected email jhonny@gmail.com, got %v", claims["email"])
	}

	if claims["role"] != "admin" {
		t.Errorf("expected role admin, got %v", claims["role"])
	}

	exp, ok := claims["exp"].(float64)
	if !ok {
		t.Fatalf("expected exp to be float64")
	}

	expiration := time.Unix(int64(exp), 0)

	expected := time.Now().Add(24 * time.Hour)

	if expiration.Before(expected.Add(-time.Minute)) || expiration.After(expected.Add(time.Minute)) {
		t.Errorf("unexpected expiration time: %v", expiration)
	}
}