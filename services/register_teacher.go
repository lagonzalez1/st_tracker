package services

import (
	"context"
	"fmt"
	"tracker/app/models"
)

func (s *AuthService) AddTeacher(ctx context.Context, req models.RegisterTeacher) (*models.RegisterTeacherResponse, error) {
	// Input validation
	if req.Name == "" || req.LocationID == nil {
		return nil, fmt.Errorf("missing location and or name")
	}
	var newID int64
	query := `INSERT INTO stu_tracker.Locations_teachers(name, room, grade_level, substitute, location_id) 
			  VALUES ($1, $2, $3, $4, $5) RETURNING id;`
	err := s.db.QueryRowContext(ctx, query, req.Name, req.Room, req.GradeLevel, req.Substitute, req.LocationID).Scan(&newID)
	if err != nil {
		return nil, fmt.Errorf("failed to insert Teacher: %w", err)
	}
	return &models.RegisterTeacherResponse{
		Status: "OK",
		ID:     &newID,
	}, nil
}
func (s *AuthService) UpdateTeacher(ctx context.Context, req models.RegisterTeacher) (*models.RegisterTeacherResponse, error) {
	// Input validation
	if req.ID == nil {
		return nil, fmt.Errorf("missing teacher ID for update")
	}
	if req.Name == "" || req.LocationID == nil {
		return nil, fmt.Errorf("missing required fields: name or location ID")
	}
	query := `
		UPDATE stu_tracker.Locations_teachers
		SET name = $1,
		    room = $2,
		    grade_level = $3,
		    substitute = $4,
		    location_id = $5,
		    created_at = NOW()
		WHERE id = $6
		RETURNING id;
	`
	var updatedID int64
	err := s.db.QueryRowContext(ctx, query,
		req.Name,
		req.Room,
		req.GradeLevel,
		req.Substitute,
		*req.LocationID,
		*req.ID,
	).Scan(&updatedID)

	if err != nil {
		return nil, fmt.Errorf("failed to update teacher: %w", err)
	}

	return &models.RegisterTeacherResponse{
		Status: "Updated",
		ID:     &updatedID,
	}, nil
}

func (s *AuthService) DeleteTeacher(ctx context.Context, req models.RegisterTeacher) (*models.RegisterTeacherResponse, error) {
	// Input validation
	if req.ID == nil {
		return nil, fmt.Errorf("missing teacher ID for delete")
	}
	query := `DELETE FROM stu_tracker.Locations_teachers WHERE id = $1;`
	_, err := s.db.ExecContext(ctx, query, req.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to delete student: %w", err)
	}
	return &models.RegisterTeacherResponse{
		Status: "Deleted",
		ID:     nil,
	}, nil
}
