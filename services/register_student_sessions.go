package services

import (
	"context"
	"fmt"
	"log"
	"tracker/app/models"
)

func (s *AuthService) CreateStudentSessions(req models.RegisterStudentSessionList) (*models.ResponseStudentSession, error) {
	// Input validation
	if len(req.SessionList) <= 0 {
		return nil, fmt.Errorf("missing required fields: Session list is empty")
	}
	ctx := context.Background()
	// Named parameters for better readability (if your DB driver supports it)
	query := `
    INSERT INTO stu_tracker.Sessions (
        tutor_id, 
        session_date, 
        location_id, 
        substitute, 
        start_time, 
        notes, 
        program_id, 
        student_count, 
        organization_id,
		semester_id,
		duration
    ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11) 
    RETURNING id;
`
	var newID int64
	// Execute query with context
	err := s.db.QueryRowContext(
		ctx,
		query,
		req.Session.TutorId,
		req.Session.SessionDate,
		req.Session.LocationId,
		req.Session.Substitute,
		req.Session.StartTime,
		req.Session.Notes,
		req.Session.ProgramId,
		req.Session.StudentCount,
		req.OrganizationID,
		req.Session.SemesterId,
		req.Session.Duration,
	).Scan(&newID)
	if err != nil {
		return nil, err
	}
	values := []interface{}{}
	assessment_values := []interface{}{}
	studentQuery := `INSERT INTO stu_tracker.Session_students(session_id, student_id, duration, subject_id) VALUES`
	studentPlaceHolderIdx, assessmentPlaceHolderIdx := 1, 1

	for i, student := range req.SessionList {
		if i > 0 {
			studentQuery += ", "
		}
		studentQuery += fmt.Sprintf("($%d, $%d, $%d, $%d)", studentPlaceHolderIdx, studentPlaceHolderIdx+1, studentPlaceHolderIdx+2, studentPlaceHolderIdx+3)
		values = append(values, newID, &student.ID, &student.Duration, &student.SubjectId)
		studentPlaceHolderIdx += 4
	}
	studentQuery += ` ON CONFLICT (session_id, student_id) DO NOTHING`
	result, err := s.db.Exec(studentQuery, values...)
	if err != nil {
		return nil, fmt.Errorf("failed to session students query: %w", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Fatalf("Rows affected error: %v", err)
	}
	assessmentsCompleted := 0
	// If there are sessions available log each session to each student.
	if len(req.SessionList) > 0 {
		assessmentQuery := `INSERT INTO stu_tracker.Assessments_students(session_id, student_id, score, assessment_id, subject_id) VALUES`
		for i, student := range req.SessionList {
			if student.AssessmentId != nil {
				if i > 0 {
					assessmentQuery += ", "
				}
				assessmentQuery += fmt.Sprintf("($%d, $%d, $%d, $%d, $%d)", assessmentPlaceHolderIdx, assessmentPlaceHolderIdx+1, assessmentPlaceHolderIdx+2, assessmentPlaceHolderIdx+3, assessmentPlaceHolderIdx+4)
				assessment_values = append(assessment_values, newID, &student.ID, &student.AssessmentScore, &student.AssessmentId, &student.SubjectId)
				assessmentPlaceHolderIdx += 5
			}
		}
		if len(assessment_values) > 0 {
			assessmentQuery += ` ON CONFLICT (student_id, assessment_id, session_id) DO NOTHING`
			assessmentResult, err := s.db.Exec(assessmentQuery, assessment_values...)
			if err != nil {
				return nil, fmt.Errorf("failed to assessment session students query: %w", err)
			}
			assessmentChanges, err := assessmentResult.RowsAffected()
			if err != nil {
				log.Fatalf("Rows affected error assessment: %v", err)
			}
			assessmentsCompleted += int(assessmentChanges)
		}
	}

	return &models.ResponseStudentSession{
		Status:          "OK",
		StudentCount:    int64(rowsAffected),
		AssessmentCount: int64(assessmentsCompleted),
	}, nil
}
