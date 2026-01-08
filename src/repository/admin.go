package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jhonnydsl/clinify-backend/src/dtos"
	"github.com/jhonnydsl/clinify-backend/src/utils"
)

type AdminRepository struct{}

func (r *AdminRepository) CreateAdmin(ctx context.Context, admin dtos.AdminInput, birthDate time.Time) (uuid.UUID, error) {
	query := `INSERT INTO clients (full_name, email, password_hash, birth_date, crp, bio, profile_image_url, office_address, phone, public_slug)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	RETURNING id`

	var id uuid.UUID

	err := DB.QueryRowContext(
		ctx,
		query, 
		admin.FullName, 
		admin.Email, 
		admin.Password, 
		birthDate,
		admin.Crp,
		admin.Bio,
		admin.ProfileImage,
		admin.OfficeAddress,
		admin.Phone,
		admin.PublicSlug,
	).Scan(&id)
	if err != nil {
		utils.LogError("createAdmin (INSERT clients)", err)
		return uuid.UUID{}, utils.InternalServerError("error creating user admin")
	}

	return id, nil
}

func (r *AdminRepository) FindAdminIDBySlug(ctx context.Context, slug string) (uuid.UUID, error) {
	query := `SELECT id FROM clients WHERE public_slug = $1 LIMIT 1`

	var id uuid.UUID

	err := DB.QueryRowContext(ctx, query, slug).Scan(&id)
	if err != nil {
		utils.LogError("FindAdminBySlug (error in SELECT clients)", err)
		return uuid.UUID{}, utils.InternalServerError("invalid client slug")
	}

	return id, nil
}

func (r *AdminRepository) CreateAppointment(ctx context.Context, input dtos.AppointmentInput, parsedDate, start, end time.Time, clientID uuid.UUID) (uuid.UUID, error) {
	query := `INSERT INTO appointments (client_id, patient_id, date, start_time, end_time, status)
	VALUES ($1, $2, $3, $4, $5, 'scheduled')
	RETURNING id;`

	var id uuid.UUID

	err := DB.QueryRowContext(
		ctx,
		query,
		clientID,
		input.PatientID,
		parsedDate,
		start,
		end,
	).Scan(&id)
	if err != nil {
		utils.LogError("createAppointment repository (error in INSERT)", err)
		return uuid.UUID{}, utils.InternalServerError("error creating appointment")
	}

	return id, nil
}

func (r *AdminRepository) GetAllAppointments(ctx context.Context, adminID uuid.UUID, page, limit int) ([]dtos.AppointmentOutput, int, error) {
	query := `SELECT a.id, a.patient_id, p.full_name, a.date, a.start_time, a.end_time, a.status
	FROM appointments a
	JOIN patients p ON p.id = a.patient_id
	WHERE a.client_id = $1
	ORDER BY p.full_name LIMIT $2 OFFSET $3;`

	queryCount := `SELECT COUNT(*) FROM appointments WHERE client_id = $1`

	offset := (page - 1) * limit

	var total int

	err := DB.QueryRowContext(ctx, queryCount, adminID).Scan(&total)
	if err != nil {
		return nil, 0, utils.InternalServerError("error getting total appointments")
	}

	rows, err := DB.QueryContext(ctx, query, adminID, limit, offset)
	if err != nil {
		utils.LogError("getAppointments repository (error in SELECT)", err)
		return nil, 0, utils.InternalServerError("error getting appointments")
	}
	defer rows.Close()

	var appointments []dtos.AppointmentOutput

	for rows.Next() {
		var (
			id uuid.UUID
			patientID uuid.UUID
			fullName string
			date time.Time
			startTime time.Time
			endTime time.Time
			status string
		)

		err := rows.Scan(&id, &patientID, &fullName, &date, &startTime, &endTime, &status)
		if err != nil {
			utils.LogError("getAppointments repository (scan error)", err)
			return nil, 0, utils.InternalServerError("error fetching appointments")
		}

		appointments = append(appointments, dtos.AppointmentOutput{
			ID: id,
			PatientID: patientID,
			FullName: fullName,
			Date: date.Format("2006-01-02"),
			StartTime: startTime.Format("15:04"),
			EndTime: endTime.Format("15:04"),
			Status: status,
		})
	}

	return appointments, total, nil
}

