package services

import (
	"fmt"
	"strings"
	"tracker/app/models"

	"github.com/lib/pq"
)

func (s *AuthService) CreatePermission(req models.RegisterPermissionRequest) (*models.RegisterPermissionResponse, error) {
	// Input validation
	if req.OrganizationId == nil || req.Role == "" || len(req.Permissions) == 0 {
		return nil, fmt.Errorf("missing required fields: org id and user type")
	}
	// Query permissions based on the provided permission IDs
	permissionQuery := `SELECT p.id, p.name FROM stu_tracker.Permissions p WHERE p.organization_id = $1 AND p.name = ANY($2);`
	rows, err := s.db.Query(permissionQuery, req.OrganizationId, pq.Array(req.Permissions))
	if err != nil {
		return nil, fmt.Errorf("error querying permissions: %w", err)
	}
	defer rows.Close()

	var permissions []models.PermissionsMap
	for rows.Next() {
		var permission models.PermissionsMap
		err := rows.Scan(&permission.ID, &permission.PermissionName)
		if err != nil {
			return nil, fmt.Errorf("error scanning row: %w", err)
		}
		permissions = append(permissions, permission)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}
	fmt.Print(permissions)
	// If no permissions were updated, return response
	if len(permissions) == 0 {
		return &models.RegisterPermissionResponse{
			Status: "No permissions were updated",
		}, nil
	}

	// Build query dynamically based on user type
	var sb strings.Builder
	if req.User == "TUTOR" {
		sb.WriteString("INSERT INTO stu_tracker.Tutor_Permissions (tutor_id, permission_id) VALUES ")
	} else if req.User == "ADMIN" {
		sb.WriteString("INSERT INTO stu_tracker.Admin_Permissions (admin_id, permission_id) VALUES ")
	} else {
		return nil, fmt.Errorf("invalid user type")
	}
	// Construct values for bulk insertion
	values := make([]string, len(permissions))
	for i, p := range permissions {
		values[i] = fmt.Sprintf("(%d, %d)", *req.ID, *p.ID)
	}
	sb.WriteString(strings.Join(values, ", "))
	sb.WriteString(";")
	query := sb.String()

	fmt.Print(sb.String())
	// Execute the query
	_, err = s.db.Exec(query)
	if err != nil {
		return nil, fmt.Errorf("failed to insert Locations: %w", err)
	}

	return &models.RegisterPermissionResponse{
		Status: "OK",
	}, nil
}
