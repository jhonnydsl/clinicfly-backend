package services

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jhonnydsl/clinify-backend/src/dtos"
	"github.com/jhonnydsl/clinify-backend/src/repository"
	"github.com/jhonnydsl/clinify-backend/src/utils"
	"github.com/patrickmn/go-cache"
)

type AdminService struct {
	Repo *repository.AdminRepository
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

	return id, nil
}

func (service *AdminService) GetPatients(ctx context.Context, page, limit int) ([]dtos.PatientOutput, int, error) {
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

	patients, total, err := service.Repo.GetPatients(ctx, page, limit)
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