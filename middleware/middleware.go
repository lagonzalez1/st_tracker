package middleware

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"
	"tracker/app/config"

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
		"sub":         claims["sub"],
		"exp":         time.Now().Add(10 * time.Minute).Unix(),
		"iat":         time.Now().Unix(),
		"type":        claims["type"],
		"permissions": claims["permissions"],
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
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(secret), nil
	})
	if token.Valid {
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			return nil, JWTValidError{
				Message: "Token invalid after claims attempt",
				Code:    501,
			}
		}
		return claims, JWTValidError{}
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

func Middleware(next http.Handler) http.Handler {
	env_config, err := config.LoadConfig()
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
		claims, jwtError := validateJWT(jwtToken, env_config.JWT)
		if jwtError.Code == 501 {
			// Not valid return Unauthorized
			fmt.Printf("jwtError code 501, %v", err)
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
			new_access_token, err := generateJWTToken(cookie_claims)
			if err != nil {
				fmt.Printf("Error trying to generate new jwt token %v", err)
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte("Unauthorized"))
				return
			}
			// Set the header field with new access token
			ctx := context.WithValue(r.Context(), "props", cookie_claims)
			w.Header().Set("X-ACCESS-TOKEN", new_access_token)
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}
		if jwtError.Code == 0 {
			// No issue has been found
			fmt.Print("JWT ok.")
			ctx := context.WithValue(r.Context(), "props", claims)
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}

		fmt.Printf("Unable to verify first token %v", jwtError)
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Unauthorized"))
	})
}
