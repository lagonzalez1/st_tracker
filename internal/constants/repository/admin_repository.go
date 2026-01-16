package repository

import (
	"context"
	"tracker/app/models"
)

type AdminRepository interface {
	Create(ctx context.Context, admin models.Admin) (*models.ResponseRequestAdmin, error)
	Update(ctx context.Context, admin models.Admin) (*models.ResponseRequestAdmin, error)
	Delete(ctx context.Context, id models.RemoveRequest) (*models.RemoveResponse, error)
}
