package services

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"
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

	var sessionID int64
	var rowsAffected int64
	//	var assessmentsCompleted *int

	// Main session insertion
	err = tx.QueryRowContext(
		ctx,
		`INSERT INTO stu_tracker.Sessions(
            tutor_id, session_date, location_id, substitute, 
            start_time, notes, program_id, student_count, 
            organization_id, semester_id, duration, in_school, substitute_id, session_token
        ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14) 
        RETURNING id;`,
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
		req.SessionToken,
	).Scan(&sessionID)
	if err != nil {
		return nil, fmt.Errorf("failed to insert session: %w", err)
	}
	// Insert students into sessions attendance.
	studentQuery, values, err := StudentSessionAttendance(sessionID, &req.SessionList)
	result, err := tx.ExecContext(ctx, studentQuery, values...)
	if err != nil {
		return nil, fmt.Errorf("failed to insert session students: %w", err)
	}
	rowsAffected, err = result.RowsAffected()
	if err != nil {
		return nil, fmt.Errorf("failed to get rows affected: %w", err)
	}

	assessmentQuery, values, err := StudentEasyScoreAssessmentSubmit(sessionID, &req.SessionList)
	assessmentQueryResult, err := tx.ExecContext(ctx, assessmentQuery, values...)
	if err != nil {
		return nil, fmt.Errorf("failed to insert session students: %w", err)
	}
	_, err = assessmentQueryResult.RowsAffected()
	if err != nil {
		return nil, fmt.Errorf("failed to get rows affected: %w", err)
	}

	var assessmentCount int64
	// Handle session insert per student
	/*
		if len(req.SessionList) > 0 {
			assessmentsCompleted, err = s.ProccessAssessmentsSessions(ctx, tx, &req.SessionList, sessionID, req.Assessments, req.SessionToken)
			if assessmentsCompleted != nil {
				assessmentCount = int64(*assessmentsCompleted)
			} else {
				assessmentCount += 0
			}
		}
		fmt.Println("AssessmentCount", assessmentCount)
		// Cleanup operations
		if req.SessionToken != nil {
			err = CleanSessionData(ctx, tx, *req.SessionToken)
			if err != nil {
				return nil, err
			}
		}

		// Commit the transaction if everything succeeded

	*/
	if err = tx.Commit(); err != nil {
		fmt.Println(err)
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}
	return &models.ResponseStudentSession{
		Status:          "OK",
		StudentCount:    rowsAffected,
		AssessmentCount: assessmentCount,
		SessionID:       &sessionID,
	}, nil
}

