package services

import (
	"context"
	"encoding/json"
	"fmt"
	"tracker/app/models"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
)

func sendAssessmentGraderTask(payload models.SQSAssessmentPayload, taskType string) error {
	SQS_LINK := "https://urltosqs"

	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return fmt.Errorf("loading aws config: %w", err)
	}
	payloadParsed, err := json.Marshal(payload)
	if err != nil {
		fmt.Println("Error encoding JSON:", err)
		return nil
	}
	client := sqs.NewFromConfig(cfg)
	_, err = client.SendMessage(context.TODO(), &sqs.SendMessageInput{
		QueueUrl:    aws.String(SQS_LINK),
		MessageBody: aws.String(string(payloadParsed)),
		MessageAttributes: map[string]types.MessageAttributeValue{
			"TaskType": {
				DataType:    aws.String("String"),
				StringValue: aws.String(taskType),
			},
		},
	})
	return err
}
