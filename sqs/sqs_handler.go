package sqs

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"strings"
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

func (s *SqsHandler) TagPayloadAssessmentGrader(ctx context.Context, key string, payload *models.RequestAssessmentGrader) ([]byte, error) {
	if err := ctx.Err(); err != nil {
		return nil, fmt.Errorf("context cancelled: %w", err)
	}

	if payload == nil {
		return nil, fmt.Errorf("payload cannot be nil")
	}

	data := struct {
		Task string                         `json:"task"`
		Body models.RequestAssessmentGrader `json:"body"`
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

func (s *SqsHandler) TagPayloadAssessmentGenerator(ctx context.Context, key string, payload *models.RequestEventGeneration) ([]byte, error) {
	if err := ctx.Err(); err != nil {
		return nil, fmt.Errorf("context cancelled: %w", err)
	}

	if payload == nil {
		return nil, fmt.Errorf("payload cannot be nil")
	}

	data := struct {
		Task string                        `json:"task"`
		Body models.RequestEventGeneration `json:"body"`
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

func (sh *SqsHandler) SendMessageToQueue(ctx context.Context, queueName string, messageBody string) (*sqs.SendMessageOutput, error) {
	fmt.Printf("[DEBUG] Queue name received: '%s'\n", queueName)
	fmt.Printf("[DEBUG] Queue name length: %d\n", len(queueName))

	urlOut, err := sh.sqs.GetQueueUrl(ctx, &sqs.GetQueueUrlInput{
		QueueName: aws.String(queueName),
	})
	if err != nil {
		fmt.Printf("[ERROR] GetQueueUrl failed: %v\n", err)
		return nil, fmt.Errorf("Error found getting url : %w", err)
	}

	fmt.Printf("[DEBUG] Got queue URL: %s\n", *urlOut.QueueUrl)

	resp, err := sh.sqs.SendMessage(ctx, &sqs.SendMessageInput{
		MessageBody: aws.String(messageBody),
		QueueUrl:    urlOut.QueueUrl,
	})

	if err != nil {
		fmt.Printf("[ERROR] SendMessage failed: %v\n", err)
		return nil, err
	}

	fmt.Printf("[SUCCESS] Message sent: %s\n", *resp.MessageId)
	return resp, nil
}

func (sh *SqsHandler) SendMessageToFIFOQueue(ctx context.Context, queueName string, messageBody string, orgid int64) (*sqs.SendMessageOutput, error) {
	fmt.Printf("[DEBUG] Queue name received: '%s'\n", queueName)
	fmt.Printf("[DEBUG] Queue name length: %d\n", len(queueName))

	urlOut, err := sh.sqs.GetQueueUrl(ctx, &sqs.GetQueueUrlInput{
		QueueName: aws.String(queueName),
	})
	if err != nil {
		fmt.Printf("[ERROR] GetQueueUrl failed: %v\n", err)
		return nil, fmt.Errorf("Error found getting url : %w", err)
	}

	input := &sqs.SendMessageInput{
		MessageBody: aws.String(messageBody),
		QueueUrl:    urlOut.QueueUrl,
	}

	if strings.HasSuffix(queueName, ".fifo") {
		messageId := fmt.Sprintf("message:queue:%d", orgid)
		input.MessageGroupId = aws.String(messageId)
		input.MessageDeduplicationId = aws.String(generateDeduplicationId(messageBody))
	}

	fmt.Printf("[DEBUG] Got queue URL: %s\n", *urlOut.QueueUrl)

	resp, err := sh.sqs.SendMessage(ctx, input)

	if err != nil {
		fmt.Printf("[ERROR] SendMessage failed: %v\n", err)
		return nil, err
	}

	fmt.Printf("[SUCCESS] Message sent: %s\n", *resp.MessageId)
	return resp, nil
}

func generateDeduplicationId(messageBody string) string {
	hash := sha256.Sum256([]byte(messageBody))
	return fmt.Sprintf("%x", hash)
}
