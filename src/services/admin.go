package services

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	auditservices "github.com/jhonnydsl/clinify-backend/src/audit/services"
	"github.com/jhonnydsl/clinify-backend/src/dtos"
	"github.com/jhonnydsl/clinify-backend/src/utils"
	"github.com/patrickmn/go-cache"
)

type Mailer interface {
	Send(to, object, body string) error
}

type AdminService struct {
	Repo AdminRepository
	Mailer Mailer
	AuditService *auditservices.AuditService
}

func (service *AdminService) AuditAppointmentCreation(ctx context.Context, actorID, appointmentID uuid.UUID, actorRole, ipAddress, userAgent string) {
	err := service.AuditService.CreateAuditLog(ctx, dtos.AuditLogInput{
        ActorID:      &actorID,
        ActorRole:    &actorRole,
        Action:       "appointment.created",
        ResourceType: "appointment",
        ResourceID:   &appointmentID,
        IPAddress:    ipAddress,
        UserAgent:    userAgent,
    })

    if err != nil {
        utils.LogError("createAppointment service (audit log error)", err)
    }
}

func (services *AdminService) CreateAdmin(ctx context.Context, admin dtos.AdminInput) (uuid.UUID, error) {
	if err := utils.ValidateAdminInput(admin); err != nil {
		utils.LogError("CreatingAdmin service (error validating admin input)", err)
		return uuid.UUID{}, utils.BadRequestError(err.Error())
	}

	hashedPassword, err := utils.HashPassword(admin.Password)
	if err != nil {
		utils.LogError("HashPassword (error to hash password)", err)
		return uuid.UUID{}, utils.InternalServerError("error creating user admin")
	}
	
	admin.Password = hashedPassword

	parsedDate, err := time.Parse("2006-01-02", admin.BirthDate)
	if err != nil {
		utils.LogError("createAdmin service (error parsing birth date)", err)
		return uuid.UUID{}, utils.BadRequestError("invalid birth date format, expected YYYY-MM-DD")
	}

	id, err := services.Repo.CreateAdmin(ctx, admin, parsedDate)
	if err != nil {
		utils.LogError("CreateAdmin service (error to call createAdmin repository)", err)
		return uuid.UUID{}, utils.InternalServerError("error create user admin")
	}

	return id, nil
}