func (r *AdminRepository) GetAppointmentsByDate(ctx context.Context, adminID uuid.UUID, date string) ([]dtos.AppointmentOutput, error) {
	query := `SELECT a.id, a.patient_id, p.full_name, a.date, a.start_time, a.end_time, a.status
	FROM appointments a
	JOIN patients p ON p.id = a.patient_id
	WHERE a.client_id = $1 AND a.date = $2 AND a.status != 'cancelled'
	ORDER BY a.start_time`

	rows, err := DB.QueryContext(ctx, query, adminID, date)
	if err != nil {
		utils.LogError("getAppointmentsByDate repository (select error)", err)
		return nil, utils.InternalServerError("error getting appointments by date")
	}
	defer rows.Close()

	appointments := make([]dtos.AppointmentOutput, 0)

	for rows.Next() {
		var (
			id uuid.UUID
			patientID uuid.UUID
			fullName string
			dateDB time.Time
			startTime time.Time
			endTime time.Time
			status string
		)

		err := rows.Scan(
			&id,
			&patientID,
			&fullName,
			&dateDB,
			&startTime,
			&endTime,
			&status,
		)
		if err != nil {
			utils.LogError("getAppointmentsByDate repository (scan error)", err)
			return nil, utils.InternalServerError("error scanning appointments")
		}

		appointments = append(appointments, dtos.AppointmentOutput{
			ID: id,
			PatientID: patientID,
			FullName: fullName,
			Date: dateDB.Format("2006-01-02"),
			StartTime: startTime.Format("15:04"),
			EndTime: endTime.Format("15:04"),
			Status: status,
		})
	}

	if err := rows.Err(); err != nil {
		utils.LogError("getAppointmentsByDate repository (rows error)", err)
		return nil, utils.InternalServerError("error iterating appointments")
	}

	return appointments, nil
}

func (r *AdminRepository) GetPatients(ctx context.Context, adminID uuid.UUID, page, limit int) ([]dtos.PatientOutput, int, error) {
	query := `SELECT id, full_name, email, phone, birth_date FROM patients
	WHERE client_id = $1
	ORDER BY full_name LIMIT $2 OFFSET $3`

	queryCount := `SELECT COUNT(*) FROM patients WHERE client_id = $1 `

	offset := (page - 1) * limit

	var total int
	
	err := DB.QueryRowContext(ctx, queryCount, adminID).Scan(&total)
	if err != nil {
		return nil, 0, utils.InternalServerError("error getting total patients")
	}
	
	rows, err := DB.QueryContext(ctx, query, adminID, limit, offset)
	if err != nil {
		utils.LogError("GetPatients repository (error in SELECT)", err)
		return nil, 0, utils.InternalServerError("error getting patients")
	}
	defer rows.Close()
	
	var patients []dtos.PatientOutput

	for rows.Next() {
		var (
			id uuid.UUID
			fullName string
			email string
			phone string
			birthDate time.Time
		)

		err = rows.Scan(&id, &fullName, &email, &phone, &birthDate)
		if err != nil {
			utils.LogError("getPatients repository (scan error)", err)
			return nil, 0, utils.InternalServerError("error fetching patients")
		}

		patients = append(patients, dtos.PatientOutput{
			ID: id,
			FullName: fullName,
			Email: email,
			Phone: phone,
			BirthDate: birthDate.Format("2006-01-02"),
		})
	}
	
	return patients, total, nil
}

func (r *AdminRepository) DeletePatient(ctx context.Context, patientID uuid.UUID) error {
	query := `DELETE FROM patients WHERE id = $1`

	res, err := DB.ExecContext(ctx, query, patientID)
	if err != nil {
		utils.LogError("deletePatient repository (error deleting patient)", err)
		return utils.InternalServerError("error deleting patient")
	}

	rows, err := res.RowsAffected()
	if err != nil {
		utils.LogError("deletePatient repository (error reading rows affected)", err)
		return utils.InternalServerError("error deleting patient")
	}

	if rows == 0 {
		return utils.NotFoundError("patient not found")
	}

	return nil
}

