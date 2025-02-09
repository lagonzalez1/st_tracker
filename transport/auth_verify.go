package transport

import (
	"errors"
	"slices"

	"github.com/golang-jwt/jwt/v5"
)

func validateRequest(claims jwt.MapClaims, permissionsRequired string) (bool, error) {
	// Extract permissions from claims
	permissions, ok := claims["permissions"].([]interface{})
	if !ok {
		return false, errors.New("permissions not found in claims")
	}

	// Convert permissions to a slice of strings for easier comparison
	var permissionsList []string
	for _, p := range permissions {
		if str, ok := p.(string); ok {
			permissionsList = append(permissionsList, str)
		}
	}

	// Check if "all" is present (wildcard for global permissions)
	if slices.Contains(permissionsList, "all") {
		return true, nil
	}

	if slices.Contains(permissionsList, permissionsRequired) {
		return true, nil
	}
	// If no match is found, return an error
	return false, errors.New("Unable to validate request")
}
