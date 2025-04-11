package services

import (
	"fmt"
	"tracker/app/models"
)

/**
 Admin types:
 For Placment info regarding substitutes and logistics
 View: locations, materail, district, programs, students, subjects, tutors
 Delete: Tutor-locations

**/

func (s *AuthService) AddSchedule(req models.RegisterSchedule) (*models.RegisterScheduleResponse, error) {
	// Input validation
	if req.TutorID == nil || req.ScheduleType == "" || req.StartDate.IsZero() {
		return nil, fmt.Errorf("missing required fields: first_name, last_name, or email")
	}
	if req.EndDate != nil && req.EndDate.Before(req.StartDate) {
		return nil, fmt.Errorf("end_date must be after start_date")
	}
	var newID int64
	query, args := buildScheduleQuery(&req)
	fmt.Println(query)
	err := s.db.QueryRow(query, args...).Scan(&newID)
	if err != nil {
		return nil, fmt.Errorf("failed to insert student: %w", err)
	}
	return &models.RegisterScheduleResponse{
		Status: "OK",
		ID:     newID,
	}, nil
}

func buildScheduleQuery(req *models.RegisterSchedule) (string, []interface{}) {
	var args []interface{}
	query := ``
	if req.ScheduleType == "exclusion" {
		args = append(args, req.TutorID)
		args = append(args, req.ProgramID)
		args = append(args, req.ScheduleType)
		args = append(args, req.StartDate)
		args = append(args, req.StartDate)
		args = append(args, req.Notes)
		query += `INSERT INTO stu_tracker.tutor_schedules 
		(tutor_id, program_id, schedule_type, start_date, end_date, notes) 
		VALUES ($1,$2,$3,$4,$5,$6) RETURNING id;`
	}
	if req.ScheduleType == "inclusion" {
		args = append(args, req.TutorID)
		args = append(args, req.ProgramID)
		args = append(args, req.ScheduleType)
		args = append(args, req.StartDate)
		args = append(args, req.EndDate)
		args = append(args, req.Notes)
		query += `INSERT INTO stu_tracker.tutor_schedules 
		(tutor_id, program_id, schedule_type, start_date, end_date, notes) 
		VALUES ($1,$2,$3,$4,$5,$6) RETURNING id;`
	}
	return query, args
}

func (s *AuthService) DeleteSchedule(req models.RemoveSchedule) (*models.RegisterScheduleResponse, error) {
	if req.ID == nil {
		return nil, fmt.Errorf("missing delete schedule")
	}
	query := `DELETE FROM stu_tracker.tutor_schedules WHERE id = $1;`
	_, err := s.db.Exec(query, req.ID)
	if err != nil {
		return nil, fmt.Errorf("unable to delete staff: %w", err)

	}
	return &models.RegisterScheduleResponse{
		ID:     *req.ID,
		Status: "Removed",
	}, nil
}
