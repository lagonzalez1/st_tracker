package services

import (
	"context"
	"fmt"
	"tracker/app/models"

	"github.com/lib/pq"
)

func (s *AuthService) AddScheduleGlobal(c context.Context, req models.RegisterGlobalSchedule) (*models.SimpleReponse, error) {
	var newID int64
	query := `INSERT INTO stu_tracker.Schedule_rule (
				job_name,
				job_description,
				organization_id,
				tutor_id,
				location_id,
				global_rule,
				cron_job,
				provider_id,
				provider_uid,
				provider_type,
				provider_employee_id,
				provider_employee_uid,
				recurrence_type,
				start_date,
				end_date,
				specific_dates,
				frequency,
				start_time,
				end_time,
				program_id,
				semester_id,
				enabled,
				archive
			)
			VALUES (
				$1,  $2,  $3,  $4,  $5,
				$6,  $7,  $8,  $9,  $10,
				$11, $12, $13, $14, $15,
				$16, $17, $18, $19, $20,
				$21, $22, $23
			)
			RETURNING id;`

	err := s.db.QueryRowContext(c, query,
		req.JobName,
		req.JobDescription,
		req.OrganizationID,
		req.TutorID,
		req.LocationID,
		req.GlobalRule,
		req.CronJob,
		req.ProviderID,
		req.ProviderUID,
		req.ProviderType,
		req.ProviderEmployeeID,
		req.ProviderEmployeeUID,
		req.RecurrenceType,
		req.StartDate,
		req.EndDate,
		pq.Array(req.SpecificDates),
		pq.Array(req.Frequency),
		req.StartTime,
		req.EndTime,
		req.ProgramID,
		req.SemesterID,
		req.Enabled,
		req.Archive,
	).Scan(&newID)
	if err != nil {
		return nil, fmt.Errorf("failed to insert Schedule_rule: %w", err)
	}
	return &models.SimpleReponse{
		Status: "OK",
		ID:     &newID,
	}, nil
}
func (s *AuthService) UpdateScheduleGlobal(c context.Context, req models.RegisterGlobalSchedule) (*models.SimpleReponse, error) {
	// First, check if the schedule exists
	var exists bool
	err := s.db.QueryRowContext(c,
		`SELECT EXISTS(SELECT 1 FROM stu_tracker.Schedule_rule WHERE id = $1 AND organization_id = $2)`,
		req.ID,
		req.OrganizationID,
	).Scan(&exists)

	if err != nil {
		return nil, fmt.Errorf("failed to check schedule existence: %w", err)
	}

	if !exists {
		return nil, fmt.Errorf("schedule not found or access denied")
	}

	query := `UPDATE stu_tracker.Schedule_rule SET
                job_name = COALESCE($1, job_name),
                job_description = COALESCE($2, job_description),
                tutor_id = COALESCE($3, tutor_id),
                location_id = COALESCE($4, location_id),
                cron_job = COALESCE($5, cron_job),
                provider_id = COALESCE($6, provider_id),
                provider_uid = COALESCE($7, provider_uid),
                provider_type = COALESCE($8, provider_type),
                provider_employee_id = COALESCE($9, provider_employee_id),
                provider_employee_uid = COALESCE($10, provider_employee_uid),
                recurrence_type = COALESCE($11, recurrence_type),
                start_date = COALESCE($12, start_date),
                end_date = COALESCE($13, end_date),
                specific_dates = COALESCE($14, specific_dates),
                frequency = COALESCE($15, frequency),
                start_time = COALESCE($16, start_time),
                end_time = COALESCE($17, end_time),
                program_id = COALESCE($18, program_id),
                semester_id = COALESCE($19, semester_id),
                enabled = COALESCE($20, enabled),
                archive = COALESCE($21, archive),
                updated_at = CURRENT_TIMESTAMP
              WHERE id = $22 AND organization_id = $23
              RETURNING id;`

	var updatedID int64
	err = s.db.QueryRowContext(c, query,
		req.JobName,
		req.JobDescription,
		req.TutorID,
		req.LocationID,
		req.CronJob,
		req.ProviderID,
		req.ProviderUID,
		req.ProviderType,
		req.ProviderEmployeeID,
		req.ProviderEmployeeUID,
		req.RecurrenceType,
		req.StartDate,
		req.EndDate,
		pq.Array(req.SpecificDates),
		pq.Array(req.Frequency),
		req.StartTime,
		req.EndTime,
		req.ProgramID,
		req.SemesterID,
		req.Enabled,
		req.Archive,
		req.ID,
		req.OrganizationID,
	).Scan(&updatedID)

	if err != nil {
		return nil, fmt.Errorf("failed to update Schedule_rule: %w", err)
	}

	return &models.SimpleReponse{
		Status: "OK",
		ID:     &updatedID,
	}, nil
}

func (s *AuthService) HardDeleteScheduleGlobal(c context.Context, scheduleID, organizationID *int64) (*models.SimpleReponse, error) {
	query := `DELETE FROM stu_tracker.Schedule_rule 
              WHERE id = $1 AND organization_id = $2
              RETURNING id;`
	var deletedID int64
	err := s.db.QueryRowContext(c, query, scheduleID, organizationID).Scan(&deletedID)
	if err != nil {
		return nil, fmt.Errorf("failed to hard delete Schedule_rule: %w", err)
	}
	return &models.SimpleReponse{
		Status: "DELETED",
		ID:     &deletedID,
	}, nil
}
