package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/jhonnydsl/clinify-backend/src/controllers"
	"github.com/jhonnydsl/clinify-backend/src/mailer"
	"github.com/jhonnydsl/clinify-backend/src/repository"
	"github.com/jhonnydsl/clinify-backend/src/services"
	"github.com/jhonnydsl/clinify-backend/src/utils/middlewares"
)

func SetupAdminRoutes(app *gin.RouterGroup, mailer *mailer.Mailer) {
	adminService := &services.AdminService{Repo: &repository.AdminRepository{}, Mailer: mailer}
	adminController := &controllers.AdminController{Service: adminService}

	app.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Server online"})
	})

	admin := app.Group("/admin")
	{
		admin.POST("", adminController.CreateAdmin)
	}

	protectedAdmin := app.Group("/admin", middlewares.AuthMiddleware(), middlewares.AdminOnlyMiddleware())
	{
		protectedAdmin.POST("/appointments", adminController.CreateAppointment)
		protectedAdmin.GET("/patients", adminController.GetPatients)			// => rota correta com paginação GET /api/v1/admin/patients?page=1&limit=10
		protectedAdmin.GET("/appointments", adminController.GetAppointments)	// => rota correta com paginação GET /api/v1/admin/appointments?page=1&limit=10
		protectedAdmin.DELETE("patients/:id", adminController.DeletePatient)
		protectedAdmin.POST("/calendar-slots", adminController.CreateCalendarSlot)
		protectedAdmin.GET("/calendar-slots", adminController.GetCalendarSlots)
		protectedAdmin.DELETE("/calendar-slots/:id", adminController.DeleteCalendarSlot)
	}
}