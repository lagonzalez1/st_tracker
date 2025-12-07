package services

import (
	"context"
	"tracker/app/models"
)

func (s *AuthService) AddScheduleV2(c context.Context, req models.RegisterScheduleV2) (*models.RegisterScheduleResponse, error) {
	// Input validation

	return &models.RegisterScheduleResponse{
		Status: "OK",
		ID:     10,
	}, nil
}
