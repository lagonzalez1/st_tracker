package services

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"tracker/app/models"

	"github.com/google/uuid"
	"github.com/rabbitmq/amqp091-go"
)

func (s *AuthService) AddFileDownloadEvent(ctx context.Context, req models.RequestDownloadData, org_id float64) (*string, error) {
	var S3OutputKey = uuid.New().String()
	req.S3OutputKey = &S3OutputKey

	jsonBody, err := json.Marshal(req)
	if err != nil {
		log.Printf("Failed to marchal json for MQ: %v", err)
		return nil, err
	}
	mq := s.mq.Channels["pgdata"]
	if mq == nil {
		err := fmt.Errorf("RabbitMQ channel 'report' not found")
		log.Printf("Failed : %v", err)
		return nil, err
	}
	err = mq.Publish(
		"ai_events_exchange",
		"pgdata",
		false, false,
		amqp091.Publishing{
			ContentType: "application/json",
			Body:        []byte(jsonBody),
		},
	)
	if err != nil {
		return nil, err
	}
	var status = "STARTED"
	var inputKey *string
	query := `INSERT INTO stu_tracker.Organization_report(organization_id, entity, s3_output_key, status) VALUES ($1,$2, $3, $4) RETURNING input_key;`
	err = s.db.QueryRowContext(ctx, query, org_id, req.Entity, *req.S3OutputKey, status).Scan(&inputKey)
	if err != nil {
		return nil, err
	}
	return inputKey, nil
}

func (s *AuthService) AddGraderEvent(ctx context.Context, req models.RequestAssessmentGrader) (bool, error) {
	jsonBody, err := json.Marshal(req)
	if err != nil {
		log.Printf("Failed to marchal json for MQ: %v", err)
		return false, err
	}
	mq := s.mq.Channels["grader"]
	if mq == nil {
		return false, fmt.Errorf("unable to get channel grader")
	}
	err = mq.Publish(
		"ai_events_exchange",
		"grader",
		false, false,
		amqp091.Publishing{
			ContentType: "application/json",
			Body:        []byte(jsonBody),
		},
	)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (s *AuthService) AddSentimentWorker(ctx context.Context, session_id *int64, org_id *int64) (bool, error) {
	req := map[string]interface{}{"session_id": *session_id, "organization_id": *org_id}
	jsonBody, err := json.Marshal(req)
	if err != nil {
		log.Printf("Failed to marchal json for MQ: %v", err)
		return false, err
	}
	mq := s.mq.Channels["generate"]
	if mq == nil {
		err := fmt.Errorf("RabbitMQ channel 'report' not found")
		log.Printf("Failed : %v", err)
		return false, err
	}
	err = mq.Publish(
		"ai_events_exchange",
		"generate",
		false, false,
		amqp091.Publishing{
			ContentType: "application/json",
			Body:        []byte(jsonBody),
		},
	)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (s *AuthService) AddStudentReportQuery(ctx context.Context, req models.RequestStudentReport) (*string, error) {
	jsonBody, err := json.Marshal(req)
	if err != nil {
		log.Printf("Failed to marchal json for MQ: %v", err)
		return nil, err
	}
	mq := s.mq.Channels["report"]
	if mq == nil {
		err := fmt.Errorf("RabbitMQ channel 'report' not found")
		log.Printf("Failed : %v", err)
		return nil, err
	}
	err = mq.Publish(
		"ai_events_exchange",
		"report",
		false, false,
		amqp091.Publishing{
			ContentType: "application/json",
			Body:        []byte(jsonBody),
		},
	)
	if err != nil {
		return nil, err
	}
	log.Println("Success: Message published to RabbitMQ")

	var status = "STARTED"
	var inputKey *string
	query := `INSERT INTO stu_tracker.Student_report(student_id, semester_id, status, s3_output_key) VALUES ($1,$2, $3,$4) RETURNING input_key;`

	err = s.db.QueryRowContext(ctx, query, req.StudentID, req.SemesterID, status, req.S3OutputKey).Scan(&inputKey)
	if err != nil {
		return nil, err
	}
	return inputKey, nil
}

func (s *AuthService) AddQueueMaterialsEvent(ctx context.Context, req *models.RequestEventGeneration) (*string, error) {
	jsonBody, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	mq := s.mq.Channels["produce"]
	if mq == nil {
		return nil, fmt.Errorf("channel is not open")
	}
	err = mq.Publish(
		"ai_events_exchange",
		"produce",
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
	var intputKey *string
	query := `INSERT INTO stu_tracker.Generate_materials_task(status, s3_output_key, organization_id, assessment_id) VALUES ($1,$2, $3,$4) RETURNING input_key;`
	err = s.db.QueryRowContext(ctx, query, status, req.RequestMaterials.S3OutputKey, req.OrganizationID, req.RequestMaterials.AssessmentId).Scan(&intputKey)
	if err != nil {
		return nil, err
	}
	return intputKey, nil
}

func (s *AuthService) AddQueueQuestionEvent(ctx context.Context, req *models.RequestEventGeneration) (*string, error) {
	jsonBody, err := json.Marshal(req)
	if err != nil {
		log.Printf("Failed to marshal JSON for MQ: %v", err)
		return nil, err
	}

	mq := s.mq.Channels["produce"]
	if mq == nil {
		err := fmt.Errorf("RabbitMQ channel 'generate' not found")
		log.Printf("Failure: %v", err)
		return nil, err
	}
	log.Printf("Attempting to publish message to exchange 'ai_events_exchange'...")

	err = mq.Publish(
		"ai_events_exchange",
		"produce",
		false, false,
		amqp091.Publishing{
			ContentType: "application/json",
			Body:        []byte(jsonBody),
		},
	)
	if err != nil {
		log.Printf("Failure: Failed to publish message to RabbitMQ: %v", err)
		return nil, err
	}
	log.Println("Success: Message published to RabbitMQ.")
	var status string = "STARTED"
	var inputKey *string
	query := `INSERT INTO stu_tracker.Generate_questions_task (status, s3_output_key, organization_id) VALUES ($1, $2, $3) RETURNING input_key;`
	err = s.db.QueryRowContext(ctx, query, status, req.RequestQuestions.S3OutputKey, req.OrganizationID).Scan(&inputKey)
	if err != nil {
		log.Printf("Failure: Failed to insert record into database: %v", err)
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
