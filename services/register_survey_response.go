package services

import (
	"context"
	"fmt"
	"tracker/app/models"
)

func (s *AuthService) AddSurveyResponse(c context.Context, req models.RegisterSurvey, org *int64) (*models.SimpleReponse, error) {
	// Input validation
	if org == nil || len(req.Survey) == 0 || req.SessionID == nil {
		return nil, fmt.Errorf("missing required fields: orgid and surveys")
	}
	// Begin a transaction
	tx, err := s.db.BeginTx(c, nil)
	if err != nil {
		return nil, fmt.Errorf("unable to begin transaction")
	}

	for i := 0; i < len(req.Survey); i++ {
		var questions = req.Survey[i].Questions
		for k := 0; k < len(questions); k++ {
			var question = req.Survey[i].Questions[k]
			inserQ := `INSERT INTO stu_tracker.Survey_response (
				organization_id, session_id, question_survey_id, response_text, response_choice, question_text
			) VALUES ($1,$2,$3,$4,$5,$6) RETURNING id;`
			_, err := tx.ExecContext(c, inserQ, org, req.SessionID, question.QuestionID, question.ResponseText, question.ResponseChoice, question.QuestionText)
			if err != nil {
				return nil, fmt.Errorf("failed to insert student: %w", err)
			}
		}
	}
	if err = tx.Commit(); err != nil {
		fmt.Println(err)
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return &models.SimpleReponse{
		Status: "OK",
	}, nil
}
