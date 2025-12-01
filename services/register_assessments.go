package services

import (
	"context"
	"fmt"
	"strings"
	"tracker/app/models"

	"github.com/lib/pq"
)

/*
id SERIAL PRIMARY KEY,
    assessment_id INT REFERENCES stu_tracker.Assessments(id) ON DELETE CASCADE,
    image_url TEXT,
    question_text TEXT NOT NULL,
    question_type VARCHAR(50) NOT NULL, -- e.g., 'multiple_choice', 'true_false', 'short_answer'
    points INT DEFAULT 1,
    order_number INT, -- optional, if you want to sort questions
    is_required BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP


CREATE TABLE stu_tracker.Choices (
    id SERIAL PRIMARY KEY,
    question_id INT REFERENCES stu_tracker.Questions(id) ON DELETE CASCADE,
    choice_text TEXT NOT NULL,
    is_correct BOOLEAN DEFAULT FALSE,
    order_number INT, -- useful for displaying options in order
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
*/

func (s *AuthService) AddAssessment(c context.Context, req models.RegisterAssessment) (*models.RegisterResponseAssessment, error) {
	// Input validation
	if req.Title == "" || req.MaxScore == nil || req.OrganizationID == nil {
		return nil, fmt.Errorf("missing required fields: first_name, last_name, email, password")
	}
	var newID int64
	query := `INSERT INTO stu_tracker.Assessments 
			(title, description, letter, cycle, alpha_identifier, external_link, max_score, subject_id, material_id, 
			organization_id, visible, program_id, version, pre, post, mid, easy_score, grade_level, questionnaire) 
			  VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19) RETURNING id;`
	err := s.db.QueryRowContext(c, query, req.Title, req.Description, req.Letter,
		req.Cycle, req.AlphaIdentifier, req.ExternalLink, req.MaxScore,
		req.SubjectId, req.MaterialID, req.OrganizationID, req.Visible,
		req.ProgramId, req.Version, req.Pre, req.Post, req.Mid, req.EasyScore, req.GradeLevel, req.Questionnaire).Scan(&newID)
	if err != nil {
		return nil, fmt.Errorf("failed to insert student: %w", err)
	}
	if len(req.Questions) > 0 && req.EasyScore {
		for idx, question := range req.Questions {
			var questionID int64
			var questionImageID = &req.Images[idx].ID
			// Insert the question
			questionQuery := `INSERT INTO stu_tracker.Questions 
				(assessment_id, image_url, question_text, question_type, points, order_number, standard_text)
				VALUES ($1, $2, $3, $4, $5, $6, $7)
				RETURNING id;`
			fixedText := strings.ReplaceAll(question.QuestionText, "\\\\", "\\")
			err := s.db.QueryRowContext(
				c,
				questionQuery,
				newID,
				questionImageID,
				fixedText,
				question.QuestionType,
				question.Points,
				question.OrderNumber,
				question.Standard,
			).Scan(&questionID)
			if err != nil {
				return nil, fmt.Errorf("failed to insert question: %w", err)
			}

			// Insert choices if any
			for _, choice := range question.Choice {
				choiceQuery := ` INSERT INTO stu_tracker.Choices 
					(question_id, choice_text, is_correct, order_number)
					VALUES ($1, $2, $3, $4);`
				_, err := s.db.ExecContext(
					c,
					choiceQuery,
					questionID,
					choice.ChoiceText,
					choice.IsCorrect,
					choice.OrderNumber,
				)
				if err != nil {
					return nil, fmt.Errorf("failed to insert choice for question %d: %w", questionID, err)
				}
			}
			if len(req.Questions) != len(req.Images) {
				return nil, fmt.Errorf("Something went wrong inserting images. Created assessment OK.", len(req.Questions), len(req.Images))
			}
		}
	}

	return &models.RegisterResponseAssessment{
		Status:       "OK",
		AssessmentId: newID,
	}, nil
}

