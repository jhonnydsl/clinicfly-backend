package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/jhonnydsl/clinify-backend/src/controllers"
	"github.com/jhonnydsl/clinify-backend/src/mailer"
	"github.com/jhonnydsl/clinify-backend/src/repository"
	"github.com/jhonnydsl/clinify-backend/src/services"
)

func SetupAuthRoutes(app *gin.RouterGroup, mailer *mailer.Mailer) {
	authService := &services.AuthService{Repo: &repository.AuthRepository{}, Mailer: mailer}
	authController := &controllers.AuthController{Service: authService}

	app.POST("/forgot-password", authController.ForgotPassword)
	app.POST("/reset-password", authController.ResetPassword)
}