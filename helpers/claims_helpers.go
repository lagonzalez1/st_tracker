package helpers

import (
	"fmt"

	"github.com/golang-jwt/jwt/v5"
)

func ExtractFloat64Claim(claims jwt.MapClaims, key string) (float64, error) {
	val, ok := claims[key]
	if !ok {
		return 0, fmt.Errorf("claim not found: %v", key)
	}
	f, ok := val.(float64)
	if !ok {
		return 0, fmt.Errorf("claim not a float: %v", key)
	}
	return f, nil
}
