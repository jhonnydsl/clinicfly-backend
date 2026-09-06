package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jhonnydsl/clinify-backend/src/dtos"
	"github.com/jhonnydsl/clinify-backend/src/services"
	"github.com/jhonnydsl/clinify-backend/src/utils"
)

type PatientController struct {
	Service *services.PatientService
	AdminService *services.AdminService
}

func (controller *PatientController) CreatePatient(c *gin.Context) {
	var patientInput dtos.PatientInput

	ctx, cancel := utils.NewDBContext()
	defer cancel()
	
	err := c.ShouldBindJSON(&patientInput)
	if err != nil {
		c.JSON(utils.GetStatusCode(err), gin.H{"error": err.Error()})
		return
	}
	
	id, err := controller.Service.CreatePatient(ctx, patientInput)
	if err != nil {
		c.JSON(utils.GetStatusCode(err), gin.H{"error": err.Error()})
		return
	}

	c.JSON(201, gin.H{
		"message": "user patient created",
		"id": 		id,
	})
}

func (controller *PatientController) GetAppointments(c *gin.Context) {
	ctx, cancel := utils.NewDBContext()
	defer cancel()

	patientIDValue, exists := c.Get("id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "patient id not found in context"})
		return
	}

	patientID, err := uuid.Parse(patientIDValue.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid patient id"})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	appointments, err := controller.Service.GetAppointments(ctx, patientID, page, limit)
	if err != nil {
		c.JSON(utils.GetStatusCode(err), gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, appointments)
}

func (controller *PatientController) GetDoctorAvailableSlots(c *gin.Context) {
	adminIDStr := c.Param("id")

	adminID, err := uuid.Parse(adminIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid doctor id"})
		return
	}

	date := c.Query("date")
	if date == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "date is required"})
		return
	}

	ctx, cancel := utils.NewDBContext()
	defer cancel()

	availableSlots, err := controller.Service.GetDoctorAvaliableSlots(ctx, adminID, date)
	if err != nil {
		c.JSON(utils.GetStatusCode(err), gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, availableSlots)
}

func (controller *PatientController) CreateAppointment(c *gin.Context) {
	var input dtos.AppointmentInput

	ctx, cancel := utils.NewDBContext()
	defer cancel()

	adminIDStr := c.Param("id")

	adminID, err := uuid.Parse(adminIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid admin id"})
		return
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(utils.GetStatusCode(err), gin.H{"error": err.Error()})
		return
	}

	patientIDValue, exists := c.Get("id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "patient id not found in context"})
		return
	}

	patientID, err := uuid.Parse(patientIDValue.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid patient id"})
		return
	}

	id, err := controller.AdminService.CreateAppointment(ctx, input, adminID, patientID, "patient", c.ClientIP(), c.GetHeader("User-Agent"))
	if err != nil {
		c.JSON(utils.GetStatusCode(err), gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "appointment created",
		"id": id,
	})
}

func (controller *PatientController) CancelAppointmentByPatient(c *gin.Context) {
	ctx, cancel := utils.NewDBContext()
	defer cancel()

	patientIDValue, exists := c.Get("id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "patient id not found in context"})
		return
	}

	patientID, err := uuid.Parse(patientIDValue.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid patient id"})
		return
	}

	appointmentIDStr := c.Param("id")

	appointmentID, err := uuid.Parse(appointmentIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid appointment id"})
		return
	}

	err = controller.Service.CancelAppointmentByPatient(ctx, appointmentID, patientID, patientID, "patient", c.ClientIP(), c.GetHeader("User-Agent"))
	if err != nil {
		c.JSON(utils.GetStatusCode(err), gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "appointment cancelled successfully",
	})
}