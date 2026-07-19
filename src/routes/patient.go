package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/jhonnydsl/clinify-backend/src/controllers"
	"github.com/jhonnydsl/clinify-backend/src/repository"
	"github.com/jhonnydsl/clinify-backend/src/services"
)

func SetupPatientRoutes(app *gin.RouterGroup) {
	patientService := &services.PatientService{
		Repo: &repository.PatientRepository{},
		AdminRepo: &repository.AdminRepository{},
	}
	patientController := &controllers.PatientController{Service: patientService}

	patient := app.Group("/patient")
	{
		patient.POST("", patientController.CreatePatient)
		patient.GET("doctors/:id/available-slots", patientController.GetDoctorAvailableSlots)
	}
}