package services

import (
	"context"
	"fmt"
	"tracker/app/models"

	"github.com/lib/pq"
)

/**
 Admin types:
 For Placment info regarding substitutes and logistics
 View: locations, materail, district, programs, students, subjects, tutors
 Delete: Tutor-locations
**/

func (s *AuthService) AddSchedule(c context.Context, req models.RegisterSchedule) (*models.RegisterScheduleResponse, error) {
	// Input validation
	if req.TutorID == nil || req.ScheduleType == "" || req.StartDate.IsZero() {
		return nil, fmt.Errorf("missing required fields: first_name, last_name, or email")
	}
	if req.EndDate != nil && req.EndDate.Before(req.StartDate) {
		return nil, fmt.Errorf("end_date must be after start_date")
	}
	var newID int64
	query, args := buildScheduleQuery(&req)
	err := s.db.QueryRowContext(c, query, args...).Scan(&newID)
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
		args = append(args, req.LocationID)
		query += `INSERT INTO stu_tracker.tutor_schedules 
		(tutor_id, program_id, schedule_type, start_date, end_date, notes, location_id) 
		VALUES ($1,$2,$3,$4,$5,$6,$7) RETURNING id;`
	}
	if req.ScheduleType == "inclusion" {
		args = append(args, req.TutorID)
		args = append(args, req.ProgramID)
		args = append(args, req.ScheduleType)
		args = append(args, req.StartDate)
		args = append(args, req.EndDate)
		args = append(args, req.Notes)
		args = append(args, pq.Array(req.WorkWeek))
		args = append(args, req.StartTime)
		args = append(args, req.EndTime)
		args = append(args, req.LocationID)
		query += `INSERT INTO stu_tracker.tutor_schedules 
		(tutor_id, program_id, schedule_type, start_date, end_date, notes, workweek, start_time, end_time, location_id) 
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10) RETURNING id;`
	}
	return query, args
}

func (s *AuthService) DeleteSchedule(c context.Context, req models.RemoveSchedule) (*models.RegisterScheduleResponse, error) {
	if req.ID == nil {
		return nil, fmt.Errorf("missing delete schedule")
	}
	query := `DELETE FROM stu_tracker.tutor_schedules WHERE id = $1;`
	_, err := s.db.ExecContext(c, query, req.ID)
	if err != nil {
		return nil, fmt.Errorf("unable to delete staff: %w", err)

	}
	return &models.RegisterScheduleResponse{
		ID:     *req.ID,
		Status: "Removed",
	}, nil
}

func (s *AuthService) AddScheduleLink(c context.Context, req models.RegisterScheduleLink) (*models.SimpleReponse, error) {
	var newID int64
	query := `INSERT INTO stu_tracker.Tutor_Schedule_Assignment (tutor_id, schedule_rule_id, location_id) VALUES ($1, $2, $3) RETURNING id;`
	err := s.db.QueryRowContext(c, query, req.TutorID, req.ScheduleID, req.LocationID).Scan(&newID)
	if err != nil {
		return nil, fmt.Errorf("failed to insert student: %w", err)
	}
	return &models.SimpleReponse{
		Status: "OK",
		ID:     &newID,
	}, nil
}

func (s *AuthService) DeleteScheduleLink(c context.Context, req models.RemoveRequest) (*models.RegisterScheduleResponse, error) {
	if req.ID == nil {
		return nil, fmt.Errorf("missing delete schedule")
	}
	query := `DELETE FROM stu_tracker.Tutor_Schedule_Assignment WHERE id = $1;`
	_, err := s.db.ExecContext(c, query, req.ID)
	if err != nil {
		return nil, fmt.Errorf("unable to delete tutor schedule assignment: %w", err)

	}
	return &models.RegisterScheduleResponse{
		ID:     *req.ID,
		Status: "Removed",
	}, nil
}
