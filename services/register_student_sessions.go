package services

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"tracker/app/models"
)

func (s *AuthService) CreateStudentSessions(req models.RegisterStudentSessionList) (*models.ResponseStudentSession, error) {
	// Input validation
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
		duration,
		in_school,
		substitute_id
    ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13) 
    RETURNING id;`
	var newID int64
	// Execute query with context
	err = s.db.QueryRowContext(
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
		req.Session.InSchool,
		req.Session.SubstituteId,
	).Scan(&newID)
	if err != nil {
		return nil, err
	}
	values := []interface{}{}
	studentQuery := `INSERT INTO stu_tracker.Session_students(session_id, student_id, duration, subject_id) VALUES`
	studentPlaceHolderIdx := 1
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
		return nil, err
	}
	assessmentsCompleted := 0
	assessmentsAnswersCompleted := 0
	// If there are sessions available log each session to each student.
	// Will need to return an Assessments_students ID

	//INSERT ID

	if len(req.SessionList) > 0 {
		for _, student := range req.SessionList {
			if student.AssessmentId != nil {
				var score int
				var questionEntires []models.AnswerFeedback
				// If easy score is turned on then grade the assessment
				if student.EasyScoreID {
					studentIdstr := strconv.Itoa(int(*student.ID))
					// Find the choices based on student_id
					if choices, ok := req.Assessments[studentIdstr]; ok {
						// Is choice map empty, if so then fetch from assessment session
						// If choice map is not empty, then get from choice obj given by user
						var isEmpty bool = isMapTrulyEmpty(choices.Choices)
						fmt.Println("Is choice object empty", isEmpty)
						if !isEmpty {
							// Compute score if the tutor has input the assessment values
							assessmentScores, err := s.ComputeScore(choices.AssessmentID, choices.Choices, choices.Grader)
							if err != nil {
								return nil, err
							}
							fmt.Println("Question entries: ", len(assessmentScores.QuestionEntries))
							fmt.Println("Question points: ", assessmentScores.Points)
							fmt.Println("Question MaxScore: ", assessmentScores.MaxScore)

							score = assessmentScores.Points
							questionEntires = assessmentScores.QuestionEntries

						} else {
							choiceList, err := s.GetAssessmentChoicesByStudent(ctx, choices.AssessmentID, student.ID, req.SessionToken)
							if err != nil {
								return nil, err
							}
							fmt.Println("Is not empty, GetAssessmentChoicesByStudent: ", choiceList)
							choicesMap := make(map[string]interface{})
							for i := 0; i < len(choiceList); i++ {
								questionID := strconv.Itoa(int(*choiceList[i].QuestionID))
								if choiceList[i].ChoiceID != nil {
									choicesMap[questionID] = *choiceList[i].ChoiceID
								} else {
									choicesMap[questionID] = *choiceList[i].AnswerText
								}
							}
							fmt.Println("Is not empty, ChoicesMap Created: ", choicesMap)

							assessmentScores, err := s.ComputeScore(choices.AssessmentID, choicesMap, choices.Grader)
							if err != nil {
								return nil, err
							}
							fmt.Println("Question entries: ", len(assessmentScores.QuestionEntries))
							fmt.Println("Question points: ", assessmentScores.Points)
							fmt.Println("Question MaxScore: ", assessmentScores.MaxScore)
							score = assessmentScores.Points
							questionEntires = assessmentScores.QuestionEntries

						}

					}
				} else {
					score = int(*student.AssessmentScore)
				}
				var insertedID int
				var insertedAssessmentAnswers int
				// Compute score if needed to based on student.EasyScoreId
				query := `
					INSERT INTO stu_tracker.Assessments_students
					(session_id, student_id, score, assessment_id, subject_id)
					VALUES ($1, $2, $3, $4, $5)
					ON CONFLICT (student_id, assessment_id, session_id) DO NOTHING
					RETURNING id
				`
				err := s.db.QueryRowContext(
					ctx,
					query,
					newID,
					student.ID,
					score,
					student.AssessmentId,
					student.SubjectId,
				).Scan(&insertedID)

				if err != nil && err != sql.ErrNoRows {
					return nil, fmt.Errorf("failed to insert individual assessment session: %w", err)
				}
				if len(questionEntires) > 0 {
					for j := 0; j < len(questionEntires); j++ {
						query2 := `
						INSERT INTO stu_tracker.Assessment_answers
						(assessment_student_id, question_id, choice_id, is_correct, answer_text)
						VALUES ($1, $2, $3, $4, $5);`
						err := s.db.QueryRowContext(
							ctx,
							query2,
							insertedID,
							questionEntires[j].QuestionID,
							questionEntires[j].ChoiceID,
							questionEntires[j].IsCorrect,
							questionEntires[j].AnswerText,
						).Scan(&insertedAssessmentAnswers)

						if err != nil && err != sql.ErrNoRows {
							return nil, fmt.Errorf("failed to insert individual assessment session: %w", err)
						}
					}
				}
				if insertedID != 0 {
					assessmentsCompleted++
				}
				if insertedAssessmentAnswers != 0 {
					assessmentsAnswersCompleted++
				}

			}
		}
	}

	if req.SessionToken != nil {
		assessmentQueryRemove := `DELETE FROM stu_tracker.Assessment_sessions WHERE session_token = $1`
		sessionQueryRemove := `DELETE FROM stu_tracker.Session_answers WHERE session_token = $1`
		_, err = s.db.ExecContext(ctx, assessmentQueryRemove, req.SessionToken)
		if err != nil {
			return nil, err
		}
		_, err = s.db.ExecContext(ctx, sessionQueryRemove, req.SessionToken)
		if err != nil {
			return nil, err
		}
	}

	return &models.ResponseStudentSession{
		Status:          "OK",
		StudentCount:    int64(rowsAffected),
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

func (s *AuthService) insertAssessmentStudent(
	ctx context.Context,
	sessionID int64,
	student models.RegisterStudentSession,
	score int,
) (int, error) {
	var insertedID int
	query := `
		INSERT INTO stu_tracker.Assessments_students
		(session_id, student_id, score, assessment_id, subject_id)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (student_id, assessment_id, session_id) DO NOTHING
		RETURNING id`
	err := s.db.QueryRowContext(ctx, query,
		sessionID, student.ID, score, student.AssessmentId, student.SubjectId,
	).Scan(&insertedID)

	if err != nil && err != sql.ErrNoRows {
		return 0, fmt.Errorf("failed to insert assessment student: %w", err)
	}
	return insertedID, nil
}

func (s *AuthService) insertAssessmentAnswers(
	ctx context.Context,
	assessmentStudentID int,
	answers []models.AnswerFeedback,
) error {
	for _, ans := range answers {
		query := `
			INSERT INTO stu_tracker.Assessment_answers
			(assessment_student_id, question_id, choice_id, is_correct)
			VALUES ($1, $2, $3, $4)
		`
		_, err := s.db.ExecContext(ctx, query,
			assessmentStudentID,
			ans.QuestionID,
			ans.ChoiceID,
			ans.IsCorrect,
		)
		if err != nil {
			return fmt.Errorf("failed to insert assessment answer: %w", err)
		}
	}
	return nil
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

func (s *AuthService) removeActiveSessions(tutor_id *int64, c context.Context) (bool, error) {
	query := `DELTE FROM stu_tracker.Assessment_sessions WHERE tutor_id = $1;`
	r, err := s.db.ExecContext(c, query, tutor_id)
	if err != nil {
		return false, err
	}
	changed, err := r.RowsAffected()
	if err != nil {
		return false, err
	}
	if changed > 0 {
		return true, nil
	}
	return false, nil
}
