package services

import (
	"context"
	"fmt"
	"tracker/app/models"
)

func (s *AuthService) AddSurvey(c context.Context, req models.RegisterRequestSurvey, org *int64) (*models.ResponseRequestSurvey, error) {
	// Input validation
	if *req.Title == "" || len(req.Questions) == 0 {
		return nil, fmt.Errorf("missing required fields: ID, Program name")
	}
	// Begin a transaction
	tx, err := s.db.BeginTx(c, nil)
	if err != nil {
		return nil, fmt.Errorf("unable to begin transaction")
	}
	var newID int64
	query := `INSERT INTO stu_tracker.Surveys(organization_id, title, description, order_by, is_active) 
			  VALUES ($1, $2, $3, $4, $5) RETURNING id;`
	err = tx.QueryRowContext(c, query, org, req.Title, req.Description, req.OrderBy, req.IsActive).Scan(&newID)
	if err != nil {
		return nil, fmt.Errorf("failed to insert student: %w", err)
	}

	for i := 0; i < len(req.Questions); i++ {
		var sid *int64
		inserQ := `INSERT INTO stu_tracker.survey_questions(
			survey_id, order_index, question_text, question_type
		) VALUES ($1,$2,$3,$4) RETURNING id;`
		err := tx.QueryRowContext(c, inserQ, newID, i+1, req.Questions[i].QuestionText, req.Questions[i].QuestionType).Scan(&sid)
		if err != nil {
			return nil, fmt.Errorf("failed to insert student: %w", err)
		}
		if *req.Questions[i].QuestionType == "multi_choice" || *req.Questions[i].QuestionType == "yes_no" {
			insertQ := `INSERT INTO stu_tracker.survey_choice (question_survey_id, choice_text) VALUES ($1,$2);`
			for j := 0; j < len(req.Questions[i].Choices); j++ {
				var choice = req.Questions[i].Choices[j]
				_, err = tx.ExecContext(c, insertQ, sid, choice.ChoiceText)
				if err != nil {
					return nil, fmt.Errorf("failed to insert student: %w", err)
				}
			}

		}
	}
	if err = tx.Commit(); err != nil {
		fmt.Println(err)
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return &models.ResponseRequestSurvey{
		Status:    "OK",
		ProgramId: newID,
	}, nil
}

func (s *AuthService) UpdateSurvey(c context.Context, req models.RegisterRequestSurvey, orgid *int64) (*models.RemoveResponse, error) {

	tx, err := s.db.BeginTx(c, nil)
	if err != nil {
		return nil, fmt.Errorf("unable to beign transaction")
	}

	query := `UPDATE stu_tracker.Surveys SET title = $1, description = $2, is_active = $3, order_by = $4 WHERE id = $5 AND organization_id = $6`
	_, err = tx.ExecContext(c, query, req.Title, req.Description, req.IsActive, req.OrderBy, req.ID, orgid)
	if err != nil {
		return nil, fmt.Errorf("unable to update program: %w", err)
	}

	if len(req.RemoveQuestions) > 0 {
		remove := `DELETE FROM stu_tracker.survey_questions WHERE id = $1 AND organization_id = $2`
		for i := 0; i < len(req.RemoveQuestions); i++ {
			_, err := tx.ExecContext(c, remove, req.RemoveQuestions[i], *orgid)
			if err != nil {
				return nil, fmt.Errorf("unable to update questions")
			}
		}
	}

	if len(req.Questions) > 0 {
		// psudo update/
		for i := 0; i < len(req.Questions); i++ {
			if req.Questions[i].ID != nil {
				// Update
				update := `UPDATE stu_tracker.survey_questions SET question_text = $1, question_type = $2 WHERE id = $3;`
				_, err := tx.ExecContext(c, update, req.Questions[i].QuestionText, req.Questions[i].QuestionType, req.Questions[i].ID)
				if err != nil {
					return nil, fmt.Errorf("unable to update questions")
				}
				if *req.Questions[i].QuestionType == "multi_choice" || *req.Questions[i].QuestionType == "yes_no" {
					// If there exist already some choices remnove before hand
					if len(req.Questions[i].Choices) > 0 {
						remove := `DELETE FROM stu_tracker.survey_choice WHERE question_survey_id = $1`
						_, err := tx.ExecContext(c, remove, req.Questions[i].ID)
						if err != nil {
							return nil, fmt.Errorf("unable to delete survey choices")
						}
					}
					// Otherwise create the choices
					for k := 0; k < len(req.Questions[i].Choices); k++ {
						insertQ := `INSERT INTO stu_tracker.survey_choice (question_survey_id, choice_text) VALUES ($1,$2);`
						_, err := tx.ExecContext(c, insertQ, req.Questions[i].ID, req.Questions[i].Choices[k].ChoiceText)
						if err != nil {
							return nil, fmt.Errorf("unable to delete survey choices")
						}
					}
				}
			} else {
				var sid *int64
				inserQ := `INSERT INTO stu_tracker.survey_questions(
						survey_id, order_index, question_text, question_type
					) VALUES ($1,$2,$3,$4) RETURNING id;`
				err := tx.QueryRowContext(c, inserQ, req.ID, i+1, req.Questions[i].QuestionText, req.Questions[i].QuestionType).Scan(&sid)
				if err != nil {
					return nil, fmt.Errorf("failed to insert student: %w", err)
				}
				if *req.Questions[i].QuestionType == "multi_choice" || *req.Questions[i].QuestionType == "yes_no" {
					insertQ := `INSERT INTO stu_tracker.survey_choice (question_survey_id, choice_text) VALUES ($1,$2);`
					for j := 0; j < len(req.Questions[i].Choices); j++ {
						var choice = req.Questions[i].Choices[j]
						_, err = tx.ExecContext(c, insertQ, sid, choice.ChoiceText)
						if err != nil {
							return nil, fmt.Errorf("failed to insert student: %w", err)
						}
					}

				}

			}

		}
	}
	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("unable to commit all changes")
	}
	return &models.RemoveResponse{
		Status: "Updated",
	}, nil
}

func (s *AuthService) DeleteSurvey(c context.Context, req models.RemoveRequest) (*models.RemoveResponse, error) {
	// Input validation
	if req.ID == nil {
		return nil, fmt.Errorf("missing required fields: ID")
	}
	query := `DELETE FROM stu_tracker.Surveys WHERE id = $1;`
	_, err := s.db.Exec(query, req.ID)
	if err != nil {
		return nil, fmt.Errorf("unable to delete program: %w", err)
	}
	return &models.RemoveResponse{
		Status: "Removed",
	}, nil
}