func (service *AdminService) CreateAppointment(ctx context.Context, input dtos.AppointmentInput, clientID, actorID uuid.UUID, actorRole, ipAddress, userAgent string) (uuid.UUID, error) {
	parsedDate, err := utils.ParseDate(input.Date)
	if err != nil {
		utils.LogError("createAppointment service (error to parse date)", err)
		return uuid.UUID{}, utils.BadRequestError("invalid format date")
	}

	startTime, err := utils.ParseTime(input.StartTime)
	if err != nil {
		utils.LogError("createAppoibtment service (error to parse star_time)", err)
		return uuid.UUID{}, utils.BadRequestError("invalid format start_time")
	}

	endTime, err := utils.ParseTime(input.EndTime)
	if err != nil {
		utils.LogError("createAppointment service (error to parse end_time)", err)
		return uuid.UUID{}, utils.BadRequestError("invalid format end_time")
	}

	if !startTime.Before(endTime) {
		return uuid.UUID{}, utils.BadRequestError("start_time must be before end_time")
	}

	start, err := utils.ParseDateTime(input.Date, input.StartTime)
	if err != nil {
		return uuid.UUID{}, utils.BadRequestError("invalid format start_time")
	}

	end, err := utils.ParseDateTime(input.Date, input.EndTime)
	if err != nil {
		return uuid.UUID{}, utils.BadRequestError("invalid format end_time")
	}

	weekday := int(parsedDate.Weekday())

	slots, err := service.Repo.GetCalendarSlotsByWeekday(ctx, clientID, weekday)
	if err != nil {
		utils.LogError("createAppointment service (error getting calendar slots)", err)
		return uuid.UUID{}, utils.InternalServerError("error getting calendar slots")
	}

	validSlot := false

	startMinutes := start.Hour()*60 + start.Minute()
	endMinutes := end.Hour()*60 + end.Minute()

	for _, slot := range slots {
		slotStartMinutes := slot.StartTime.Hour()*60 + slot.StartTime.Minute()
		slotEndMinutes := slot.EndTime.Hour()*60 + slot.EndTime.Minute()

		if startMinutes >= slotStartMinutes && endMinutes <= slotEndMinutes {
			validSlot = true
			break
		}
	}

	if !validSlot {
		return uuid.UUID{}, utils.BadRequestError("appointment is outside the doctor's availability")
	}

	appointments, err := service.Repo.GetAppointmentsByDate(ctx, clientID, input.Date)
	if err != nil {
		utils.LogError("createAppointment service (error getting appointments)", err)
		return uuid.UUID{}, utils.InternalServerError("error getting appointments")
	}

	for _, appointment := range appointments {
		appointmentStart, err := utils.ParseDateTime(appointment.Date, appointment.StartTime)
		if err != nil {
			utils.LogError("createAppointment service (error parsing appointment start_time)", err)
			return uuid.UUID{}, utils.InternalServerError("error validating appointments")
		}

		appointmentEnd, err := utils.ParseDateTime(appointment.Date, appointment.EndTime)
		if err != nil {
			utils.LogError("createAppointment service (error parsing appointment end_time)", err)
			return uuid.UUID{}, utils.InternalServerError("error validating appointments")
		}

		if start.Before(appointmentEnd) && end.After(appointmentStart) {
			return uuid.UUID{}, utils.BadRequestError("this time slot is already booked")
		}
	}

	id, err := service.Repo.CreateAppointment(ctx, input, parsedDate, start, end, clientID)
	if err != nil {
		utils.LogError("createAppointment service (error call to createAppointment repository)", err)
		return uuid.UUID{}, utils.InternalServerError("error creating appointment")
	}

	service.AuditAppointmentCreation(ctx, actorID, id, actorRole, ipAddress, userAgent)

	patientUUID, err := uuid.Parse(input.PatientID)
	if err != nil {
		return uuid.UUID{}, utils.BadRequestError("invalid patient id format")
	}

	email, err := service.Repo.GetPatientEmailByID(ctx, patientUUID)
	if err != nil {
		utils.LogError("createAppointment service (error call to getPatientsByEmail repository)", err)
		return uuid.UUID{}, utils.InternalServerError("error getting email")
	}

	body := utils.BuildAppointmentEmailBody(input.Date, input.StartTime, input.EndTime)

	go func() {
		if err := service.Mailer.Send(email, "Confirmação de Agendamento", body); err != nil {
			utils.LogError("error sending email", err)
		}
	}()

	return id, nil
}

func (service *AdminService) GetAppointments(ctx context.Context, adminID uuid.UUID, status string, page, limit int) ([]dtos.AppointmentOutput, int, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	switch status {
	case "", "scheduled", "cancelled", "completed":
	default:
		return nil, 0, utils.BadRequestError("invalid appointment status")
	}

	cacheKey := fmt.Sprintf("appointments_admin_%s_status_%s_page_%d_limit_%d", adminID.String(), status, page, limit)

	if cached, found := utils.Cache.Get(cacheKey); found {
		cachedRes := cached.(*utils.AppointmentsCache)
		return cachedRes.Data, cachedRes.Total, nil
	}

	appointments, total, err := service.Repo.GetAllAppointments(ctx, adminID, status, page, limit)
	if err != nil {
		return nil, 0, err
	}

	utils.Cache.Set(cacheKey, &utils.AppointmentsCache {
		Data: appointments,
		Total: total,
	}, cache.DefaultExpiration)

	return appointments, total, nil
}

func (service *AdminService) GetPatients(ctx context.Context, adminID uuid.UUID, page, limit int) ([]dtos.PatientOutput, int, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	cacheKey := fmt.Sprintf("patients_page_%d_limit_%d", page, limit)

	if cached, found := utils.Cache.Get(cacheKey); found {
		cachedRes := cached.(*utils.PatientsCache)
		return cachedRes.Data, cachedRes.Total, nil
	}

	patients, total, err := service.Repo.GetPatients(ctx, adminID, page, limit)
	if err != nil {
		return nil, 0, err
	}

	utils.Cache.Set(cacheKey, &utils.PatientsCache {
		Data: patients,
		Total: total,
	}, cache.DefaultExpiration)

	return patients, total, nil
}

func (service *AdminService) DeletePatient(ctx context.Context, patientID uuid.UUID) error {
	if patientID == uuid.Nil {
		return utils.BadRequestError("invalid patient id")
	}

	return service.Repo.DeletePatient(ctx, patientID)
}

