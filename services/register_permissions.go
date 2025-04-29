package services

import (
	"context"
	"database/sql"
	"fmt"
	"tracker/app/models"

	"github.com/lib/pq"
)

func (s *AuthService) CreatePermission(ctx context.Context, req models.RegisterPermissionRequest) (*models.RegisterPermissionResponse, error) {
	// Validate input
	if req.OrganizationId == nil || req.Role == "" || req.ID == nil {
		return nil, fmt.Errorf("missing required fields: organization ID, role, or user ID")
	}

	// Begin transaction
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Step 1: Get all valid permissions from both arrays
	validPermissions, err := s.getValidPermissions(ctx, tx, req.Permissions, req.UpdatePermissions)
	if err != nil {
		return nil, fmt.Errorf("failed to validate permissions: %w", err)
	}

	if len(validPermissions) == 0 {
		return &models.RegisterPermissionResponse{Status: "No valid permissions provided"}, nil
	}

	// Step 2: Delete permissions not in the update list (maintain referential integrity)
	if err := s.cleanStalePermissions(ctx, tx, req, validPermissions); err != nil {
		return nil, fmt.Errorf("failed to clean stale permissions: %w", err)
	}

	// Step 3: Insert new permissions
	if err := s.insertNewPermissions(ctx, tx, req, validPermissions); err != nil {
		return nil, fmt.Errorf("failed to insert new permissions: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return &models.RegisterPermissionResponse{
		Status: "OK",
	}, nil
}

// Helper function to validate permissions against database
func (s *AuthService) getValidPermissions(ctx context.Context, tx *sql.Tx, newPerms []string, updatePerms []string) ([]models.PermissionsMap, error) {
	query := `
        SELECT p.id, p.name 
        FROM stu_tracker.Permissions p 
        WHERE p.name = ANY($1) OR p.name = ANY($2)`

	rows, err := tx.QueryContext(ctx, query, pq.Array(newPerms), pq.Array(updatePerms))
	if err != nil {
		return nil, fmt.Errorf("permission query failed: %w", err)
	}
	defer rows.Close()

	var permissions []models.PermissionsMap
	for rows.Next() {
		var p models.PermissionsMap
		if err := rows.Scan(&p.ID, &p.PermissionName); err != nil {
			return nil, fmt.Errorf("permission scan failed: %w", err)
		}
		permissions = append(permissions, p)
	}
	return permissions, rows.Err()
}

// Helper function to remove stale permissions
func (s *AuthService) cleanStalePermissions(ctx context.Context, tx *sql.Tx, req models.RegisterPermissionRequest, validPerms []models.PermissionsMap) error {
	var deleteQuery string
	switch req.User {
	case "TUTOR":
		deleteQuery = `DELETE FROM stu_tracker.Tutor_Permissions 
                      WHERE tutor_id = $1 AND permission_id <> ALL($2)`
	case "ADMIN":
		deleteQuery = `DELETE FROM stu_tracker.Admin_Permissions 
                      WHERE admin_id = $1 AND permission_id <> ALL($2)`
	default:
		return fmt.Errorf("invalid user type")
	}

	permIDs := make([]int64, len(validPerms))
	for i, p := range validPerms {
		permIDs[i] = *p.ID
	}

	_, err := tx.ExecContext(ctx, deleteQuery, *req.ID, pq.Array(permIDs))
	return err
}

// Helper function to insert new permissions
func (s *AuthService) insertNewPermissions(ctx context.Context, tx *sql.Tx, req models.RegisterPermissionRequest, validPerms []models.PermissionsMap) error {
	var insertQuery string
	switch req.User {
	case "TUTOR":
		insertQuery = `INSERT INTO stu_tracker.Tutor_Permissions 
                      (tutor_id, permission_id) VALUES ($1, $2)
                      ON CONFLICT (tutor_id, permission_id) DO NOTHING`
	case "ADMIN":
		insertQuery = `INSERT INTO stu_tracker.Admin_Permissions 
                      (admin_id, permission_id) VALUES ($1, $2)
                      ON CONFLICT (admin_id, permission_id) DO NOTHING`
	default:
		return fmt.Errorf("invalid user type")
	}

	stmt, err := tx.PrepareContext(ctx, insertQuery)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, perm := range validPerms {
		if _, err := stmt.ExecContext(ctx, *req.ID, *perm.ID); err != nil {
			return err
		}
	}
	return nil
}
