package config

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// It loads configuration from the environment, shared config files, or IAM role.
func ConnectS3() (*s3.Client, error) {
	region := "us-west-1"
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(region))
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config: %w", err)
	}
	value, exist := os.LookupEnv("APP_ENV")
	if exist && value == "dev" {
		client := s3.NewFromConfig(cfg, func(o *s3.Options) {
			o.BaseEndpoint = aws.String("http://localstack:4566")
		})
		return client, nil
	}
	return s3.NewFromConfig(cfg), nil
}