func (service *AdminService) CreateCalendarSlot(ctx context.Context, input dtos.CalendarSlotsInput, adminID uuid.UUID) (uuid.UUID, error) {
	start, err := utils.ParseTime(input.StartTime)
	if err != nil {
		utils.LogError("createCalendarSlot service (error to parse date)", err)
		return uuid.UUID{}, utils.BadRequestError("invalid format date")
	}

	end, err := utils.ParseTime(input.EndTime)
	if err != nil {
		utils.LogError("createCalendarSlot service (error to parse date)", err)
		return uuid.UUID{}, utils.BadRequestError("invalid format date")
	}

	if !end.After(start) {
		utils.LogError("createCalendarSlot service (end time must be after start time", err)
		return uuid.UUID{}, utils.BadRequestError("end time must be after start time")
	}

	// Metodo para evitar sobreposição de slots
	slots, err := service.Repo.GetCalendarSlotsByWeekday(ctx, adminID, input.Weekday)
	if err != nil {
		utils.LogError("createCalendarSlot service (error getting calendar slots)", nil)
		return uuid.UUID{}, utils.InternalServerError("error getting calendar slots")
	}

	startMinutes := start.Hour()*60 + start.Minute()
	endMinutes := end.Hour()*60 + end.Minute()

	for _, slot := range slots {
		slotStart := slot.StartTime.Hour()*60 + slot.StartTime.Minute()
		slotEnd := slot.EndTime.Hour()*60 + slot.EndTime.Minute()

		if startMinutes < slotEnd && endMinutes > slotStart {
			return uuid.UUID{}, utils.BadRequestError("calendar slot overlaps with an existing slot")
		}
	}

	id, err := service.Repo.CreateCalendarSlot(ctx, input, start, end, adminID)
	if err != nil {
		utils.LogError("createCalendarSlot service (error call to repository)", err)
		return uuid.UUID{}, utils.InternalServerError("error creating calendar slot")
	}

	return id, nil
}

func (service *AdminService) GetCalendarSlots(ctx context.Context, adminID uuid.UUID) ([]dtos.CalendarSlotsOutput, error) {
	cacheKey := fmt.Sprintf("calendar_slots_%s", adminID)

	if cached, found := utils.Cache.Get(cacheKey); found {
		return cached.(*utils.SlotsCache).Data, nil
	}

	slotsOutputs, err := service.Repo.GetCalendarSlots(ctx, adminID)
	if err != nil {
		utils.LogError("getCalendarSlots service (error call to repo)", err)
		return nil, utils.InternalServerError("error getting slots")
	}

	utils.Cache.Set(cacheKey, &utils.SlotsCache{
		Data: slotsOutputs,
	}, cache.DefaultExpiration)

	return slotsOutputs, nil
}

func (service *AdminService) DeleteCalendarSlot(ctx context.Context, slotID uuid.UUID) error {
	if slotID == uuid.Nil {
		return utils.BadRequestError("invalid slot id")
	}

	return service.Repo.DeleteCalendarSlot(ctx, slotID)
}

func (service *AdminService) UpdateCalendarSlot(ctx context.Context, slotID, adminID uuid.UUID, input dtos.CalendarSlotsInput) ([]dtos.CalendarSlotsOutput, error) {
	if err := utils.ValidateCalendarSlotInput(input); err != nil {
		return nil, utils.BadRequestError(err.Error())
	}

	start, err := utils.ParseTime(input.StartTime)
	if err != nil {
		utils.LogError("updateCalendarSlot service (error parsing start time)", err)
		return nil, utils.BadRequestError("invalid start time format")
	}

	end, err := utils.ParseTime(input.EndTime)
	if err != nil {
		utils.LogError("updateCalendarSlot service (error parsing end time)", err)
		return nil, utils.BadRequestError("invalid end time format")
	}

	slots, err := service.Repo.GetCalendarSlotsByWeekday(ctx, adminID, input.Weekday)
	if err != nil {
		utils.LogError("updateCalendarSlot service (error getting calendar slots)", err)
		return nil, utils.InternalServerError("error getting calendar slots")
	}

	startMinutes := start.Hour()*60 + start.Minute()
	endMinutes := end.Hour()*60 + end.Minute()

	for _, slot := range slots {
		if slot.ID == slotID {
			continue
		}

		slotStart := slot.StartTime.Hour()*60 + slot.StartTime.Minute()
		slotEnd := slot.EndTime.Hour()*60 + slot.EndTime.Minute()

		if startMinutes < slotEnd && endMinutes > slotStart {
			return nil, utils.BadRequestError("calendar slot overlaps with an existing slot")
		}
	}

	err = service.Repo.UpdateCalendarSlot(ctx, slotID, adminID, input)
	if err != nil {
		return nil, err
	}

	slotsOutput, err := service.Repo.GetCalendarSlots(ctx, adminID)
	if err != nil {
		utils.LogError("updateCalendarSlot service (error getting update slots)", err)
		return nil, err
	}

	return slotsOutput, nil
}

