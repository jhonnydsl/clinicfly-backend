package services

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jhonnydsl/clinify-backend/src/dtos"
	"github.com/jhonnydsl/clinify-backend/src/mailer"
	"github.com/jhonnydsl/clinify-backend/src/repository"
	"github.com/jhonnydsl/clinify-backend/src/utils"
)

type PatientService struct {
	Repo *repository.PatientRepository
	AdminRepo *repository.AdminRepository
	Mailer *mailer.Mailer
}

func (service *PatientService) CreatePatient(ctx context.Context, patient dtos.PatientInput) (uuid.UUID, error) {
	if err := utils.ValidatePatientInput(patient); err != nil {
		utils.LogError("createPatient service (error validating patient input)", err)
		return uuid.UUID{}, utils.BadRequestError(err.Error())
	}

	clientUUID, err := service.AdminRepo.FindAdminIDBySlug(ctx, patient.PublicSlug)
	if err != nil {
		utils.LogError("createPatient service (error get client_id by public_slug)", err)
		return uuid.UUID{}, utils.BadRequestError("invalid admin url")
	}

	hashedPassword, err := utils.HashPassword(patient.Password)
	if err != nil {
		utils.LogError("hashPassword (error to hash password)", err)
		return uuid.UUID{}, utils.InternalServerError("error creating user patient")
	}

	patient.Password = hashedPassword

	parsedDate, err := time.Parse("2006-01-02", patient.BirthDate)
	if err != nil {
		utils.LogError("createAdmin service (error parsing birth date)", err)
		return uuid.UUID{}, utils.BadRequestError("invalid birth date format, expected YYYY-MM-DD")
	}

	id, err := service.Repo.CreatePatient(ctx, patient, parsedDate, clientUUID)
	if err != nil {
		utils.LogError("createPatient service (error to call createPatient repository)", err)
		return uuid.UUID{}, utils.InternalServerError("error creating user patient")
	}

	return id, nil
}

func (service *PatientService) GetDoctorAvaliableSlots(ctx context.Context, adminID uuid.UUID, date string) ([]string, error){
	parsedDate, err := utils.ParseDate(date)
	if err != nil {
		utils.LogError("getAvailableSlots service (error parsed date)", err)
		return nil, utils.BadRequestError("invalid date format")
	}

	weekday := int(parsedDate.Weekday())

	slots, err := service.AdminRepo.GetCalendarSlotsByWeekday(ctx, adminID, weekday)
	if err != nil {
		utils.LogError("getDoctorAvailableSlots service (error getting calendar slots)", err)
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

	appointments, err := service.AdminRepo.GetAppointmentsByDate(ctx, adminID, parsedDate.Format("2006-01-02"))
	if err != nil {
		utils.LogError("getDoctorAvailableSlots service (error getting appointments)", err)
		return nil, utils.InternalServerError("error getting appointments")
	}

	occupied := make(map[string]bool)

	for _, appt := range appointments {
		occupied[appt.StartTime] = true
	}

	available := []string{}

	for _, slot := range possibleSlots {
		formatted := slot.Format("15:04")

		if !occupied[formatted] {
			available = append(available, formatted)
		}
	}

	return available, nil
}

func (service *PatientService) CancelAppointmentByPatient(ctx context.Context, appointmentID, patientID uuid.UUID) error {
	appointment, err := service.Repo.GetAppointmentByID(ctx, appointmentID, patientID)
	if err != nil {
		utils.LogError("cancelAppointmentByPatient service (error getting appointment)", err)
		return err
	}

	err = service.Repo.CancelAppointmentByPatient(ctx, appointmentID, patientID)
	if err != nil {
		utils.LogError("cancelAppointmentByPatient service (error cancelling appointment)", err)
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
			"Cancelamento do Agendamento",
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