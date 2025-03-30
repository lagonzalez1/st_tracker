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
		"exp":   time.Now().Add(10 * time.Minute).Unix(),
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
	fmt.Println("ValidateJWT - Token", token.Valid)
	if token == nil {
		return nil, JWTValidError{
			Message: "Token is nill",
			Code:    501,
		}
	}
	if token.Valid {
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			return nil, JWTValidError{
				Message: "Token invalid after claims attempt",
				Code:    500,
			}
		}
		fmt.Print("Token is valid..")
		return claims, JWTValidError{Message: "OK", Code: 0}
	}
	if !token.Valid {
		return nil, JWTValidError{
			Message: "Token invalid after claims attempt",
			Code:    500,
		}
	}
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
	return nil, JWTValidError{
		Message: "Parse issues with given token",
		Code:    501,
	}
}

func Middleware(s *services.AuthService) func(http.Handler) http.Handler {
	env_config, err := config.LoadConfig()
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if err != nil {
				fmt.Println("Unable to lead env")
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte("Unable to load config files."))
			}
			authHeader := strings.Split(r.Header.Get("Authorization"), "Bearer ")
			if len(authHeader) != 2 {
				fmt.Println("Malformed token")
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte("Malformed Token"))
			}
			jwtToken := authHeader[1]
			// Check if valid JWT
			fmt.Println("JWT token", jwtToken)
			claims, jwtError := validateJWT(jwtToken, env_config.JWT)
			fmt.Println("validateJWT Code----", jwtError)
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
					fmt.Printf("Error trying to generate new jwt token %v", err)
					w.WriteHeader(http.StatusUnauthorized)
					w.Write([]byte("Unauthorized"))
					return
				}
				id, ok := cookie_claims["id"].(float64) // JWT stores numbers as float64
				if !ok {
					fmt.Println("Error trying to get id claims")
					w.WriteHeader(http.StatusUnauthorized)
					w.Write([]byte("Unable to parse id claims id as float64"))
					return
				}
				userType, ok := cookie_claims["type"].(string)
				if !ok {
					fmt.Println("Error trying to get type claims")
					w.WriteHeader(http.StatusUnauthorized)
					w.Write([]byte("Unable to parse user type in claims"))
					return
				}
				orgID, ok := cookie_claims["orgid"].(float64) // Convert orgid from float64 to int64
				if !ok {
					fmt.Println("Error trying to get orgid claims")
					w.WriteHeader(http.StatusUnauthorized)
					w.Write([]byte("Unable to parse orgid"))
					return
				}
				// Need a fast lookup like redis
				// Attach permissions
				// This can be faster if using redis or something in memory
				permissions, err := GetPermissions(int64(id), userType, int64(orgID), s)
				if err != nil {
					fmt.Println("Error trying to get permissions claims")
					w.WriteHeader(http.StatusUnauthorized)
					w.Write([]byte("Unable to get user permissions"))
					return
				}
				fmt.Println("Sending new access token ---", new_access_token)
				cookie_claims["permissions"] = permissions
				// Set the header field with new access token
				ctx := context.WithValue(r.Context(), "props", cookie_claims)
				w.Header().Set("X-ACCESS-TOKEN", new_access_token)
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}
			if jwtError.Code == 0 {
				id, ok := claims["id"].(float64) // JWT stores numbers as float64
				if !ok {
					fmt.Println(err)
					w.WriteHeader(http.StatusUnauthorized)
					w.Write([]byte("Unable to parse id"))
					return
				}
				userType, ok := claims["type"].(string)
				if !ok {
					fmt.Println(err)
					w.WriteHeader(http.StatusUnauthorized)
					w.Write([]byte("Unable to parse user type"))
					return
				}

				orgID, ok := claims["orgid"].(float64) // Convert orgid from float64 to int64
				if !ok {
					fmt.Println(err)
					w.WriteHeader(http.StatusUnauthorized)
					w.Write([]byte("Unable to parse orgid"))
					return
				}
				permissions, err := GetPermissions(int64(id), userType, int64(orgID), s)
				if err != nil {
					fmt.Println(err)
					w.WriteHeader(http.StatusUnauthorized)
					w.Write([]byte("Unable to get user permissions"))
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

func GetPermissions(id int64, userType string, orgid int64, s *services.AuthService) ([]models.PermissionsList, error) {
	rows, err := s.GetPermissionsById(orgid, userType, id)
	if err != nil {
		return nil, fmt.Errorf("error querying permissions: %w", err)
	}
	return rows, nil

}