func (r *AdminRepository) GetPatientEmailByID(ctx context.Context, patientID uuid.UUID) (string, error) {
	query := `SELECT email FROM patients WHERE id = $1`

	var email string

	err := DB.QueryRowContext(ctx, query, patientID).Scan(&email)
	if err != nil {
		utils.LogError("getPatientsByEmail repository (error SELECT)", err)
		return "", utils.InternalServerError("error getting email")
	}

	return email, nil
}

func (r *AdminRepository) CreateCalendarSlot(ctx context.Context, input dtos.CalendarSlotsInput, start, end time.Time, adminID uuid.UUID) (uuid.UUID, error) {
	query := `INSERT INTO calendar_slots (client_id, weekday, start_time, end_time)
	VALUES ($1, $2, $3, $4)
	RETURNING id`

	var id uuid.UUID

	err := DB.QueryRowContext(ctx, query, adminID, input.Weekday, start, end).Scan(&id)
	if err != nil {
		utils.LogError("creatCalendarSlot repository (error in INSERT)", err)
		return uuid.UUID{}, utils.InternalServerError("error creating calendar slot")
	}

	return id, nil
}

func (r *AdminRepository) GetCalendarSlots(ctx context.Context, adminID uuid.UUID) ([]dtos.CalendarSlotsOutput, error) {
	query := `SELECT id, weekday, start_time, end_time FROM calendar_slots WHERE client_id = $1`

	var slotsOutput []dtos.CalendarSlotsOutput

	rows, err := DB.QueryContext(ctx, query, adminID)
	if err != nil {
		utils.LogError("getCalendarSlots repository (SELECT error)", err)
		return nil, utils.InternalServerError("error getting slots")
	}
	defer rows.Close()

	for rows.Next() {
		var (
			id uuid.UUID
			weekday string
			startTime time.Time
			endTime time.Time
		)

		err := rows.Scan(&id, &weekday, &startTime, &endTime)
		if err != nil {
			utils.LogError("GetCalendarSlots repository (scan error)", err)
			return nil, utils.InternalServerError("error fetching slots")
		}

		slotsOutput = append(slotsOutput, dtos.CalendarSlotsOutput{
			ID: id,
			Weekday: weekday,
			StartTime: startTime.Format("15:04"),
			EndTime: endTime.Format("15:04"),
		})
	}

	return slotsOutput, nil
}

func (r *AdminRepository) DeleteCalendarSlot(ctx context.Context, slotID uuid.UUID) error {
	query := `DELETE FROM calendar_slots WHERE id = $1`

	res, err := DB.ExecContext(ctx, query, slotID)
	if err != nil {
		utils.LogError("deleteCalendarSlot repository (error deleting slot)", err)
		return utils.InternalServerError("error deleting slot")
	}

	rows, err := res.RowsAffected()
	if err != nil {
		utils.LogError("deleteCalendarSlot repository (error reading rows affected)", err)
		return utils.InternalServerError("error deleting slot")
	}

	if rows == 0 {
		return utils.NotFoundError("slot not found")
	}

	return nil
}

func (r *AdminRepository) GetCalendarSlotsByWeekday(ctx context.Context, adminID uuid.UUID, weekday int) ([]dtos.CalendarSlotDB, error) {
	query := `SELECT id, client_id, weekday, start_time, end_time FROM calendar_slots
	WHERE client_id = $1 AND weekday = $2 ORDER BY start_time`

	rows, err := DB.QueryContext(ctx, query, adminID, weekday)
	if err != nil {
		utils.LogError("getCalendarSlotsByWeekday repository (error getting slots)", err)
		return nil, utils.InternalServerError("error getting slots")
	}
	defer rows.Close()

	slots := make([]dtos.CalendarSlotDB, 0)

	for rows.Next() {
		var slot dtos.CalendarSlotDB

		err := rows.Scan(
			&slot.ID,
			&slot.AdminID,
			&slot.Weekday,
			&slot.StartTime,
			&slot.EndTime,
		)
		if err != nil {
			utils.LogError("getCalendarSlotsByWeekday repository (scan error)", err)
			return nil, utils.InternalServerError("error fetching slots")
		}

		slots = append(slots, slot)
	}

	if err := rows.Err(); err != nil {
		return nil, utils.InternalServerError("error getting slots")
	}

	return slots, nil
}