func (service *AdminService) GetAvaliableSlots(ctx context.Context, adminID uuid.UUID, date string) ([]string, error) {
	parsedDate, err := utils.ParseDate(date)
	if err != nil {
		utils.LogError("getAvaliableSlots service (error parsed date)", err)
		return nil, utils.BadRequestError("invalid date format")
	}

	weekday := int(parsedDate.Weekday())

	slots, err := service.Repo.GetCalendarSlotsByWeekday(ctx, adminID, weekday)
	if err != nil {
		utils.LogError("getAvaliableSlots service (error call repository)", err)
		return nil, utils.InternalServerError("error getting calendar slots")
	}

	var possibleSlots []time.Time
	interval := 60 * time.Minute

	for _, slot := range slots {
		current := slot.StartTime

		for current.Add(interval).Equal(slot.EndTime) || current.Add(interval).Before(slot.EndTime) {
			possibleSlots = append(possibleSlots, current)
			current = current.Add(interval)
		}
	}

	appointments, err := service.Repo.GetAppointmentsByDate(ctx, adminID, parsedDate.Format("2006-01-02"))
	if err != nil {
		utils.LogError("getAvaliableSlots service (error call to repository)", err)
		return nil, utils.InternalServerError("error getting appointments")
	}

	occupied := make(map[string]bool)
	for _, appt := range appointments {
		occupied[appt.StartTime] = true
	}

	var available []string
	for _, slot := range possibleSlots {
		formatted := slot.Format("15:04")

		if !occupied[formatted] {
			available = append(available, formatted)
		}
	}

	return available, nil
}

func (service *AdminService) CancelAppointmentByAdmin(ctx context.Context, appointmentID, adminID uuid.UUID) error {
	appointment, err := service.Repo.GetAllAppointmentByID(ctx, appointmentID, adminID)
	if err != nil {
		utils.LogError("cancelAppointmentByAdmin service (error getting appointment)", err)
		return err
	}

	err = service.Repo.CancelAppointmentByAdmin(ctx, appointmentID, adminID)
	if err != nil {
		utils.LogError("cancelAppointmentByAdmin service (error cancelling appointment)", err)
		return err
	}

	body := utils.BuildAppointmentCancellationEmailBody(
		appointment.Date,
		appointment.StartTime,
		appointment.EndTime,
	)

	go func() {
		if err := service.Mailer.Send(
			appointment.PatientEmail,
			"Cancelamento de Agendamento",
			body,
		); err != nil {
			utils.LogError("error sending patient cancellation email", err)
		}

		if err := service.Mailer.Send(
			appointment.AdminEmail,
			"Cancelamento de Agendamento",
			body,
		); err != nil {
			utils.LogError("error sending admin cancellation email", err)
		}
	}()

	return nil
}

