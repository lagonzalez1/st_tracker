package sqs

import (
	"context"
	"encoding/json"
	"fmt"
	"tracker/app/models"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

type SqsHandler struct {
	sqs *sqs.Client
}

func New(sqs *sqs.Client) *SqsHandler {
	return &SqsHandler{sqs: sqs}
}

func (s *SqsHandler) TagPayloadStudentReport(ctx context.Context, key string, payload *models.RequestStudentReport) ([]byte, error) {
	// Check if context is cancelled
	if err := ctx.Err(); err != nil {
		return nil, fmt.Errorf("context cancelled: %w", err)
	}

	if payload == nil {
		return nil, fmt.Errorf("payload cannot be nil")
	}

	data := struct {
		Task string                      `json:"task"`
		Body models.RequestStudentReport `json:"body"`
	}{
		Task: key,
		Body: *payload,
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payload: %w", err)
	}

	return jsonData, nil
}

func (s *SqsHandler) TagPayloadTutorDownload(ctx context.Context, key string, payload *models.RequestDownloadData) ([]byte, error) {
	if err := ctx.Err(); err != nil {
		return nil, fmt.Errorf("context cancelled: %w", err)
	}

	if payload == nil {
		return nil, fmt.Errorf("payload cannot be nil")
	}

	data := struct {
		Task string                     `json:"task"`
		Body models.RequestDownloadData `json:"body"`
	}{
		Task: key,
		Body: *payload,
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payload: %w", err)
	}

	return jsonData, nil
}

func (sh *SqsHandler) SendMessageToQueue(ctx context.Context, url *string, messageBody string) (*sqs.SendMessageOutput, error) {
	return sh.sqs.SendMessage(ctx, &sqs.SendMessageInput{
		MessageBody: aws.String(messageBody),
		QueueUrl:    url,
	})
}
