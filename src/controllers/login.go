package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jhonnydsl/clinify-backend/src/dtos"
	"github.com/jhonnydsl/clinify-backend/src/services"
	"github.com/jhonnydsl/clinify-backend/src/utils"
)

type LoginController struct {
	Service *services.LoginService
}

func (controller *LoginController) LoginUser(c *gin.Context) {
	var loginInput dtos.LoginInput
	
	err := c.ShouldBindJSON(&loginInput)
	if err != nil {
		c.JSON(utils.GetStatusCode(err), gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := utils.NewDBContext()
	defer cancel()

	login, err := controller.Service.LoginUser(ctx, loginInput.Email, loginInput.Password, c.ClientIP(), c.GetHeader("User-Agent"))
	if err != nil {
		c.JSON(utils.GetStatusCode(err), gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": login})
}