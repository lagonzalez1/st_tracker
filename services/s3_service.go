package services

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func (s *AuthService) GenerateMaterialsPresignedUrl(ctx context.Context, id string) (string, error) {
	presigner := s3.NewPresignClient(s.s3)
	presignParams := &s3.GetObjectInput{
		Bucket:                     aws.String("tracker-client-storage"),
		Key:                        aws.String(id),
		ResponseContentType:        aws.String("application/pdf"),
		ResponseContentDisposition: aws.String("inline"),
	}
	presignOpts := func(po *s3.PresignOptions) {
		po.Expires = 5 * time.Minute
	}
	presignResult, err := presigner.PresignGetObject(context.TODO(), presignParams, presignOpts)
	if err != nil {
		return "nil", err
	}
	return presignResult.URL, nil
}

func (s *AuthService) GenerateAssessmentsPresignedUrl(ctx context.Context, id string) (string, error) {
	presigner := s3.NewPresignClient(s.s3)
	presignParams := &s3.GetObjectInput{
		Bucket: aws.String("tracker-client-storage"),
		Key:    aws.String(id),
	}
	presignOpts := func(po *s3.PresignOptions) {
		po.Expires = 5 * time.Minute
	}
	presignResult, err := presigner.PresignGetObject(context.TODO(), presignParams, presignOpts)
	if err != nil {
		return "nil", err
	}
	return presignResult.URL, nil
}

func (s *AuthService) DeleteObjectS3(ctx context.Context, id string) error {

	presignParams := &s3.DeleteObjectInput{
		Bucket: aws.String("tracker-client-storage"),
		Key:    aws.String(id),
	}
	_, err := s.s3.DeleteObject(ctx, presignParams)
	if err != nil {
		return err
	}
	return nil
}
