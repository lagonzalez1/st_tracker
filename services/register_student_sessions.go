package services

import (
	"context"
	"database/sql"
	"fmt"
	"tracker/app/models"
)

func (s *AuthService) CreateStudentSessions(req models.RegisterStudentSessionList) (*models.ResponseStudentSession, error) {
	// Input validation remains the same
	if len(req.SessionList) <= 0 {
		return nil, fmt.Errorf("missing required fields: Session list is empty")
	}

	duplicate, err := s.CheckDuplicateSession(req)
	if err != nil {
		return nil, err
	}
	if duplicate {
		return nil, fmt.Errorf("found a duplicate entry")
	}

	ctx := context.Background()

	// Begin a transaction
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}

	// Defer rollback in case of failure
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	var newID int64
	var rowsAffected int64
	var assessmentsCompleted, assessmentsAnswersCompleted int

	// Main session insertion
	err = tx.QueryRowContext(
		ctx,
		`INSERT INTO stu_tracker.Sessions (
            tutor_id, session_date, location_id, substitute, 
            start_time, notes, program_id, student_count, 
            organization_id, semester_id, duration, in_school, substitute_id
        ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13) 
        RETURNING id`,
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
		req.Session.InSchool,
		req.Session.SubstituteId,
	).Scan(&newID)
	if err != nil {
		return nil, fmt.Errorf("failed to insert session: %w", err)
	}

	// Student sessions insertion
	if len(req.SessionList) > 0 {
		values := []interface{}{}
		studentQuery := `INSERT INTO stu_tracker.Session_students(
            session_id, student_id, duration, subject_id, 
            timeframe, timeframe_start, timeframe_end
        ) VALUES `

		studentPlaceHolderIdx := 1
		for i, student := range req.SessionList {
			if i > 0 {
				studentQuery += ", "
			}
			studentQuery += fmt.Sprintf("($%d, $%d, $%d, $%d, $%d, $%d, $%d)",
				studentPlaceHolderIdx, studentPlaceHolderIdx+1, studentPlaceHolderIdx+2,
				studentPlaceHolderIdx+3, studentPlaceHolderIdx+4, studentPlaceHolderIdx+5,
				studentPlaceHolderIdx+6)

			values = append(values,
				newID, student.ID, student.Duration, student.SubjectId,
				student.Timeframe, student.TimeframeStart, student.TimeframeEnd)
			studentPlaceHolderIdx += 7
		}

		studentQuery += ` ON CONFLICT (session_id, student_id) DO NOTHING`

		result, err := tx.ExecContext(ctx, studentQuery, values...)
		if err != nil {
			return nil, fmt.Errorf("failed to insert session students: %w", err)
		}

		rowsAffected, err = result.RowsAffected()
		if err != nil {
			return nil, fmt.Errorf("failed to get rows affected: %w", err)
		}
	}

	// Assessment processing
	for _, student := range req.SessionList {
		if student.AssessmentId != nil {
			var score int
			var questionEntries []models.AnswerFeedback

			if student.EasyScoreID {
				// [Keep your existing score calculation logic]
				// Just replace s.db with tx for all database operations
			} else {
				score = int(*student.AssessmentScore)
			}

			var insertedID int
			err := tx.QueryRowContext(
				ctx,
				`INSERT INTO stu_tracker.Assessments_students
                (session_id, student_id, score, assessment_id, subject_id)
                VALUES ($1, $2, $3, $4, $5)
                ON CONFLICT (student_id, assessment_id, session_id) DO NOTHING
                RETURNING id`,
				newID, student.ID, score, student.AssessmentId, student.SubjectId,
			).Scan(&insertedID)

			if err != nil && err != sql.ErrNoRows {
				return nil, fmt.Errorf("failed to insert assessment: %w", err)
			}

			if len(questionEntries) > 0 {
				for _, entry := range questionEntries {
					_, err := tx.ExecContext(
						ctx,
						`INSERT INTO stu_tracker.Assessment_answers
                        (assessment_student_id, question_id, choice_id, is_correct, answer_text)
                        VALUES ($1, $2, $3, $4, $5)`,
						insertedID, entry.QuestionID, entry.ChoiceID, entry.IsCorrect, entry.AnswerText,
					)
					if err != nil {
						return nil, fmt.Errorf("failed to insert assessment answer: %w", err)
					}
					assessmentsAnswersCompleted++
				}
			}

			if insertedID != 0 {
				assessmentsCompleted++
			}
		}
	}

	// Cleanup operations
	if req.SessionToken != nil {
		_, err = tx.ExecContext(ctx,
			`DELETE FROM stu_tracker.Assessment_sessions WHERE session_token = $1`,
			req.SessionToken)
		if err != nil {
			return nil, fmt.Errorf("failed to clean assessment sessions: %w", err)
		}

		_, err = tx.ExecContext(ctx,
			`DELETE FROM stu_tracker.Session_answers WHERE session_token = $1`,
			req.SessionToken)
		if err != nil {
			return nil, fmt.Errorf("failed to clean session answers: %w", err)
		}
	}

	// Commit the transaction if everything succeeded
	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return &models.ResponseStudentSession{
		Status:          "OK",
		StudentCount:    rowsAffected,
		AssessmentCount: int64(assessmentsCompleted),
	}, nil
}

func (s *AuthService) CheckDuplicateSession(req models.RegisterStudentSessionList) (bool, error) {
	var exists bool
	// Check if the same submission has been made for the same location and tutor at the same time.
	err := s.db.QueryRow(`
    SELECT EXISTS (
        SELECT 1 FROM stu_tracker.Sessions 
        WHERE session_date = $1 AND start_time = $2 AND tutor_id = $3 AND location_id = $4)`,
		req.Session.SessionDate, req.Session.StartTime, req.Session.TutorId, req.Session.LocationId).Scan(&exists)

	if err != nil {
		return true, fmt.Errorf("unable to query to check duplicate sessions: %v", err)
	}
	if exists {
		return true, fmt.Errorf("duplicate session found")
	}
	if !exists {
		return false, nil
	}
	return false, nil
}
func isMapTrulyEmpty(m map[string]interface{}) bool {
	if len(m) == 0 {
		return true
	}

	// Check if all values are empty maps/interfaces
	for _, v := range m {
		switch val := v.(type) {
		case map[string]interface{}:
			if !isMapTrulyEmpty(val) {
				return false
			}
		default:
			// If any non-map, non-zero value exists
			if val != nil {
				return false
			}
		}
	}
	return true
}
