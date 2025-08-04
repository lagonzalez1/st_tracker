package services

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"tracker/app/models"

	"github.com/rabbitmq/amqp091-go"
)

func (s *AuthService) AddQueueQuestionEvent(ctx context.Context, req models.RequestQuestions) (*string, error) {
	jsonBody, err := json.Marshal(req)
	if err != nil {
		log.Printf("Failed to marshal JSON for MQ: %v", err)
		return nil, err
	}

	mq := s.mq.Channels
	err = mq["generate"].Publish(
		"ai_events_exchange",
		"generate",
		false, false,
		amqp091.Publishing{
			ContentType: "application/json",
			Body:        []byte(jsonBody),
		},
	)
	if err != nil {
		return nil, err
	}

	var status string = "STARTED"
	var inputKey *string
	query := `INSERT INTO stu_tracker.Generate_questions_task (status, s3_output_key, organization_id) VALUES ($1, $2, $3) RETURNING input_key;`
	err = s.db.QueryRowContext(ctx, query, status, req.S3OutputKey, req.OrganizationID).Scan(&inputKey)
	if err != nil {
		return nil, err
	}

	return inputKey, nil
}

func (s *AuthService) DeleteGeneratedAssessment(ctx context.Context, req models.RemoveGeneratedQuestion) (*models.RemoveResponse, error) {
	query := `SELECT s3_output_key FROM stu_tracker.Generate_questions_task WHERE input_key = $1 AND organization_id = $2;`
	var s3_output_key *string
	err := s.db.QueryRowContext(ctx, query, req.InputKey, req.OrganizationID).Scan(&s3_output_key)
	if err != nil {
		return nil, err
	}

	err = s.DeleteObjectS3(ctx, *s3_output_key)
	if err != nil {
		return nil, err
	}

	deleteQuery := `DELETE FROM stu_tracker.Generate_questions_task WHERE input_key = $1 AND organization_id = $2;`
	res, err := s.db.ExecContext(ctx, deleteQuery, req.InputKey, req.OrganizationID)
	if err != nil {
		return nil, err
	}
	count, err := res.RowsAffected()
	if err != nil {
		return nil, fmt.Errorf("could not get rows affected: %w", err)
	}
	if count == 0 {
		return nil, fmt.Errorf("no task found for input_key=%s org=%d", req.InputKey, req.OrganizationID)
	}
	return &models.RemoveResponse{
		Status: "Deleted",
	}, nil
}
