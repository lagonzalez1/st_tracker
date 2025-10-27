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

func ExtractInt64Claim(claims jwt.MapClaims, key string) (int64, error) {
	val, ok := claims[key]
	if !ok {
		return 0, fmt.Errorf("claim not found: %v", key)
	}
	f, ok := val.(float64)
	if !ok {
		return 0, fmt.Errorf("claim not a float: %v", key)
	}
	return int64(f), nil
}

func ExtractStringClaims(claims jwt.MapClaims, key string) (string, error) {
	val, ok := claims[key]
	if !ok {
		return "", fmt.Errorf("claim not found: %v", key)
	}
	f, ok := val.(string)
	if !ok {
		return "", fmt.Errorf("claim not a string: %v", key)
	}
	return f, nil
}
