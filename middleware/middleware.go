package middleware

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"
	"tracker/app/config"
	"tracker/app/models"
	"tracker/app/services"

	"github.com/golang-jwt/jwt/v5"
)

func generateJWTToken(claims jwt.MapClaims) (string, error) {
	env_config, err := config.LoadConfig()
	if err != nil {
		return fmt.Sprintf("unable to load config env"), err
	}
	jwt_token := env_config.JWT
	secret_key := []byte(jwt_token)
	// Need to create a role to return here ?
	jwt_object := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":   claims["sub"],
		"exp":   time.Now().Add(1 * time.Hour).Unix(),
		"iat":   time.Now().Unix(),
		"type":  claims["type"],
		"id":    claims["id"],
		"orgid": claims["orgid"],
	})
	token_string, err := jwt_object.SignedString(secret_key)
	if err != nil {
		return fmt.Sprintf("Unable to create JWT token"), err
	}
	return token_string, nil
}

type JWTValidError struct {
	Message string
	Code    int
}

func validateJWT(tokenString, secret string) (jwt.MapClaims, JWTValidError) {
	if tokenString == "" {
		return nil, JWTValidError{
			Message: "Token is ErrTokenSignatureInvalid and or ErrTokenNotValidYet",
			Code:    501,
		}
	}
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(secret), nil
	})

	// Check for token expired ?
	if err != nil {
		switch {
		case errors.Is(err, jwt.ErrTokenExpired):
			return nil, JWTValidError{
				Message: "Token expired",
				Code:    500,
			}
		case errors.Is(err, jwt.ErrTokenSignatureInvalid) || errors.Is(err, jwt.ErrTokenNotValidYet):
			return nil, JWTValidError{
				Message: "Token is ErrTokenSignatureInvalid and or ErrTokenNotValidYet",
				Code:    501,
			}
		default:
			return nil, JWTValidError{
				Message: "Token is not handled",
				Code:    501,
			}
		}
	}
	if token == nil {
		return nil, JWTValidError{
			Message: "Token is nill",
			Code:    501,
		}
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, JWTValidError{Message: "OK", Code: 0}
	}

	return nil, JWTValidError{
		Message: "Parse issues with given token",
		Code:    501,
	}
}

func Middleware(s *services.AuthService) func(http.Handler) http.Handler {
	env_config, err := config.LoadConfig()
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
			defer cancel()
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte("Unable to load config files."))
			}
			if r.Header.Get("Authorization") == "" {
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte("Malformed Token"))
			}
			authHeader := strings.Split(r.Header.Get("Authorization"), "Bearer ")
			if len(authHeader) != 2 {
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte("Malformed Token"))
			}
			jwtToken := authHeader[1]
			// Check if valid JWT
			claims, jwtError := validateJWT(jwtToken, env_config.JWT)
			fmt.Println("1. Claims-", claims)

			if jwtError.Code == 501 {
				// Not valid return Unauthorized
				fmt.Printf("jwtError code 401, %v", err)
				fmt.Println(jwtError.Message)
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte("Unauthorized"))
				return
			}
			if jwtError.Code == 500 {
				cookie_auth, err := r.Cookie("_auth")
				if err != nil {
					fmt.Println(err)
					w.WriteHeader(http.StatusUnauthorized)
					w.Write([]byte("Unable to parse cookie"))
					return
				}
				// Validate cookie to ensure new creation of token is valid
				cookie_claims, jwtError := validateJWT(cookie_auth.Value, env_config.JWT)
				// If cookie token is expired reject.
				if jwtError.Code == 500 || jwtError.Code == 501 {
					fmt.Printf("jwtError on cookie claims code 501, %v", err)
					w.WriteHeader(http.StatusUnauthorized)
					w.Write([]byte("Unauthorized"))
					return
				}
				// Create new access token
				fmt.Println("Cookie claims----", cookie_claims)
				new_access_token, err := generateJWTToken(cookie_claims)
				if err != nil {
					fmt.Println(err)
					w.WriteHeader(http.StatusUnauthorized)
					w.Write([]byte("Unable to parse cookie"))
					return
				}
				permissions, err := ValidateRequestPermissions(ctx, cookie_claims, s)
				if err != nil {
					fmt.Println(err)
					w.WriteHeader(http.StatusUnauthorized)
					w.Write([]byte("Unable to parse cookie"))
					return
				}
				fmt.Println("Sending new access token ---", new_access_token)
				cookie_claims["permissions"] = permissions
				// Set the header field with new access token
				ctx := context.WithValue(r.Context(), "props", cookie_claims)
				w.Header().Set("X-Access-Token", new_access_token)
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}
			if jwtError.Code == 0 {
				permissions, err := ValidateRequestPermissions(ctx, claims, s)
				if err != nil {
					fmt.Println(err)
					w.WriteHeader(http.StatusUnauthorized)
					w.Write([]byte("Unable to parse cookie"))
					return
				}
				claims["permissions"] = permissions
				// No issue has been found
				fmt.Print("JWT ok.")
				ctx := context.WithValue(r.Context(), "props", claims)
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}
			fmt.Printf("Unable to verify token %v", jwtError)
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Unauthorized"))
		})
	}
}

func ValidateRequestPermissions(c context.Context, claims jwt.MapClaims, s *services.AuthService) ([]models.PermissionsList, error) {
	if claims == nil {
		return nil, fmt.Errorf("unable to get request")
	}
	id, ok := claims["id"].(float64) // JWT stores numbers as float64
	if !ok {
		return nil, fmt.Errorf("unable to parse id")
	}
	userType, ok := claims["type"].(string)
	if !ok {
		return nil, fmt.Errorf("unable to parse user type")
	}
	orgID, ok := claims["orgid"].(float64) // Convert orgid from float64 to int64
	if !ok {
		return nil, fmt.Errorf("unable to parse organization id")
	}
	permissions, err := GetPermissions(c, int64(id), userType, int64(orgID), s)
	if err != nil {
		return nil, fmt.Errorf("unable to get permissions from validation")
	}
	return permissions, nil
}

func GetPermissions(c context.Context, id int64, userType string, orgid int64, s *services.AuthService) ([]models.PermissionsList, error) {
	rows, err := s.GetPermissionsById(c, orgid, userType, id)
	if err != nil {
		return nil, fmt.Errorf("error querying permissions: %w", err)
	}
	return rows, nil

}
