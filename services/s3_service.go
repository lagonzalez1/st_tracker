package services

import (
	"context"
	"mime/multipart"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/uuid"
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

func (s *AuthService) CreateObjectS3(c context.Context, file multipart.File, keyFound *string) (*string, error) {
	key := uuid.New()
	keyString := key.String()
	var stringPtr *string
	// If file exist update such key.
	if keyFound == nil {
		stringPtr = &keyString
	} else {
		stringPtr = keyFound
	}
	_, err := s.s3.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String("tracker-client-storage"),
		Key:    stringPtr,
		Body:   file,
	})
	if err != nil {
		return nil, err
	}
	return stringPtr, nil
}
