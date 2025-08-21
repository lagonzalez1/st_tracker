package services

import (
	"context"
	"fmt"
)

func (s *AuthService) GetQuestionGenerationStatus(ctx context.Context, inputKey *string) (*string, *string, error) {
	if inputKey == nil {
		return nil, nil, fmt.Errorf("missing paramater input key")
	}
	var status *string
	var outputKey *string
	query := `SELECT status, s3_output_key FROM stu_tracker.Generate_questions_task WHERE input_key = $1;`
	err := s.db.QueryRowContext(ctx, query, inputKey).Scan(&status, &outputKey)
	if err != nil {
		return nil, nil, err
	}
	return status, outputKey, nil
}

func (s *AuthService) GetStudentReportStatus(ctx context.Context, inputKey *string) (*string, *string, error) {
	if inputKey == nil {
		return nil, nil, fmt.Errorf("missing paramater input key")
	}
	var status *string
	var outputKey *string
	query := `SELECT status, s3_output_key FROM stu_tracker.Student_report WHERE input_key = $1;`
	err := s.db.QueryRowContext(ctx, query, inputKey).Scan(&status, &outputKey)
	if err != nil {
		return nil, nil, err
	}
	return status, outputKey, nil
}
