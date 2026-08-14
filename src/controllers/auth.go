package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jhonnydsl/clinify-backend/src/dtos"
	"github.com/jhonnydsl/clinify-backend/src/services"
	"github.com/jhonnydsl/clinify-backend/src/utils"
)

type AuthController struct {
	Service *services.AuthService
}

func (controller *AuthController) ForgotPassword(c *gin.Context) {
	var input dtos.ForgotPasswordInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	ctx, cancel := utils.NewDBContext()
	defer cancel()

	err := controller.Service.ForgotPassword(ctx, input.Email)
	if err != nil {
		c.JSON(utils.GetStatusCode(err), gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "if the email exists, a password reset link will be sent",
	})
}

func (controller *AuthController) ResetPassword(c *gin.Context) {
	var input dtos.ResetPasswordInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	ctx, cancel := utils.NewDBContext()
	defer cancel()

	err := controller.Service.ResetPassword(ctx, input)
	if err != nil {
		c.JSON(utils.GetStatusCode(err), gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "password reset successfully",
	})
}