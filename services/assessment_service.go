package services

import (
	"context"
	"fmt"
	"strconv"
	"tracker/app/models"
)

func (s *AuthService) ComputeScore(assessmentID *int64, choices map[string]interface{}, grader map[string]bool) (*models.AssessmentScore, error) {
	if assessmentID == nil {
		return nil, fmt.Errorf("assessmentID cannot be nil")
	}

	maxScore, err := s.GetAssessmentMaxScore(assessmentID)
	if err != nil {
		return nil, fmt.Errorf("unable to get max score for assessment %d: %w", *assessmentID, err)
	}
	if maxScore == nil {
		return nil, fmt.Errorf("maxScore was nil for assessment %d", *assessmentID)
	}

	if len(choices) == 0 {
		fmt.Println("Empty choices, returning zero score")
		return &models.AssessmentScore{
			Points:          0.0,
			MaxScore:        *maxScore,
			QuestionEntries: []models.AnswerFeedback{},
		}, nil
	}

	points, questionEntries, err := s.GradeAssessmentWithCorrectAnswers(assessmentID, choices, grader)
	if err != nil {
		return nil, fmt.Errorf("grading failed for assessment %d: %w", *assessmentID, err)
	}
	if points == nil {
		return nil, fmt.Errorf("grade result returned nil points for assessment %d", *assessmentID)
	}

	return &models.AssessmentScore{
		Points:          *points,
		MaxScore:        *maxScore,
		QuestionEntries: questionEntries,
	}, nil
}

func (s *AuthService) GradeAssessment(assessment_id *int64, choices map[string]int) (*int, error) {
	if assessment_id == nil {
		return nil, fmt.Errorf("assessment is null")
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	query := `SELECT 
		c.id AS choice_id,
		c.question_id,
		c.is_correct,
		q.points	
	FROM stu_tracker.Choices c
	INNER JOIN stu_tracker.Questions q 
	ON c.question_id = q.id
	WHERE q.assessment_id = $1
	ORDER BY c.question_id, c.order_number;`

	rows, err := s.db.QueryContext(ctx, query, assessment_id)
	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}
	defer rows.Close()
	var questions []models.AssessmentGrader
	for rows.Next() {
		var r models.AssessmentGrader
		err := rows.Scan(
			&r.ChoiceID,
			&r.QuestionID,
			&r.IsCorrect,
			&r.Points,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		questions = append(questions, r)
	}
	var totalScore int // use a normal int to accumulate
	for _, q := range questions {
		questionIDstr := strconv.Itoa(int(*q.QuestionID))
		if selectedChoiceID, ok := choices[questionIDstr]; ok {
			if q.IsCorrect && selectedChoiceID == int(*q.ChoiceID) {
				totalScore += q.Points
			}
		}
	}

	return &totalScore, nil
}

func (s *AuthService) GradeAssessmentWithCorrectAnswers(
	assessmentID *int64,
	choices interface{},
	grader map[string]bool,
) (*float32, []models.AnswerFeedback, error) {
	if assessmentID == nil {
		return nil, nil, fmt.Errorf("assessment_id cannot be nil")
	}

	choicesMap, ok := choices.(map[string]interface{})
	if !ok {
		return nil, nil, fmt.Errorf("choices must be a map[string]interface{}")
	}
	fmt.Println("Choices Map 2", choicesMap)
	questions, err := s.fetchCorrectAnswers(*assessmentID)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to fetch correct answers: %w", err)
	}
	fmt.Println("Questions answers", questions)
	var totalScore float32
	var allAnswers []models.AnswerFeedback

	for questionID, questionSet := range questions {
		value, ok := choicesMap[questionID]
		if !ok {
			continue
		}

		if len(questionSet) == 1 {
			feedback, score, err := s.gradeQuestion(questionSet[0], questionID, value, grader)
			if err != nil {
				return nil, nil, fmt.Errorf("grading failed for question %s: %w", questionID, err)
			}
			totalScore += float32(score)
			allAnswers = append(allAnswers, feedback)

		} else if len(questionSet) > 1 {
			rawSlice, ok := value.([]interface{})
			if !ok {
				feedback, score, err := s.gradeQuestion(questionSet[0], questionID, value, grader)
				if err != nil {
					return nil, nil, fmt.Errorf("grading failed for question %s: %w", questionID, err)
				}
				totalScore += float32(score)
				allAnswers = append(allAnswers, feedback)
				continue
			}

			lookup := make(map[int64]bool)
			for _, val := range rawSlice {
				switch v := val.(type) {
				case int64:
					lookup[v] = true
				case int:
					lookup[int64(v)] = true
				case float64:
					lookup[int64(v)] = true
				}
			}

			for _, q := range questionSet {
				if q.ChoiceID != nil && lookup[*q.ChoiceID] {
					feedback, score, err := s.gradeQuestion(q, questionID, *q.ChoiceID, grader)
					if err != nil {
						return nil, nil, fmt.Errorf("grading failed for question %s: %w", questionID, err)
					}
					divisor := float32(len(questionSet))
					totalScore += float32(score) / divisor
					allAnswers = append(allAnswers, feedback)
				}
			}
		}
	}

	return &totalScore, allAnswers, nil
}

