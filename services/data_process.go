package services

import (
	"context"
	"tracker/app/models"

	"github.com/google/uuid"
)

func (s *AuthService) ProcessStudentReport(ctx context.Context, req models.RequestStudentReport) (*string, error) {
	var status = "STARTED"
	var inputKey *string
	query := `INSERT INTO stu_tracker.Student_report(student_id, semester_id, status, s3_output_key) VALUES ($1,$2, $3,$4) RETURNING input_key;`
	err := s.db.QueryRowContext(ctx, query, req.StudentID, req.SemesterID, status, req.S3OutputKey).Scan(&inputKey)
	if err != nil {
		return nil, err
	}
	return inputKey, nil
}

func (s *AuthService) ProcessDownloadEvent(ctx context.Context, req *models.RequestDownloadData, orgid *int64) (*string, error) {
	var status = "STARTED"
	var inputKeyStr *string
	var s3_output_key = uuid.New().String()
	req.S3OutputKey = &s3_output_key
	query := `INSERT INTO stu_tracker.Organization_report(organization_id, entity, s3_output_key, status) VALUES ($1,$2, $3, $4) RETURNING input_key;`
	err := s.db.QueryRowContext(ctx, query, orgid, req.Entity, s3_output_key, status).Scan(&inputKeyStr)
	if err != nil {
		return nil, err
	}
	return inputKeyStr, nil
}
