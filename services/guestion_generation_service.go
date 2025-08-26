package services

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"tracker/app/models"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
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
	if *status == "DONE" && outputKey != nil {
		out, err := s.s3.GetObject(ctx, &s3.GetObjectInput{
			Bucket: aws.String("client-student-report"),
			Key:    aws.String(*outputKey),
		})
		if err != nil {
			return nil, nil, err
		}
		body, err := io.ReadAll(out.Body)
		if err != nil {
			return nil, nil, err
		}
		var model models.StudentReport
		if err := json.Unmarshal(body, &model); err != nil {
			return nil, nil, err
		}
		fmt.Println(model)
	}

	return status, outputKey, nil
}