// Helper function to fetch correct answers from DB
func (s *AuthService) fetchCorrectAnswers(assessment_id int64) (map[string][]models.AssessmentGrader, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	query := `SELECT 
        c.id AS choice_id,
        c.question_id,
        c.is_correct,
        q.points,
		q.question_type
        FROM stu_tracker.Choices c
        INNER JOIN stu_tracker.Questions q ON c.question_id = q.id
        WHERE q.assessment_id = $1 AND c.is_correct = TRUE
        ORDER BY c.question_id, c.order_number;`

	rows, err := s.db.QueryContext(ctx, query, assessment_id)
	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}
	defer rows.Close()

	var questions []models.AssessmentGrader
	for rows.Next() {
		var r models.AssessmentGrader
		if err := rows.Scan(
			&r.ChoiceID,
			&r.QuestionID,
			&r.IsCorrect,
			&r.Points,
			&r.QuestionType,
		); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		questions = append(questions, r)
	}

	q := make(map[string][]models.AssessmentGrader)

	if len(questions) > 0 {
		for _, question := range questions {
			questionIdStr := strconv.Itoa(int(*question.QuestionID))
			// Exist in map
			if payload, ok := q[questionIdStr]; ok {
				currentArray := payload
				currentArray = append(currentArray, question)
				q[questionIdStr] = currentArray
			} else {
				var array []models.AssessmentGrader
				array = append(array, question)
				q[questionIdStr] = array
			}
		}
	}

	return q, nil
}

func (s *AuthService) gradeQuestion(q models.AssessmentGrader, questionIDstr string, selectedChoice interface{}, graderMap map[string]bool) (models.AnswerFeedback, float32, error) {
	fmt.Printf("Type: %T, Value: %v\n", selectedChoice, selectedChoice)
	switch val := selectedChoice.(type) {
	case float64:
		isCorrect := q.IsCorrect && (int(val) == int(*q.ChoiceID))
		feedback := models.AnswerFeedback{
			QuestionID: *q.QuestionID,
			ChoiceID:   q.ChoiceID,
			IsCorrect:  isCorrect,
			AnswerText: nil,
		}
		if isCorrect {
			return feedback, float32(q.Points), nil
		}
		return feedback, 0, nil
	case int64:
		isCorrect := q.IsCorrect && (int(val) == int(*q.ChoiceID))
		feedback := models.AnswerFeedback{
			QuestionID: *q.QuestionID,
			ChoiceID:   q.ChoiceID,
			IsCorrect:  isCorrect,
			AnswerText: nil,
		}
		if isCorrect {
			return feedback, float32(q.Points), nil
		}
		return feedback, 0.0, nil
	case string:
		if graderMap != nil {
			correctBool, ok := graderMap[questionIDstr]
			if ok {
				feedback := models.AnswerFeedback{
					QuestionID: *q.QuestionID,
					IsCorrect:  correctBool,
					AnswerText: &val,
					ChoiceID:   nil,
				}
				if correctBool {
					return feedback, float32(q.Points), nil
				}
			}
		}
		return models.AnswerFeedback{
			QuestionID: *q.QuestionID,
			IsCorrect:  false,
			AnswerText: nil,
			ChoiceID:   nil,
		}, 0, nil
	case nil:
		if graderMap != nil {
			correctBool, ok := graderMap[questionIDstr]
			if ok {
				feedback := models.AnswerFeedback{
					QuestionID: *q.QuestionID,
					IsCorrect:  correctBool,
					AnswerText: nil,
					ChoiceID:   nil,
				}
				if correctBool {
					return feedback, float32(q.Points), nil
				}
			}
		}
		return models.AnswerFeedback{
			QuestionID: *q.QuestionID,
			IsCorrect:  false,
			AnswerText: nil,
			ChoiceID:   nil,
		}, 0, nil

	default:
		return models.AnswerFeedback{}, float32(0), fmt.Errorf("invalid choice type %T for question", selectedChoice)
	}

}

