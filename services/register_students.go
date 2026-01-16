package services

import (
	"context"
	"fmt"
	"strings"
	"tracker/app/models"
)

func (s *AuthService) AddStudent(ctx context.Context, req models.RegisterRequestStudents) (*models.ResponseRequestStudents, error) {
	// Input validation
	println(req.LastName)
	if req.FirstName == "" || req.LastName == "" {
		return nil, fmt.Errorf("missing required fields: first_name, last_name, or email")
	}
	var newID int64
	query := `INSERT INTO 
				stu_tracker.Students
				(first_name, last_name, middle_name,
				email, grade_level, active,
				location_id, period, created_by,
				direct_partnership, tutor_id, semester_id,
				teacher_id, timeframe, timeframe_start, timeframe_end, duration_required, race, gender)
              VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19) RETURNING id;`

	err := s.db.QueryRowContext(ctx, query, req.FirstName, req.LastName,
		req.MiddleName, req.Email, req.GradeLevel, req.Active, req.LocationId, req.Period,
		req.CreatedBy, req.DirectPartnership, req.TutorID, req.SemesterID, req.TeacherID, req.Timeframe,
		req.TimeframeStart, req.TimeframeEnd, req.DurationRequired, req.Race, req.Gender).Scan(&newID)
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
	query := `UPDATE stu_tracker.Students SET first_name = $1, last_name = $2, 
			middle_name = $3, email = $4, grade_level = $5, active = $6, location_id = $7,
			period = $8, direct_partnership = $9, tutor_id = $10, semester_id = $11,
			teacher_id = $12, timeframe = $13, timeframe_start = $14, timeframe_end = $15, duration_required = $16,
			race = $17, gender = $18
            WHERE id = $19;`

	_, err := s.db.ExecContext(ctx, query, req.FirstName, req.LastName, req.MiddleName,
		req.Email, req.GradeLevel, req.Active, req.LocationId, req.Period, req.DirectPartnership,
		req.TutorID, req.SemesterID, req.TeacherID, req.Timeframe, req.TimeframeStart,
		req.TimeframeEnd, req.DurationRequired, req.Race, req.Gender, req.ID)
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

func (s *AuthService) AddStudentGroup(ctx context.Context, req models.RegisterStudentGroup) (*models.SimpleReponse, error) {
	var newID int64
	query := `INSERT INTO stu_tracker.Student_groups(
		tutor_id, location_id, semester_id, title, description
	) VALUES ($1, $2, $3, $4, $5) RETURNING id;`
	err := s.db.QueryRowContext(ctx, query, req.TutorID, req.LocationID, req.SemesterID, req.Title, req.Description).Scan(&newID)
	if err != nil {
		return nil, err
	}
	return &models.SimpleReponse{
		ID:     &newID,
		Status: "OK",
	}, nil
}

func (s *AuthService) AddStudentGroupAttendies(ctx context.Context, req models.RegisterStudentGroupAttendies) (*models.SimpleReponse, error) {
	values := make([]string, 0, len(req.Students))
	args := make([]interface{}, 0, len(req.Students))

	for i, student := range req.Students {
		values = append(values, fmt.Sprintf("($%d, $%d)", i*2+1, i*2+2))
		args = append(args, student.ID, req.GroupID)
	}
	fmt.Print(len(values))
	fmt.Print(len(args))

	query := fmt.Sprintf(`INSERT INTO stu_tracker.Student_group_attendees(student_id, student_group_id) VALUES %s ON CONFLICT DO NOTHING;`, strings.Join(values, ", "))
	fmt.Println(query)
	_, err := s.db.ExecContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	return &models.SimpleReponse{
		ID:     nil,
		Status: "OK",
	}, nil
}

func (s *AuthService) UpdateStudentGroup(ctx context.Context, req models.RegisterStudentGroup) (*models.ResponseUpdate, error) {
	// Input validation
	if req.ID == nil {
		return nil, fmt.Errorf("missing required fields: first_name, last_name, or email")
	}
	query := `UPDATE stu_tracker.Student_groups SET tutor_id = $1, location_id = $2, 
			semester_id = $3, title = $4, description = $5
            WHERE id = $6;`
	_, err := s.db.ExecContext(ctx, query, req.TutorID, req.LocationID, req.SemesterID, req.Title, req.Description, req.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to update student: %w", err)
	}
	return &models.ResponseUpdate{
		Status: "Updated",
	}, nil

}

func (s *AuthService) DeleteStudentGroup(ctx context.Context, req models.RemoveRequest) (*models.RemoveResponse, error) {
	// Input validation
	if req.ID == nil {
		return nil, fmt.Errorf("missing required fields: id")
	}
	query := `DELETE FROM stu_tracker.Student_groups WHERE id = $1`
	_, err := s.db.ExecContext(ctx, query, req.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to delete student: %w", err)
	}
	return &models.RemoveResponse{
		Status: "Removed",
	}, nil
}
