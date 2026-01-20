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

func (s *AuthService) GetMaterialsGenerationStatus(ctx context.Context, inputKey *string) (*string, *string, error) {
	if inputKey == nil {
		return nil, nil, fmt.Errorf("missing paramater input key")
	}
	var status *string
	var outputKey *string
	query := `SELECT status, s3_output_key FROM stu_tracker.Generate_materials_task WHERE input_key = $1;`
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

func (s *AuthService) GetOrganizationReportStatus(ctx context.Context, inputKey *string) (*string, *string, error) {
	if inputKey == nil {
		return nil, nil, fmt.Errorf("missing paramater input key")
	}
	var status *string
	var outputKey *string
	var signedUrl *string
	query := `SELECT status, s3_output_key FROM stu_tracker.Organization_report WHERE input_key = $1;`
	err := s.db.QueryRowContext(ctx, query, inputKey).Scan(&status, &outputKey)
	if err != nil {
		return nil, nil, err
	}

	if *status == "DONE" && outputKey != nil {
		var key = "reports/" + *outputKey + ".csv"
		url, err := s.GeneratePresignedUrl(ctx, key)
		if err != nil {
			return nil, nil, fmt.Errorf("unable to get s3 presigned ulr")
		}
		signedUrl = &url
	}

	return status, signedUrl, nil
}
