package services

import (
	"context"
	"fmt"
	"tracker/app/models"
)

func (s *AuthService) AddStudent(ctx context.Context, req models.RegisterRequestStudents) (*models.ResponseRequestStudents, error) {
	// Input validation
	println(req.LastName)
	if req.FirstName == "" || req.LastName == "" {
		return nil, fmt.Errorf("missing required fields: first_name, last_name, or email")
	}
	var newID int64
	query := `INSERT INTO stu_tracker.Students(first_name, last_name, middle_name, email, grade_level, active, location_id, period, created_by, direct_partnership, tutor_id, semester_id)
              VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12) RETURNING id;`

	err := s.db.QueryRow(query, req.FirstName, req.LastName, req.MiddleName, req.Email, req.GradeLevel, req.Active, req.LocationId, req.Period, req.CreatedBy, req.DirectPartnership, req.TutorID, req.SemesterID).Scan(&newID)
	if err != nil {
		return nil, fmt.Errorf("failed to insert student: %w", err)
	}
	return &models.ResponseRequestStudents{
		Status:    "OK",
		StudentID: newID,
	}, nil
}

func (s *AuthService) UpdateStudent(ctx context.Context, req models.RegisterRequestStudents) (*models.ResponseUpdate, error) {
	// Input validation
	println(req.LastName)
	if req.ID == nil || req.LastName == "" {
		return nil, fmt.Errorf("missing required fields: first_name, last_name, or email")
	}
	query := `UPDATE stu_tracker.Students SET first_name = $1, last_name = $2, middle_name = $3, email = $4, grade_level = $5, active = $6, location_id = $7, period = $8, direct_partnership = $9, tutor_id = $10, semester_id = $11
              WHERE id = $12`

	_, err := s.db.ExecContext(ctx, query, req.FirstName, req.LastName, req.MiddleName, req.Email, req.GradeLevel, req.Active, req.LocationId, req.Period, req.DirectPartnership, req.TutorID, req.SemesterID, req.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to update student: %w", err)
	}
	return &models.ResponseUpdate{
		Status: "Updated",
	}, nil
}

func (s *AuthService) DeleteStudent(ctx context.Context, req models.RemoveRequest) (*models.RemoveResponse, error) {
	// Input validation
	if req.ID == nil {
		return nil, fmt.Errorf("missing required fields: id")
	}
	query := `DELETE FROM stu_tracker.Students WHERE id = $1`

	_, err := s.db.ExecContext(ctx, query, req.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to delete student: %w", err)
	}
	return &models.RemoveResponse{
		Status: "Removed",
	}, nil
}