func (s *AuthService) GetAssessmentMaxScore(assessment_id *int64) (*int, error) {
	if assessment_id == nil {
		return nil, fmt.Errorf("assessment id is null")
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	var score *int
	query := `SELECT SUM(q.points) FROM stu_tracker.Questions q
	WHERE q.assessment_id = $1;`
	row := s.db.QueryRowContext(ctx, query, assessment_id)
	err := row.Scan(&score)
	if err != nil {
		return nil, err
	}

	return score, nil
}

func (s *AuthService) GetAssessmentChoicesByStudent(c context.Context, student_id *int64, session_token *string) ([]models.StudentAssessmentChoices, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	query := `SELECT question_id, choice_id, answer_text, assessment_id
	FROM stu_tracker.Session_answers 
	WHERE session_token = $1 AND student_id = $2;`
	rows, err := s.db.QueryContext(ctx, query, session_token, student_id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var assessmentChoices []models.StudentAssessmentChoices
	for rows.Next() {
		var choices models.StudentAssessmentChoices
		err := rows.Scan(
			&choices.QuestionID,
			&choices.ChoiceID,
			&choices.AnswerText,
			&choices.AssessmentID,
		)
		if err != nil {
			return nil, err
		}
		assessmentChoices = append(assessmentChoices, choices)
	}
	return assessmentChoices, nil
}

func (s *AuthService) GetStudentAssessmentId(c context.Context, student_id *int64, session_token *string) (*int64, error) {
	var id int64
	query := `SELECT assessment_id
	FROM stu_tracker.Assessment_sessions 
	WHERE session_token = $1 AND student_id = $2;`
	err := s.db.QueryRowContext(c, query, session_token, student_id).Scan(&id)
	if err != nil {
		return nil, err
	}
	return &id, nil
}

func (s *AuthService) GetAssessmentBySessionId(c context.Context, session_id string, tutor_id *int64) ([]models.StudentAssessmentSession, error) {
	var students []models.StudentAssessmentSession
	query := `
	SELECT
		s.id,
		s.first_name,
		s.last_name,
		ss.assessment_id,
		ss.is_active,
		ss.completed,
		ss.grade_assessment,
		ast.title
	FROM
		stu_tracker.Assessment_sessions ss
	INNER JOIN 
		stu_tracker.Students s	
	ON 
		s.id = ss.student_id
	JOIN 
		stu_tracker.Assessments ast
	ON 
		ast.id = ss.assessment_id
	WHERE
		ss.session_token = $1 AND ss.tutor_id = $2;`

	rows, err := s.db.QueryContext(c, query, session_id, tutor_id)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var student models.StudentAssessmentSession
		err := rows.Scan(
			&student.ID,
			&student.FirstName,
			&student.FirstName,
			&student.AssessmentID,
			&student.IsActive,
			&student.Completed,
			&student.GradeAssessment,
			&student.AssessmentTitle,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning row: %w", err)
		}
		students = append(students, student)
	}

	return students, nil
}

func (s AuthService) DeleteAssessmentSession(c context.Context, req models.DeleteAssessmentSession) (*models.DeleteAssessmentSessionResponse, error) {
	if req.SessionID == "" {
		return nil, fmt.Errorf("no session id provided")
	}
	query := `DELETE FROM stu_tracker.Assessment_sessions WHERE session_token = $1;`
	_, err := s.db.ExecContext(c, query, req.SessionID)
	if err != nil {
		return nil, err
	}
	sessionAnswerQuery := `DELETE FROM stu_tracker.Session_answers WHERE session_token = $1;`
	_, err = s.db.ExecContext(c, sessionAnswerQuery, req.SessionID)
	if err != nil {
		return nil, err
	}
	return &models.DeleteAssessmentSessionResponse{
		Status: "Deleted",
	}, nil
}

func (s *AuthService) CreateStudentAssessmentResponse(c context.Context, req models.RegisterStudentAssessment) (*models.StudentSubmitResponse, error) {
	var exists bool
	checkQueryValidSession := `
			SELECT EXISTS (
				SELECT 1 FROM stu_tracker.Assessment_sessions
				WHERE tutor_id = $1 AND assessment_id = $2 AND student_id = $3
			)
	`
	err := s.db.QueryRowContext(c, checkQueryValidSession, req.TutorID, req.AssessmentID, req.StudentID).Scan(&exists)
	if err != nil {
		return nil, fmt.Errorf("failed to check existing session: %w", err)
	}
	if exists {
		return nil, fmt.Errorf("session not found given session id")
	}
	fmt.Println("checkQueryValidSession", exists)

	var existDuplicate bool
	checkDuplicateSubmit := `SELECT EXISTS (SELECT 1 FROM stu_tracker.Session_answers WHERE session_token = $1 AND student_id = $2 AND assessment_id = $3)`
	err = s.db.QueryRowContext(c, checkDuplicateSubmit, req.SessionID, req.StudentID, req.AssessmentID).Scan(&existDuplicate)
	if err != nil {
		return nil, fmt.Errorf("failed to check existing session: %w", err)
	}
	fmt.Println("Exist duplicate", existDuplicate)
	if existDuplicate {
		return nil, fmt.Errorf("no more than one submission per session")
	}

	for questionIDStr, choiceID := range req.Answers {
		switch val := choiceID.(type) {
		case []interface{}:
			questionID, err := strconv.ParseInt(questionIDStr, 10, 64)
			if err != nil {
				return nil, err
			}
			idx := 1
			values := []interface{}{}
			insertQuery := `INSERT INTO stu_tracker.Session_answers
				(assessment_id, student_id, session_token, question_id, choice_id) VALUES `
			for i, choiceID := range val {
				if i > 0 {
					insertQuery += ", "
				}
				insertQuery += fmt.Sprintf("($%d, $%d, $%d, $%d, $%d)", idx, idx+1, idx+2, idx+3, idx+4)
				values = append(values, req.AssessmentID, req.StudentID, req.SessionID, questionID, choiceID)
				idx += 5
			}
			_, err = s.db.ExecContext(c, insertQuery, values...)
			if err != nil {
				return nil, err
			}
		case float64:
			insertQuery := `
			INSERT INTO stu_tracker.Session_answers
				(assessment_id, student_id, session_token, question_id, choice_id)
			VALUES ($1, $2, $3, $4, $5);`
			questionID, err := strconv.ParseInt(questionIDStr, 10, 64)

			if err != nil {
				return nil, fmt.Errorf("invalid question ID '%s': %w", questionIDStr, err)
			}
			_, err = s.db.ExecContext(c, insertQuery,
				req.AssessmentID,
				req.StudentID,
				req.SessionID,
				questionID,
				choiceID,
			)
			if err != nil {
				return nil, err
			}
		case int64:
			insertQuery := `
			INSERT INTO stu_tracker.Session_answers
				(assessment_id, student_id, session_token, question_id, choice_id)
			VALUES ($1, $2, $3, $4, $5);`
			questionID, err := strconv.ParseInt(questionIDStr, 10, 64)

			if err != nil {
				return nil, fmt.Errorf("invalid question ID '%s': %w", questionIDStr, err)
			}
			_, err = s.db.ExecContext(c, insertQuery,
				req.AssessmentID,
				req.StudentID,
				req.SessionID,
				questionID,
				choiceID,
			)
			if err != nil {
				return nil, err
			}
		case string:
			insertQuery := `
			INSERT INTO stu_tracker.Session_answers
				(assessment_id, student_id, session_token, question_id, answer_text)
			VALUES ($1, $2, $3, $4, $5);`
			questionID, err := strconv.ParseInt(questionIDStr, 10, 64)

			if err != nil {
				return nil, fmt.Errorf("invalid question ID '%s': %w", questionIDStr, err)
			}
			_, err = s.db.ExecContext(c, insertQuery,
				req.AssessmentID,
				req.StudentID,
				req.SessionID,
				questionID,
				val,
			)
			if err != nil {
				return nil, err
			}
		}
	}
	completedQuery := `UPDATE stu_tracker.Assessment_sessions SET completed = $1 WHERE student_id = $2 AND session_token = $3;`
	_, err = s.db.ExecContext(c, completedQuery, true, req.StudentID, req.SessionID)
	if err != nil {
		return nil, fmt.Errorf("failed to insert admin: %w", err)
	}
	return &models.StudentSubmitResponse{
		Status:  "OK",
		Answers: len(req.Answers),
	}, nil
}

func (s *AuthService) DeleteStudentSession(c context.Context, req models.DeleteStudentSession) (*models.DeleteStudentSessionResponse, error) {
	if req.SessionID == "" || req.StudentID == nil {
		return nil, fmt.Errorf("missing session id and or student id")
	}
	tx, err := s.db.BeginTx(c, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	// Defer rollback in case of failure
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	query := `DELETE FROM stu_tracker.Assessment_sessions WHERE session_token = $1 AND student_id = $2;`
	_, err = tx.ExecContext(c, query, req.SessionID, req.StudentID)
	if err != nil {
		return nil, err
	}
	sessionAnswerQuery := `DELETE FROM stu_tracker.Session_answers WHERE session_token = $1 AND student_id = $2;`
	_, err = tx.ExecContext(c, sessionAnswerQuery, req.SessionID, req.StudentID)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit DeleteStudentSession transaction: %w", err)
	}

	return &models.DeleteStudentSessionResponse{
		Status: "Deleted",
	}, nil
}
