package utils

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/jhonnydsl/clinify-backend/src/dtos"
)

func TestAPIError(t *testing.T) {
	tests := []struct {
		name       string
		fn         func(string) *dtos.APIError
		message    string
		wantStatus int
	}{
		{
			name:       "not found",
			fn:         NotFoundError,
			message:    "user not found",
			wantStatus: http.StatusNotFound,
		},
		{
			name:       "bad request",
			fn:         BadRequestError,
			message:    "invalid input",
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "conflict",
			fn:         ConflictError,
			message:    "email already exists",
			wantStatus: http.StatusConflict,
		},
		{
			name:       "internal server error",
			fn:         InternalServerError,
			message:    "database error",
			wantStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.fn(tt.message)

			if err == nil {
				t.Fatal("expected API error, got nil")
			}

			if err.StatusCode != tt.wantStatus {
				t.Errorf("expected status %d, got %d", tt.wantStatus, err.StatusCode)
			}

			if err.Message != tt.message {
				t.Errorf("expected message %q, got %q", tt.message, err.Message)
			}
		})
	}
}

func TestGetStatusCode(t *testing.T) {
	tests := []struct {
		name       string
		err        error
		wantStatus int
	}{
		{
			name:       "api error",
			err:        NotFoundError("user not found"),
			wantStatus: http.StatusNotFound,
		},
		{
			name:       "bad request api error",
			err:        BadRequestError("invalid input"),
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "generic error",
			err:        fmt.Errorf("something went wrong"),
			wantStatus: http.StatusInternalServerError,
		},
		{
			name:       "nil error",
			err:        nil,
			wantStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetStatusCode(tt.err)

			if got != tt.wantStatus {
				t.Errorf("expected status %d, got %d", tt.wantStatus, got)
			}
		})
	}
}