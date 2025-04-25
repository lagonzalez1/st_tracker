package services

import (
	"context"
	"fmt"
	"strconv"
	"tracker/app/models"
)

func (s *AuthService) ComputeScore(assessment_id *int64, choices map[string]int) (*models.AssessmentScore, error) {

	points, questionEntries, err := s.GradeAssessmentWithCorrectAnswers(assessment_id, choices)
	if err != nil {
		return nil, fmt.Errorf("unable to get assessments by assessment_id: %v", err)
	}

	maxScore, err := s.GetAssessmentMaxScore(assessment_id)
	if err != nil {
		return nil, fmt.Errorf("unable to get assessments by assessment_id: %v", err)
	}

	return &models.AssessmentScore{
		Points:          *points,
		Total:           *maxScore,
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
	INNER JOIN stu_tracker.Questions q ON c.question_id = q.id
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

func (s *AuthService) GradeAssessmentWithCorrectAnswers(assessment_id *int64, choices map[string]int) (*int, []models.AnswerFeedback, error) {
	if assessment_id == nil {
		return nil, nil, fmt.Errorf("assessment is null")
	}

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
		return nil, nil, fmt.Errorf("query failed: %w", err)
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
			return nil, nil, fmt.Errorf("failed to scan row: %w", err)
		}
		questions = append(questions, r)
	}

	var totalScore int
	var allAnswers []models.AnswerFeedback

	for _, q := range questions {
		questionIDstr := strconv.Itoa(int(*q.QuestionID))
		// 3: 7, 4: 10
		if selectedChoiceID, ok := choices[questionIDstr]; ok {
			if q.IsCorrect && selectedChoiceID == int(*q.ChoiceID) {
				totalScore += q.Points
				allAnswers = append(allAnswers, models.AnswerFeedback{
					QuestionID: *q.QuestionID,
					ChoiceID:   *q.ChoiceID,
					IsCorrect:  true,
				})
			} else {
				allAnswers = append(allAnswers, models.AnswerFeedback{
					QuestionID: *q.QuestionID,
					ChoiceID:   *q.ChoiceID,
					IsCorrect:  false,
				})
			}
		}
	}

	return &totalScore, allAnswers, nil
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