func (service *AdminService) UpdateAppointment(ctx context.Context, appointmentID, adminID uuid.UUID, input dtos.AppointmentUpdateInput) error {
	parsedDate, err := utils.ParseDate(input.Date)
	if err != nil {
		utils.LogError("updateAppointment service (error parsing date)", err)
		return utils.BadRequestError("invalid date format")
	}

	startTime, err := utils.ParseTime(input.StartTime)
	if err != nil {
		utils.LogError("updateAppointment service (error parsing start time)", err)
		return utils.BadRequestError("invalid start time format")
	}

	endTime, err := utils.ParseTime(input.EndTime)
	if err != nil {
		utils.LogError("updateAppointment service (error parsing end time)", err)
		return utils.BadRequestError("invalid end time format")
	}

	if !startTime.Before(endTime) {
		return utils.BadRequestError("start time must be before end time")
	}

	appointment, err := service.Repo.GetAllAppointmentByID(ctx, appointmentID, adminID)
	if err != nil {
		utils.LogError("updateAppointment service (error getting appointment)", err)
		return err
	}

	weekday := int(parsedDate.Weekday())

	slots, err := service.Repo.GetCalendarSlotsByWeekday(ctx, adminID, weekday)
	if err != nil {
		utils.LogError("updateAppointment service (error getting calendar slots)", err)
		return err
	}

	withinSchedule := false

	for _, slot := range slots {
		if !startTime.Before(slot.StartTime) && !endTime.After(slot.EndTime) {
			withinSchedule = true
			break
		}
	}

	if !withinSchedule {
		return utils.BadRequestError("appointment time is outside doctor's availability")
	}

	appointments, err := service.Repo.GetAppointmentsByDate(ctx, adminID, parsedDate.Format("2006-01-02"))
	if err != nil {
		utils.LogError("updateAppointment service (error getting appointments)", err)
		return err
	}

	for _, existingAppointment := range appointments {
		if existingAppointment.ID == appointmentID {
			continue
		}

		existingStart, err := utils.ParseTime(existingAppointment.StartTime)
		if err != nil {
			utils.LogError("updateAppointment service (error parsing existing start time)", err)
			return utils.InternalServerError("error validating appointment")
		}

		existingEnd, err := utils.ParseTime(existingAppointment.EndTime)
		if err != nil {
			utils.LogError("updateAppointment service (error parsing existing end time)", err)
			return utils.InternalServerError("error validating appointment")
		}

		if startTime.Before(existingEnd) && endTime.After(existingStart) {
			return utils.BadRequestError("appointment time is already occupied")
		}
	}

	err = service.Repo.UpdateAppointment(ctx, appointmentID, adminID, parsedDate, startTime, endTime)
	if err != nil {
		utils.LogError("updateAppointment service (error updating appointment)", err)
		return err
	}

	body := utils.BuildAppointmentUpdateEmailBody(
		appointment.Date,
		appointment.StartTime,
		appointment.EndTime,
		parsedDate.Format("2006-01-02"),
		startTime.Format("15:04"),
		endTime.Format("15:04"),
	)

	go func() {
		if err := service.Mailer.Send(
			appointment.PatientEmail,
			"Alteração de Agendamento",
			body,
		); err != nil {
			utils.LogError("error sending patient appointment update email", err)
		}

		if err := service.Mailer.Send(
			appointment.AdminEmail,
			"Alteração de Agendamento",
			body,
		); err != nil {
			utils.LogError("error sending admin appointment update email", err)
		}
	}()

	return nil
}

func (service *AdminService) GetAdminProfile(ctx context.Context, adminID uuid.UUID) (dtos.AdminProfileOutput, error) {
	profile, err := service.Repo.GetAdminProfile(ctx, adminID)
	if err != nil {
		utils.LogError("getAdminProfile service (error getting profile)", err)
		return dtos.AdminProfileOutput{}, err
	}

	return profile, nil
}

func (service *AdminService) UpdateAdminProfile(ctx context.Context, adminID uuid.UUID, input dtos.AdminProfileInput) (dtos.AdminProfileOutput, error) {
	if err := utils.ValidateAdminProfileInput(input); err != nil {
		return dtos.AdminProfileOutput{}, utils.BadRequestError(err.Error())
	}

	if input.Email != nil {
		exists, err := service.Repo.EmailExists(ctx, *input.Email, adminID)
		if err != nil {
			return dtos.AdminProfileOutput{}, err
		}

		if exists {
			return dtos.AdminProfileOutput{}, utils.ConflictError("invalid profile data")
		}
	}

	if input.PublicSlug != nil {
		exists, err := service.Repo.PublicSlugExists(ctx, *input.PublicSlug, adminID)
		if err != nil {
			return dtos.AdminProfileOutput{}, err
		}

		if exists {
			return dtos.AdminProfileOutput{}, utils.ConflictError("public slug already in use")
		}
	}

	err := service.Repo.UpdateAdminProfile(ctx, adminID, input)
	if err != nil {
		return dtos.AdminProfileOutput{}, err
	}

	profile, err := service.Repo.GetAdminProfile(ctx, adminID)
	if err != nil {
		return dtos.AdminProfileOutput{}, err
	}

	return profile, nil
}