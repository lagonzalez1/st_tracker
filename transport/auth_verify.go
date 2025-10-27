package transport

import (
	"errors"
	"fmt"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

// permissionsRequired are in the format "write:permissions"
func validateRequest(claims jwt.MapClaims, permissionsRequired string) (bool, error) {
	// Extract permissions from claims
	permissions, ok := claims["permissions"]
	if !ok {
		return false, errors.New("unable to parse permissions claims")
	}
	permMap, ok := permissions.(map[string]bool)
	if !ok {
		return false, errors.New("unable to parse claims to map")
	}

	parts := strings.Split(permissionsRequired, ":")
	action := parts[0]
	rootCheck := fmt.Sprintf("%s:*", action)
	_, ok = permMap[rootCheck]
	if ok {
		return true, nil
	}
	_, ok = permMap[permissionsRequired]
	if ok {
		return true, nil
	}
	// If no match is found, return an error
	return false, fmt.Errorf("unable to validate request, invalid permission")
}
