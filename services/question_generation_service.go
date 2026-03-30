package services

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"tracker/app/models"
)

func (s *AuthService) GetQuestionGenerationStatus(ctx context.Context, inputKey *string) (*string, []byte, error) {
	if inputKey == nil || *inputKey == "" {
		return nil, nil, fmt.Errorf("missing parameter input key")
	}

	var status sql.NullString
	var outputKey sql.NullString
	var jsonOutput []byte

	query := `SELECT status, s3_output_key, json_output FROM stu_tracker.Generate_questions_task WHERE input_key = $1;`

	// Need to pass a pointer for jsonOutput to be scanned into
	err := s.db.QueryRowContext(ctx, query, inputKey).Scan(&status, &outputKey, &jsonOutput)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil, fmt.Errorf("no record found for input_key: %s", *inputKey)
		}
		return nil, nil, fmt.Errorf("database error: %w", err)
	}

	// Return actual values (not pointers to sql.Null types)
	statusStr := status.String
	return &statusStr, jsonOutput, nil
}

func (s *AuthService) GetMaterialsGenerationStatus(ctx context.Context, inputKey *string) (*string, *string, []byte, error) {
	if inputKey == nil {
		return nil, nil, nil, fmt.Errorf("missing paramater input key")
	}

	var status sql.NullString
	var outputKey sql.NullString
	var jsonOutput []byte

	query := `SELECT status, s3_output_key, json_output FROM stu_tracker.Generate_materials_task WHERE input_key = $1;`
	err := s.db.QueryRowContext(ctx, query, inputKey).Scan(&status, &outputKey, &jsonOutput)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil, nil, fmt.Errorf("no record found for input_key: %s", *inputKey)
		}
		return nil, nil, nil, fmt.Errorf("database error: %w", err)
	}
	statusStr := status.String
	outputKeyStr := outputKey.String
	return &statusStr, &outputKeyStr, jsonOutput, nil
}

func (s *AuthService) GetStudentReportStatus(ctx context.Context, inputKey *string) (*models.StudentReportResponse, error) {
	if inputKey == nil {
		return nil, fmt.Errorf("input key missing")
	}
	var res models.StudentReportResponse
	var rawReport []byte
	query := `SELECT status, s3_output_key, json_report FROM stu_tracker.Student_report WHERE input_key = $1;`
	// You must list the fields in the same order as your SELECT statement
	err := s.db.QueryRowContext(ctx, query, inputKey).Scan(
		&res.Status,
		&res.OutputKey,
		&rawReport,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Or a custom "Not Found" error
		}
		return nil, err
	}
	if *res.Status == "DONE" && len(rawReport) > 0 {
		var reportData models.StudentReportData
		if err := json.Unmarshal(rawReport, &reportData); err != nil {
			return nil, fmt.Errorf("failed to unmarshal report: %w", err)
		}
		res.JsonReport = &reportData
	}

	return &res, nil
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
