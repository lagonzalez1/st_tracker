package config

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

// It loads configuration from the environment, shared config files, or IAM role.
func ConnectSQS() (*sqs.Client, error) {
	region := "us-west-1"
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(region))
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config: %w", err)
	}
	client := sqs.NewFromConfig(cfg, func(o *sqs.Options) {
		value, exist := os.LookupEnv("APP_ENV")
		if exist && value == "dev" {
			o.BaseEndpoint = aws.String("http://localstack:4566")
		}
	})
	return client, nil
}
