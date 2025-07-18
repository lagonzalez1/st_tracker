package services

import (
	"context"
	"fmt"
	"mime/multipart"
	"tracker/app/models"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/uuid"
)

func (s *AuthService) UploadMaterialFile(c context.Context, file multipart.File, keyFound *string) (*string, error) {
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

func (s *AuthService) AddMaterial(c context.Context, req models.RegisterRequestMaterials) (*models.ResponseRequestMaterials, error) {
	// Input validation
	if req.Title == "" || req.OrganizationId == nil {
		return nil, fmt.Errorf("missing required fields: first_name, last_name, or email")
	}
	var newID int64
	query := `INSERT INTO stu_tracker.Materials
			(title, external_link, description, pre, mid, post, visible, version, organization_id, location_id, program_id, s3_reference)
              VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12) RETURNING id;`

	err := s.db.QueryRowContext(c, query, req.Title, req.ExternalLink, req.Description,
		req.Pre, req.Mid, req.Post, req.Visible, req.Version, *req.OrganizationId, req.LocationId, req.ProgramId, req.SReference).Scan(&newID)
	if err != nil {
		return nil, fmt.Errorf("failed to insert Locations: %w", err)
	}
	return &models.ResponseRequestMaterials{
		Status:     "OK",
		MaterialId: newID,
	}, nil
}

func (s *AuthService) UpdateMaterial(c context.Context, req models.RegisterRequestMaterials) (*models.ResponseUpdate, error) {
	// Input validation
	if req.ID == nil {
		return nil, fmt.Errorf("missing required fields: id")
	}

	query := `UPDATE stu_tracker.Materials SET
			  title = $1, external_link = $2, description = $3, pre = $4, mid = $5, 
			  post = $6, visible = $7, version = $8, organization_id = $9, location_id = $10, program_id = $11, s3_reference = $12
              WHERE id = $13`

	_, err := s.db.ExecContext(c, query, req.Title, req.ExternalLink, req.Description,
		req.Pre, req.Mid, req.Post, req.Visible, req.Version, *req.OrganizationId, req.LocationId, req.ProgramId, req.SReference, *req.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to insert Locations: %w", err)
	}
	return &models.ResponseUpdate{
		Status: "Updated",
	}, nil
}

func (s *AuthService) DoesReferenceExist(c context.Context, id *int64) (*string, error) {
	if id == nil {
		return nil, fmt.Errorf("missing required fields: id")
	}
	var reference *string
	selectQuery := `SELECT s3_reference FROM stu_tracker.Materials WHERE id = $1;`
	err := s.db.QueryRowContext(c, selectQuery, id).Scan(&reference)
	if err != nil {
		return nil, err
	}
	return reference, nil
}

func (s *AuthService) DeleteMaterial(c context.Context, req models.RemoveRequest) (*models.RemoveResponse, error) {
	if req.ID == nil {
		return nil, fmt.Errorf("missing required fields: id")
	}
	// Reference variable from DB
	var reference *string
	selectQuery := `SELECT s3_reference FROM stu_tracker.Materials WHERE id = $1;`
	err := s.db.QueryRowContext(c, selectQuery, req.ID).Scan(&reference)
	if err != nil {
		return nil, err
	}
	if reference != nil {
		// Delete from s3 given the correct reference to s3 key
		err = s.DeleteObjectS3(c, *reference)
		if err != nil {
			return nil, err
		}
	}
	query := `DELETE FROM stu_tracker.Materials WHERE id = $1`
	_, err = s.db.ExecContext(c, query, *req.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to insert Locations: %w", err)
	}
	return &models.RemoveResponse{
		Status: "Removed",
	}, nil

}
