package controllers

import (
	"math"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jhonnydsl/clinify-backend/src/dtos"
	"github.com/jhonnydsl/clinify-backend/src/services"
	"github.com/jhonnydsl/clinify-backend/src/utils"
)

type AdminController struct {
	Service *services.AdminService
}

func (controller *AdminController) CreateAdmin(c *gin.Context) {
	var adminInput dtos.AdminInput

	ctx, cancel := utils.NewDBContext()
	defer cancel()

	err := c.ShouldBindJSON(&adminInput)
	if err != nil {
		c.JSON(utils.GetStatusCode(err), gin.H{"error": err.Error()})
		return
	}

	id, err := controller.Service.CreateAdmin(ctx, adminInput)
	if err != nil {
		c.JSON(utils.GetStatusCode(err), gin.H{"error": err.Error()})
		return
	}

	c.JSON(201, gin.H{
		"message": "user admin created",
		"id":		id,
	})
}

func (controller *AdminController) CreateAppointment(c *gin.Context) {
	var input dtos.AppointmentInput

	ctx, cancel := utils.NewDBContext()
	defer cancel()

	clientIDValue, exists := c.Get("id")
	if !exists {
		c.JSON(401, gin.H{"error": "client id not found in context"})
		return
	}

	clientID, err := uuid.Parse(clientIDValue.(string))
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid client id"})
		return
	}

	err = c.ShouldBindJSON(&input)
	if err != nil {
		c.JSON(utils.GetStatusCode(err), gin.H{"error": err.Error()})
		return
	}

	id, err := controller.Service.CreateAppointment(ctx, input, clientID)
	if err != nil {
		c.JSON(utils.GetStatusCode(err), gin.H{"error": err.Error()})
		return
	}

	c.JSON(201, gin.H{
		"message": "appointment created",
		"id": 		id,
	})
}

func (controller *AdminController) GetPatients(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	adminIDStr, exists := c.Get("id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid client id"})
		return
	}

	adminID, err := uuid.Parse(adminIDStr.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid client id"})
		return
	}

	ctx, cancel := utils.NewDBContext()
	defer cancel()

	patients, total, err := controller.Service.GetPatients(ctx, adminID, page, limit)
	if err != nil {
		c.JSON(utils.GetStatusCode(err), gin.H{"error": err.Error()})
		return
	}

	totalPages := int(math.Ceil(float64(total) / float64(limit)))

	c.JSON(http.StatusOK, gin.H{
		"data": patients,
		"page": page,
		"limit": limit,
		"total": total,
		"total_pages": totalPages,
	})
}

func (controller *AdminController) DeletePatient(c *gin.Context) {
	ctx, cancel := utils.NewDBContext()
	defer cancel()

	idParam := c.Param("id")

	patientID, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid patient id"})
		return
	}

	err = controller.Service.DeletePatient(ctx, patientID)
	if err != nil {
		c.JSON(utils.GetStatusCode(err), gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "patient deleted successfully"})
}