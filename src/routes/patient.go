package routes

import (
	"github.com/gin-gonic/gin"
	auditrepository "github.com/jhonnydsl/clinify-backend/src/audit/repository"
	auditservices "github.com/jhonnydsl/clinify-backend/src/audit/services"
	"github.com/jhonnydsl/clinify-backend/src/controllers"
	"github.com/jhonnydsl/clinify-backend/src/mailer"
	"github.com/jhonnydsl/clinify-backend/src/repository"
	"github.com/jhonnydsl/clinify-backend/src/services"
	"github.com/jhonnydsl/clinify-backend/src/utils/middlewares"
)

func SetupPatientRoutes(app *gin.RouterGroup, mailer *mailer.Mailer) {
	patientService := &services.PatientService{
		Repo: &repository.PatientRepository{},
		AdminRepo: &repository.AdminRepository{},
		Mailer: mailer,
	}
	adminService := &services.AdminService{
		Repo: &repository.AdminRepository{},
		Mailer: mailer,
		AuditService: &auditservices.AuditService{
			Repo: &auditrepository.AuditRepository{},
		},
	}

	patientController := &controllers.PatientController{
		Service: patientService,
		AdminService: adminService,
	}

	patient := app.Group("/patient")
	{
		patient.POST("", patientController.CreatePatient)
		patient.GET("/doctors/:id/available-slots", patientController.GetDoctorAvailableSlots)
	}

	protectedPatient := app.Group(
		"/patient",
		middlewares.AuthMiddleware(),
		middlewares.PatientOnlyMiddleware(),
	)
	{
		protectedPatient.POST("/doctors/:id/appointments", patientController.CreateAppointment)
		protectedPatient.PATCH("/appointments/:id/cancel", patientController.CancelAppointmentByPatient)
		protectedPatient.GET("/appointments", patientController.GetAppointments)
	}
}