package services

import (
	"context"
	"fmt"
	"tracker/app/models"

	"github.com/google/uuid"
)

func (s *AuthService) CreateAssessmentSession(c context.Context, req models.RegisterStudentAssessmentSession) (*models.ResponseStudentAssessmentSession, error) {
	// Input validation
	if len(req.StudentAssessmentSession) <= 0 {
		return nil, fmt.Errorf("no student sessions provided")
	}
	sessionToken := uuid.New() // generate a UUID
	for _, student := range req.StudentAssessmentSession {
		var exists bool
		checkQuery := `
			SELECT EXISTS (
				SELECT 1 FROM stu_tracker.Assessment_sessions
				WHERE tutor_id = $1 AND assessment_id = $2 AND student_id = $3
			)
		`
		err := s.db.QueryRowContext(c, checkQuery, student.TutorId, student.AssessmentId, student.ID).Scan(&exists)
		if err != nil {
			return nil, fmt.Errorf("failed to check existing session: %w", err)
		}
		if exists {
			deleteQuery := `DELETE FROM stu_tracker.Assessment_sessions
								WHERE tutor_id = $1 AND assessment_id = $2 AND student_id = $3`
			_, err = s.db.ExecContext(c, deleteQuery, student.TutorId, student.AssessmentId, student.ID)
			if err != nil {
				return nil, err
			}
		}
		var toGrade bool
		gradeCheck := `SELECT EXISTS (SELECT 1 FROM stu_tracker.Questions WHERE assessment_id = $1 AND question_type = 'short_answer');`
		err = s.db.QueryRowContext(c, gradeCheck, student.AssessmentId).Scan(&toGrade)
		if err != nil {
			return nil, fmt.Errorf("unable to check if question type")
		}
		query := `
        INSERT INTO stu_tracker.Assessment_sessions (
            tutor_id, student_id, first_name, last_name,
            assessment_id, semester_id, session_token, grade_assessment
        )
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
    `
		_, err = s.db.ExecContext(c, query, student.TutorId, student.ID, student.FirstName, student.LastName, student.AssessmentId, student.SemesterId, sessionToken, toGrade)
		if err != nil {
			return nil, fmt.Errorf("failed to add semester: %w", err)
		}
	}
	return &models.ResponseStudentAssessmentSession{
		Status:         "OK",
		SessionsActive: len(req.StudentAssessmentSession),
		SessionID:      sessionToken.String(),
	}, nil
}