func (s *AuthService) CheckDuplicateSession(req models.RegisterStudentSessionList) (bool, error) {
	var exists bool
	// Check if the same submission has been made for the same location and tutor at the same time.
	duplicateQuery := `
    SELECT EXISTS (
        SELECT 1 FROM stu_tracker.Sessions 
        WHERE session_date = $1 AND start_time = $2 AND tutor_id = $3 AND location_id = $4)`
	err := s.db.QueryRow(duplicateQuery,
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

func (s AuthService) ProccessAssessmentsSessions(ctx context.Context, tx *sql.Tx,
	SessionList *[]models.RegisterStudentSession, sessionID int64,
	AssessmentList map[string]*models.AssessmentPayload, sessionToken *string) (*int, error) {
	var assessmentsCompleted int
	completed := &assessmentsCompleted
	for _, student := range *SessionList {
		if student.AssessmentId == nil {
			continue
		}
		var score float32
		var questionEntires []models.AnswerFeedback
		studentIdStr := strconv.Itoa(int(*student.ID))
		choices, ok := AssessmentList[studentIdStr]
		if student.EasyScoreID && ok {
			if !isMapTrulyEmpty(choices.Choices) {
				assessmentScores, err := s.ComputeScore(choices.AssessmentID, choices.Choices, choices.Grader)
				fmt.Println("ComputeScore ProcessAssessmentsSessions", assessmentScores.Points, len(assessmentScores.QuestionEntries))
				if err != nil {
					return nil, err
				}
				score += assessmentScores.Points
				questionEntires = assessmentScores.QuestionEntries
			}
		} else if student.EasyScoreID {
			// sessionToken can be null
			choiceList, err := s.GetAssessmentChoicesByStudent(ctx, student.ID, sessionToken)
			if err != nil {
				return nil, err
			}
			assessmentId, err := s.GetStudentAssessmentId(ctx, student.ID, sessionToken)
			if err != nil {
				return nil, fmt.Errorf("student failed to complete assessment. Please remove and re-submit")
			}

			choiceMap := make(map[string]interface{})
			for i := 0; i < len(choiceList); i++ {
				questionID := strconv.Itoa(int(*choiceList[i].QuestionID))
				if choiceList[i].ChoiceID != nil {
					exist, found := choiceMap[questionID]
					if !found {
						choiceMap[questionID] = int64(*choiceList[i].ChoiceID)
					} else {
						switch v := exist.(type) {
						case int64:
							choiceMap[questionID] = []int64{v, *choiceList[i].ChoiceID}
						case []int64:
							choiceMap[questionID] = append(v, *choiceList[i].ChoiceID)
						}
					}
				} else if choiceList[i].AnswerText != nil {
					choiceMap[questionID] = *choiceList[i].AnswerText
				}
			}
			var grader = &map[string]bool{}
			if choices == nil {
				fmt.Println("Grader is null")
			} else {
				grader = &choices.Grader
			}
			assessmentScores, err := s.ComputeScore(assessmentId, choiceMap, *grader)
			if err != nil {
				return nil, err
			}
			if assessmentScores != nil {
				score += assessmentScores.Points
				questionEntires = assessmentScores.QuestionEntries
			}

		} else if student.AssessmentScore != nil {
			score += float32(*student.AssessmentScore)
		}
		questionnaireId, err := QuestionnaireExist(ctx, student.ID, student.AssessmentId, sessionToken, tx)
		if err != nil {
			fmt.Printf("Error in QuestionnairExist %v", err)
			return nil, err
		}
		insertedId, err := InsertAssessmentStudents(ctx, student, score, sessionID, tx, questionnaireId)
		if err != nil {
			return nil, err
		}
		if insertedId != nil {
			insertedQId, err := InsertAssessmentQuestionsEntries(ctx, questionEntires, insertedId, tx)
			if err != nil {
				fmt.Println(err)
				return nil, err
			}
			fmt.Println("InsertAssessmentQuestionsEntries ID:", insertedQId)

			if insertedQId != nil {
				*completed += 1
			}
		}

	}
	return completed, nil
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

func StudentSessionAttendance(SessionID int64, SessionList *[]models.RegisterStudentSession) (string, []interface{}, error) {
	values := []interface{}{}
	studentQuery := `INSERT INTO stu_tracker.Session_students(
            session_id, student_id, duration, subject_id, 
            timeframe, timeframe_start, timeframe_end, absent
        ) VALUES `
	studentPlaceHolderIdx := 1
	for i, student := range *SessionList {
		if i > 0 {
			studentQuery += ", "
		}
		studentQuery += fmt.Sprintf("($%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d)",
			studentPlaceHolderIdx, studentPlaceHolderIdx+1, studentPlaceHolderIdx+2,
			studentPlaceHolderIdx+3, studentPlaceHolderIdx+4, studentPlaceHolderIdx+5,
			studentPlaceHolderIdx+6, studentPlaceHolderIdx+7)

		values = append(values,
			SessionID, student.ID, student.Duration, student.SubjectId,
			student.Timeframe, student.TimeframeStart, student.TimeframeEnd, student.Absent)
		studentPlaceHolderIdx += 8
	}
	studentQuery += ` ON CONFLICT (session_id, student_id) DO NOTHING`
	return studentQuery, values, nil

}

func StudentEasyScoreAssessmentSubmit(SessionID int64, SessionList *[]models.RegisterStudentSession) (string, []interface{}, error) {
	values := []interface{}{}
	studentQuery := `INSERT INTO stu_tracker.Assessments_students(
            session_id, student_id, score, assessment_id, 
            subject_id
        ) VALUES `
	studentPlaceHolderIdx := 1
	for i, student := range *SessionList {
		if student.EasyScoreID == false && student.AssessmentId != nil {
			if i > 0 {
				studentQuery += ", "
			}
			studentQuery += fmt.Sprintf("($%d, $%d, $%d, $%d, $%d)",
				studentPlaceHolderIdx, studentPlaceHolderIdx+1, studentPlaceHolderIdx+2,
				studentPlaceHolderIdx+3, studentPlaceHolderIdx+4)
			values = append(values,
				SessionID, student.ID, student.AssessmentScore, student.AssessmentId,
				student.SubjectId)
			studentPlaceHolderIdx += 5
		}
		studentQuery += ` ON CONFLICT (student_id, assessment_id, session_id) DO NOTHING`
	}

	return studentQuery, values, nil

}

func CleanSessionData(ctx context.Context, tx *sql.Tx, sessionToken string) error {
	if sessionToken == "" {
		return nil
	}
	queries := []struct {
		sql  string
		desc string
	}{
		{
			sql:  `DELETE FROM stu_tracker.Assessment_sessions WHERE session_token = $1`,
			desc: "assessment sessions",
		},
		{
			sql:  `DELETE FROM stu_tracker.Session_answers WHERE session_token = $1`,
			desc: "session answers",
		},
	}

	for _, query := range queries {
		_, err := tx.ExecContext(ctx, query.sql, sessionToken)
		if err != nil {
			return fmt.Errorf("failed to clean %s: %w", query.desc, err)
		}
	}

	return nil
}

func InsertAssessmentQuestionsEntries(ctx context.Context, questionEntries []models.AnswerFeedback, insertID *int64, tx *sql.Tx) (*int, error) {
	if len(questionEntries) <= 0 {
		return nil, nil
	}
	var insertedAssessmentAnswers int

	for j := 0; j < len(questionEntries); j++ {
		query2 := `
		INSERT INTO stu_tracker.Assessment_answers
		(assessment_student_id, question_id, choice_id, is_correct, answer_text)
		VALUES ($1, $2, $3, $4, $5)`
		_, err := tx.ExecContext(
			ctx,
			query2,
			*insertID,
			questionEntries[j].QuestionID,
			questionEntries[j].ChoiceID,
			questionEntries[j].IsCorrect,
			questionEntries[j].AnswerText,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to insert assessment answer: %w", err)
		}
		insertedAssessmentAnswers++
	}
	return &insertedAssessmentAnswers, nil
}

func QuestionnaireExist(ctx context.Context, studentId *int64, assessmentId *int64, sessionToken *string, tx *sql.Tx) (*int64, error) {
	if sessionToken == nil {
		return nil, nil
	}
	var questionnaire_id *int64
	var exist bool
	query := `SELECT EXISTS (
		SELECT 1 FROM stu_tracker.Pre_assessment_questionnaire 
		WHERE session_token = $1 AND student_id = $2 AND assessment_id = $3)`

	err := tx.QueryRowContext(ctx, query, sessionToken, studentId, assessmentId).Scan(&exist)
	if err != nil {
		fmt.Printf("Error select exist %v", err)
		return nil, err
	}
	if !exist {
		fmt.Printf("Does not exist %v", err)
		return nil, nil
	}
	_query := `SELECT id FROM stu_tracker.Pre_assessment_questionnaire 
		WHERE session_token = $1 AND student_id = $2 AND assessment_id = $3;`
	err = tx.QueryRowContext(ctx, _query, sessionToken, studentId, assessmentId).Scan(&questionnaire_id)
	if err != nil {
		return nil, err
	}
	return questionnaire_id, nil
}

func InsertAssessmentStudents(ctx context.Context, student models.RegisterStudentSession, score float32, sessionID int64, tx *sql.Tx, questionnaireId *int64) (*int64, error) {
	var insertedID int64
	query := `
		INSERT INTO stu_tracker.Assessments_students
		(session_id, student_id, score, assessment_id, subject_id, questionnaire_id)
		VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (student_id, assessment_id, session_id) DO NOTHING
		RETURNING id
				`
	err := tx.QueryRowContext(
		ctx,
		query,
		sessionID,
		student.ID,
		score,
		student.AssessmentId,
		student.SubjectId,
		questionnaireId,
	).Scan(&insertedID)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil // Not an error, just nothing inserted
		}
		return nil, fmt.Errorf("failed to insert individual assessment session: %w", err)
	}
	return &insertedID, nil
}