func (s *AuthService) UpdateAssessment(c context.Context, req models.RegisterAssessment) (*models.RegisterResponseAssessment, error) {
	// Input validation
	if req.Title == "" || req.MaxScore == nil || req.OrganizationID == nil {
		return nil, fmt.Errorf("missing required fields: first_name, last_name, email, password")
	}
	query := `UPDATE stu_tracker.Assessments SET title = $1, description = $2, letter = $3, 
	cycle = $4, alpha_identifier = $5, external_link = $6, max_score = $7, subject_id = $8, 
	material_id = $9, organization_id = $10, visible = $11, program_id = $12, version = $13, pre = $14, post = $15, mid = $16, easy_score = $17,
	grade_level = $18, questionnaire = $19
	WHERE id = $20;`
	_, err := s.db.ExecContext(c, query, req.Title, req.Description, req.Letter, req.Cycle,
		req.AlphaIdentifier, req.ExternalLink, req.MaxScore, req.SubjectId, req.MaterialID,
		req.OrganizationID, req.Visible, req.ProgramId, req.Version, req.Pre, req.Post, req.Mid,
		req.EasyScore, req.GradeLevel, req.Questionnaire, req.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to insert student: %w", err)
	}

	// Update any questions by adding or updating current questions in place.
	if req.EasyScore && len(req.Questions) > 0 {
		for idx, question := range req.Questions {
			// Question is already on file update..
			var questionImageID = &req.Images[idx].ID
			if question.QuestionID != nil {
				// Update the question
				updateQuery := `UPDATE stu_tracker.Questions 
					SET assessment_id = $1,
						image_url = $2,
						question_text = $3,
						question_type = $4,
						points = $5,
						order_number = $6,
						standard_text = $7
					WHERE id = $8;`

				_, err := s.db.ExecContext(
					c,
					updateQuery,
					req.ID,
					questionImageID,
					question.QuestionText,
					question.QuestionType,
					question.Points,
					question.OrderNumber,
					question.Standard,
					question.QuestionID,
				)
				if err != nil {
					return nil, fmt.Errorf("failed to update question: %w", err)
				}
				for _, choice := range question.Choice {
					if choice.ChoiceID != nil {
						// Update existing choice
						updateChoiceQuery := `UPDATE stu_tracker.Choices
							SET choice_text = $1,
								is_correct = $2,
								order_number = $3
							WHERE id = $4;`

						_, err := s.db.ExecContext(
							c,
							updateChoiceQuery,
							choice.ChoiceText,
							choice.IsCorrect,
							choice.OrderNumber,
							choice.ChoiceID,
						)
						if err != nil {
							return nil, fmt.Errorf("failed to update choice ID %d: %w", *choice.ChoiceID, err)
						}
					} else {
						// Insert new choice
						insertChoiceQuery := `INSERT INTO stu_tracker.Choices 
							(question_id, choice_text, is_correct, order_number)
							VALUES ($1, $2, $3, $4);`

						_, err := s.db.ExecContext(
							c,
							insertChoiceQuery,
							question.QuestionID,
							choice.ChoiceText,
							choice.IsCorrect,
							choice.OrderNumber,
						)
						if err != nil {
							return nil, fmt.Errorf("failed to insert new choice for question %d: %w", *question.QuestionID, err)
						}
					}
				}

			} else {
				var questionID int64
				// Insert the question
				questionQuery := `INSERT INTO stu_tracker.Questions 
				(assessment_id, image_url, question_text, question_type, points, order_number, standard_text)
				VALUES ($1, $2, $3, $4, $5, $6, $7)
				RETURNING id;`
				var questionImageID = &req.Images[idx].ID
				fixedText := strings.ReplaceAll(question.QuestionText, "\\\\", "\\")

				err := s.db.QueryRowContext(
					c,
					questionQuery,
					req.ID,
					questionImageID,
					fixedText,
					question.QuestionType,
					question.Points,
					question.OrderNumber,
					question.Standard,
				).Scan(&questionID)
				if err != nil {
					return nil, fmt.Errorf("failed to insert question: %w", err)
				}
				// Insert choices if any
				for _, choice := range question.Choice {
					choiceQuery := ` INSERT INTO stu_tracker.Choices 
					(question_id, choice_text, is_correct, order_number)
					VALUES ($1, $2, $3, $4);`
					_, err := s.db.ExecContext(
						c,
						choiceQuery,
						questionID,
						choice.ChoiceText,
						choice.IsCorrect,
						choice.OrderNumber,
					)
					if err != nil {
						return nil, fmt.Errorf("failed to insert choice for question %d: %w", questionID, err)
					}
				}
			}
		}
	}
	// Remove any questions from array of ids
	if len(req.RemoveQuestions) > 0 {
		query := `DELETE FROM stu_tracker.Questions WHERE id = ANY($1)`
		_, err := s.db.ExecContext(c, query, pq.Array(req.RemoveQuestions))
		if err != nil {
			return nil, fmt.Errorf("failed to delete questions: %w", err)
		}
	}

	return &models.RegisterResponseAssessment{
		Status:       "OK",
		AssessmentId: *req.ID,
	}, nil
}

func (s *AuthService) DeleteAssessment(c context.Context, req models.RemoveRequest) (*models.RemoveResponse, error) {
	// Input validation
	if req.ID == nil {
		return nil, fmt.Errorf("missing required fields: id")
	}
	query := `DELETE FROM stu_tracker.Assessments WHERE id = $1`
	_, err := s.db.ExecContext(c, query, *req.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to insert Locations: %w", err)
	}
	return &models.RemoveResponse{
		Status: "Removed",
	}, nil
}
