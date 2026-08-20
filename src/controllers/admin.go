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

func (controller *AdminController) GetAppointments(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	status := c.Query("status")

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

	appointments, total, err := controller.Service.GetAppointments(ctx, adminID, status, page, limit)
	if err != nil {
		c.JSON(utils.GetStatusCode(err), gin.H{"error": err.Error()})
		return
	}

	totalPages := int(math.Ceil(float64(total) / float64(limit)))

	c.JSON(http.StatusOK, gin.H{
		"data": appointments,
		"status": status,
		"page": page,
		"limit": limit,
		"total": total,
		"total_pages": totalPages,
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
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid patient id"})
		return
	}

	err = controller.Service.DeletePatient(ctx, patientID)
	if err != nil {
		c.JSON(utils.GetStatusCode(err), gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "patient deleted successfully"})
}

func (controller *AdminController) CreateCalendarSlot(c *gin.Context) {
	var input dtos.CalendarSlotsInput

	ctx, cancel := utils.NewDBContext()
	defer cancel()
	
	err := c.ShouldBindJSON(&input)
	if err != nil {
		c.JSON(utils.GetStatusCode(err), gin.H{"error": err.Error()})
		return
	}
	
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


	id, err := controller.Service.CreateCalendarSlot(ctx, input, clientID)
	if err != nil {
		c.JSON(utils.GetStatusCode(err), gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "slot created succesfully",
		"id": 		id,
	})
}

func (controller *AdminController) GetCalendarSlots(c *gin.Context) {
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

	slotsOutputs, err := controller.Service.GetCalendarSlots(ctx, adminID)
	if err != nil {
		c.JSON(utils.GetStatusCode(err), gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, slotsOutputs)
}

func (controller *AdminController) DeleteCalendarSlot(c *gin.Context) {
	ctx, cancel := utils.NewDBContext()
	defer cancel()

	idParam := c.Param("id")

	slotID, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid slot id"})
		return
	}

	err = controller.Service.DeleteCalendarSlot(ctx, slotID)
	if err != nil {
		c.JSON(utils.GetStatusCode(err), gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "slot deleted successfully"})
}

func (controller *AdminController) UpdateCalendarSlot(c *gin.Context) {
	var input dtos.CalendarSlotsInput

	ctx, cancel := utils.NewDBContext()
	defer cancel()

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	clientIDValue, exists := c.Get("id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "client id not found in context"})
		return
	}

	clientID, err := uuid.Parse(clientIDValue.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid client id"})
		return
	}

	slotID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid calendar slot id"})
		return
	}

	slots, err := controller.Service.UpdateCalendarSlot(ctx, slotID, clientID, input)
	if err != nil {
		c.JSON(utils.GetStatusCode(err), gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, slots)
}

// Retornando para dar continuidade ao projeto

func (controller *AdminController) GetAvaliableSlots(c *gin.Context) {
	adminIDStr := c.Param("id")

	ctx, cancel := utils.NewDBContext()
	defer cancel()

	adminID, err := uuid.Parse(adminIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid client ID"})
		return
	}

	date := c.Query("date")
	if date == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "date is required"})
		return
	}

	availableSlots, err := controller.Service.GetAvaliableSlots(ctx, adminID, date)
	if err != nil {
		c.JSON(utils.GetStatusCode(err), gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, availableSlots)
}

func (controller *AdminController) CancelAppointmentByAdmin(c *gin.Context) {
	ctx, cancel := utils.NewDBContext()
	defer cancel()

	adminIDValue, exists := c.Get("id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "client id not found in context"})
		return
	}

	adminID, err := uuid.Parse(adminIDValue.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid client id"})
		return
	}

	appointmentIDStr := c.Param("id")

	appointmentID, err := uuid.Parse(appointmentIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid appointment id"})
		return
	}
	err = controller.Service.CancelAppointmentByAdmin(ctx, appointmentID, adminID)
	if err != nil {
		c.JSON(utils.GetStatusCode(err), gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "appointment cancelled successfully",
	})
}

func (controller *AdminController) UpdateAppointment(c *gin.Context) {
	ctx, cancel := utils.NewDBContext()
	defer cancel()

	adminIDValue, exists := c.Get("id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "admin id not found in context"})
		return
	}

	adminIDStr, ok := adminIDValue.(string)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid admin id in context"})
		return
	}

	adminID, err := uuid.Parse(adminIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid admin id"})
		return
	}

	appointmentID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid appointment id"})
		return
	}

	var input dtos.AppointmentUpdateInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	err = controller.Service.UpdateAppointment(ctx, appointmentID, adminID, input)
	if err != nil {
		c.JSON(utils.GetStatusCode(err), gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "appointment updated successfully",
	})
}

func (controller *AdminController) GetAdminProfile(c *gin.Context) {
	adminIDStr, exists := c.Get("id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "client id not found in context"})
		return
	}

	adminID, err := uuid.Parse(adminIDStr.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid client id"})
		return
	}

	ctx, cancel := utils.NewDBContext()
	defer cancel()

	profile, err := controller.Service.GetAdminProfile(ctx, adminID)
	if err != nil {
		c.JSON(utils.GetStatusCode(err), gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, profile)
}