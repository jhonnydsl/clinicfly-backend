package services

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jhonnydsl/clinify-backend/src/dtos"
	"github.com/jhonnydsl/clinify-backend/src/mailer"
	"github.com/jhonnydsl/clinify-backend/src/repository"
	"github.com/jhonnydsl/clinify-backend/src/utils"
	"github.com/patrickmn/go-cache"
)

type AdminService struct {
	Repo *repository.AdminRepository
	Mailer *mailer.Mailer
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

func (service *AdminService) CreateAppointment(ctx context.Context, input dtos.AppointmentInput, clientID uuid.UUID) (uuid.UUID, error) {
	parsedDate, err := utils.ParseDate(input.Date)
	if err != nil {
		utils.LogError("createAppointment service (error to parse date)", err)
		return uuid.UUID{}, utils.BadRequestError("invalid format date")
	}

	start, err := utils.ParseTime(input.StartTime)
	if err != nil {
		utils.LogError("createAppoibtment service (error to parse star_time)", err)
		return uuid.UUID{}, utils.BadRequestError("invalid format start_time")
	}

	end, err := utils.ParseTime(input.EndTime)
	if err != nil {
		utils.LogError("createAppointment service (error to parse end_time)", err)
		return uuid.UUID{}, utils.BadRequestError("invalid format end_time")
	}

	if !start.Before(end) {
		return uuid.UUID{}, utils.BadRequestError("start_time must be before end_time")
	}

	id, err := service.Repo.CreateAppointment(ctx, input, parsedDate, start, end, clientID)
	if err != nil {
		utils.LogError("createAppointment service (error call to createAppointment repository)", err)
		return uuid.UUID{}, utils.InternalServerError("error creating appointment")
	}

	patientUUID, err := uuid.Parse(input.PatientID)
	if err != nil {
		return uuid.UUID{}, utils.BadRequestError("invalid patient id format")
	}

	email, err := service.Repo.GetPatientsByEmail(ctx, patientUUID)
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

func (service *AdminService) GetAppointments(ctx context.Context, adminID uuid.UUID, page, limit int) ([]dtos.AppointmentOutput, int, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	cacheKey := fmt.Sprintf("appointments_page_%d_limit_%d", page, limit)

	if cached, found := utils.Cache.Get(cacheKey); found {
		cachedRes := cached.(*utils.AppointmentsCache)
		return cachedRes.Data, cachedRes.Total, nil
	}

	appointments, total, err := service.Repo.GetAllAppointments(ctx, adminID, page, limit)
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