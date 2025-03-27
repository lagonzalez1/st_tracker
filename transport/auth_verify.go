package transport

import (
	"errors"
	"strings"
	"tracker/app/models"

	"github.com/golang-jwt/jwt/v5"
)

// permissionsRequired are in the format "write:permissions"
func validateRequest(claims jwt.MapClaims, permissionsRequired string) (bool, error) {
	// Extract permissions from claims
	permissions, ok := claims["permissions"].([]models.PermissionsList)
	if !ok {
		return false, errors.New("unable to find user type")
	}
	parts := strings.Split(permissionsRequired, ":")
	action := parts[0]
	perm := parts[1]
	// Convert permissions to a slice of strings for easier comparison
	for _, p := range permissions {
		parts := strings.Split(p.Name, ":")
		userAction := parts[0]
		userPermission := parts[1]
		if userAction == action && userPermission == perm {
			return true, nil
		}
		if userAction == action && userPermission == "*" {
			return true, nil
		}
	}
	// If no match is found, return an error
	return false, errors.New("unable to validate request")
}
