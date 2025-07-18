package config

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// ConnectS3Client initializes and returns an AWS S3 client.
// It loads configuration from the environment, shared config files, or IAM role.
func ConnectS3Client() (*s3.Client, error) {
	region := "us-west-1"

	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(region))
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config: %w", err)
	}

	// Create S3 client from config
	s3Client := s3.NewFromConfig(cfg)
	return s3Client, nil
}
