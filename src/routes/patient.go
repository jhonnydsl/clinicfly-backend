package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/jhonnydsl/clinify-backend/src/controllers"
	"github.com/jhonnydsl/clinify-backend/src/mailer"
	"github.com/jhonnydsl/clinify-backend/src/repository"
	"github.com/jhonnydsl/clinify-backend/src/services"
)

func SetupPatientRoutes(app *gin.RouterGroup, mailer *mailer.Mailer) {
	patientService := &services.PatientService{
		Repo: &repository.PatientRepository{},
		AdminRepo: &repository.AdminRepository{},
	}
	adminService := &services.AdminService{
		Repo: &repository.AdminRepository{},
		Mailer: mailer,
	}

	patientController := &controllers.PatientController{
		Service: patientService,
		AdminService: adminService,
	}

	patient := app.Group("/patient")
	{
		patient.POST("", patientController.CreatePatient)
		patient.GET("doctors/:id/available-slots", patientController.GetDoctorAvailableSlots)

		patient.POST("/doctors/:id/appointments", patientController.CreateAppointment)
	}
}