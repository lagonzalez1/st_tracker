package services

import (
	"context"
	"fmt"
	"strconv"
	"tracker/app/models"
)

func (s *AuthService) ComputeScore(assessment_id *int64, choices map[string]interface{}, grader map[string]bool) (*models.AssessmentScore, error) {
	fmt.Println("Choices interface", choices)
	fmt.Println("Grader interface", grader)
	points, questionEntries, err := s.GradeAssessmentWithCorrectAnswers(assessment_id, choices, grader)
	if err != nil {
		return nil, fmt.Errorf("unable to get assessments by assessment_id: %v", err)
	}

	maxScore, err := s.GetAssessmentMaxScore(assessment_id)
	if err != nil {
		return nil, fmt.Errorf("unable to get assessments by assessment_id: %v", err)
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

func (s *AuthService) GradeAssessmentWithCorrectAnswers(assessment_id *int64, choices interface{}, grader map[string]bool) (*int, []models.AnswerFeedback, error) {
	// Validate inputs
	if assessment_id == nil {
		return nil, nil, fmt.Errorf("assessment_id cannot be nil")
	}
	// Type assert the interfaces to maps
	choicesMap, ok := choices.(map[string]interface{})
	if !ok {
		return nil, nil, fmt.Errorf("choices must be a map[string]interface{}")
	}
	// Fetch correct answers from database
	questions, err := s.fetchCorrectAnswers(*assessment_id)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to fetch correct answers: %w", err)
	}
	// Grade each question
	totalScore := 0
	var allAnswers []models.AnswerFeedback
	fmt.Println("fetchCorrectAnswers", questions)
	for _, q := range questions {
		questionIDstr := strconv.Itoa(int(*q.QuestionID))
		feedback, score, err := s.gradeQuestion(q, questionIDstr, choicesMap, grader)
		fmt.Println("Grade question result: ", feedback, score)
		if err != nil {
			return nil, nil, fmt.Errorf("grading failed for question %s: %w", questionIDstr, err)
		}
		totalScore += score
		allAnswers = append(allAnswers, feedback)
	}

	return &totalScore, allAnswers, nil

}

// Helper function to fetch correct answers from DB
func (s *AuthService) fetchCorrectAnswers(assessment_id int64) ([]models.AssessmentGrader, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	query := `SELECT 
        c.id AS choice_id,
        c.question_id,
        c.is_correct,
        q.points
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
		); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		questions = append(questions, r)
	}

	return questions, nil
}

// Helper function to grade a single question
func (s *AuthService) gradeQuestion(q models.AssessmentGrader, questionIDstr string, choicesMap map[string]interface{}, graderMap map[string]bool) (models.AnswerFeedback, int, error) {
	selectedChoice, exists := choicesMap[questionIDstr]
	fmt.Printf("Type: %T, Value: %v, Exists: %t\n", selectedChoice, selectedChoice, exists)
	if !exists {
		return models.AnswerFeedback{
			QuestionID: *q.QuestionID,
			ChoiceID:   nil,
			IsCorrect:  false,
			AnswerText: nil,
		}, 0, nil
	}
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
			return feedback, q.Points, nil
		}
		return feedback, 0, nil

	case string:
		correctBool := graderMap[questionIDstr]
		feedback := models.AnswerFeedback{
			QuestionID: *q.QuestionID,
			IsCorrect:  correctBool,
			AnswerText: &val,
			ChoiceID:   nil,
		}
		if correctBool {
			return feedback, q.Points, nil
		}
		return feedback, 0, nil
	case nil:
		correctBool := graderMap[questionIDstr]
		feedback := models.AnswerFeedback{
			QuestionID: *q.QuestionID,
			IsCorrect:  correctBool,
			AnswerText: nil,
			ChoiceID:   nil,
		}
		if correctBool {
			return feedback, q.Points, nil
		}
		return feedback, 0, nil

	default:
		selectedChoice, exists := graderMap[questionIDstr]
		if !exists {
			return models.AnswerFeedback{
				QuestionID: *q.QuestionID,
				IsCorrect:  false,
				AnswerText: nil,
				ChoiceID:   nil,
			}, 0, nil
		}
		if exists {
			return models.AnswerFeedback{
				QuestionID: *q.QuestionID,
				IsCorrect:  selectedChoice,
				ChoiceID:   nil,
				AnswerText: nil,
			}, 0, nil
		}

		return models.AnswerFeedback{}, 0, fmt.Errorf("invalid choice type %T for question", selectedChoice)
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

func (s *AuthService) GetAssessmentChoicesByStudent(c context.Context, assessment_id *int64, student_id *int64, session_token *string) ([]models.StudentAssessmentChoices, error) {
	if assessment_id == nil {
		return nil, fmt.Errorf("assessment id is null")
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	query := `SELECT question_id, choice_id, answer_text 
	FROM stu_tracker.Session_answers 
	WHERE session_token = $1 AND student_id = $2 AND assessment_id = $3;`
	rows, err := s.db.QueryContext(ctx, query, session_token, student_id, assessment_id)
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
		)
		if err != nil {
			return nil, err
		}
		assessmentChoices = append(assessmentChoices, choices)
	}
	return assessmentChoices, nil
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
	var existDuplicate bool
	checkDuplicateSubmit := `SELECT EXISTS (SELECT 1 FROM stu_tracker.Session_answers WHERE session_token = $1 AND student_id = $2 AND assessment_id = $3)`
	err = s.db.QueryRowContext(c, checkDuplicateSubmit, req.SessionID, req.StudentID, req.AssessmentID).Scan(&existDuplicate)
	if err != nil {
		return nil, fmt.Errorf("failed to check existing session: %w", err)
	}
	if existDuplicate {
		return nil, fmt.Errorf("no more than one submission per session")
	}

	for questionIDStr, choiceID := range req.Answers {
		switch val := choiceID.(type) {
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
