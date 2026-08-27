package services_test

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/jhonnydsl/clinify-backend/src/dtos"
	"github.com/jhonnydsl/clinify-backend/src/mocks"
	"github.com/jhonnydsl/clinify-backend/src/services"
	"golang.org/x/crypto/bcrypt"
)

func TestCreateAdminInvalidInput(t *testing.T) {
	mockRepo := &mocks.MockAdminRepository{}

	service := services.AdminService{
		Repo: mockRepo,
	}

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

	service := services.AdminService{
		Repo: mockRepo,
	}

	admin := dtos.AdminInput{
		FullName: "Jhonny Lima",
		Email: "jhonny@gmail.com",
		Password: "123456",
		BirthDate: "1990-01-01",
		Phone: "11999999999",
	}

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

	service := services.AdminService{
		Repo: mockRepo,
	}

	admin := dtos.AdminInput{
		FullName: "Jhonny Lima",
		Email: "jhonny@gmail.com",
		Password: "123456",
		BirthDate: "1990-01-01",
		Phone: "11999999999",
	}

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