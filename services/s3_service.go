package services

import (
	"bytes"
	"context"
	"fmt"
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

func (s *AuthService) GeneratePresignedUrl(ctx context.Context, id string) (string, error) {
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

func (s *AuthService) GeneratePutPresignedUrl(ctx context.Context, contentType *string, minMultiplier int64) (*string, *string, error) {
	key := uuid.New()
	keyString := key.String()
	presigner := s3.NewPresignClient(s.s3)
	duration := time.Minute * time.Duration(minMultiplier)
	presignUrl, err := presigner.PresignPutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String("tracker-client-storage"),
		Key:         aws.String(keyString),
		ContentType: aws.String(*contentType),
	}, s3.WithPresignExpires(duration))
	if err != nil {
		return nil, nil, err
	}
	return &presignUrl.URL, &keyString, nil
}

func (s *AuthService) GeneratePutPresignedUrlMaterials(ctx context.Context, key *string, contentType string, minMultiplier int64) (*string, error) {
	presigner := s3.NewPresignClient(s.s3)
	duration := time.Minute * time.Duration(minMultiplier)
	presignUrl, err := presigner.PresignPutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String("tracker-client-storage"),
		Key:         aws.String(*key),
		ContentType: aws.String(contentType),
	}, s3.WithPresignExpires(duration))
	if err != nil {
		return nil, err
	}
	return &presignUrl.URL, nil
}

func (s *AuthService) GetS3Object(ctx context.Context, id string, bucket string) (*string, error) {
	response, err := s.s3.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(id),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get object %w", err)
	}
	defer response.Body.Close()
	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(response.Body)
	if err != nil {
		return nil, fmt.Errorf("unable to write to bytes buffer")
	}
	var res = buf.String()
	var resptr *string
	resptr = &res
	return resptr, nil
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

func (s *AuthService) CreateObjectS3(c context.Context, file multipart.File, keyFound *string, path *string) (*string, error) {
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
		Key:    aws.String(*path + *stringPtr),
		Body:   file,
	})
	if err != nil {
		return nil, err
	}
	return stringPtr, nil
}
