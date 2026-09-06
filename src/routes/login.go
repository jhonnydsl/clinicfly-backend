package routes

import (
	"github.com/gin-gonic/gin"
	auditservices "github.com/jhonnydsl/clinify-backend/src/audit/services"
	"github.com/jhonnydsl/clinify-backend/src/controllers"
	"github.com/jhonnydsl/clinify-backend/src/repository"
	"github.com/jhonnydsl/clinify-backend/src/services"
)

func SetupLoginRoutes(app *gin.RouterGroup) {
	auditServices := &auditservices.AuditService{}

	loginService := &services.LoginService{
		Repo: &repository.LoginRepository{},
		AuditService: auditServices,
	}
	loginController := &controllers.LoginController{Service: loginService}

	app.POST("/login", loginController.LoginUser)
}