package middleware

import (
	"context"
	"encoding/json"
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

/*
	Valkey Cache

		Callstack:
			Request ->
				Middleware ->
					Auth Middleware
					Pass claims
					Reject claims.
			Response

		Check cache for hit/miss to pass claims to next call.
		If JWT expiration is valid pass claims, invalidate and clear cache.
		Revalidate cache for miss, when expiration event occurs.

		JWT Payload { id: email: type: orgid, exp: iat } -> { str: str: str: int: int: int }


*/

type CacheClaims struct {
	ID    int64  `json:"id"`
	Sub   string `json:"sub"`
	Type  string `json:"type"`
	OrgId int64  `json:"orgid"`
	Exp   int64  `json:"exp"`
	Iat   int64  `json:"iat"`
}

type JWTValidError struct {
	Message string
	Code    int
}

func Middleware(s *services.AuthService) func(http.Handler) http.Handler {
	env_config, err := config.LoadConfig()
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
			defer cancel()
			if err != nil {
				http.Error(w, "Unable to load config files.", http.StatusUnauthorized)
				return
			}

			// Short duration token
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "Authorization header is required", http.StatusUnauthorized)
				return
			}
			// Split the token and validate format
			parts := strings.SplitN(authHeader, "Bearer ", 2)
			if len(parts) != 2 {
				http.Error(w, "Malformed Token - must be 'Bearer <token>'", http.StatusUnauthorized)
				return
			}

			jwtToken := strings.TrimSpace(parts[1]) // Also trim any whitespace
			if jwtToken == "" {
				http.Error(w, "Empty Token", http.StatusUnauthorized)
				return
			}
			cacheClaims, isHit, _ := GetCacheClaims(ctx, jwtToken, s)
			if isHit {
				cachePermissions, isValid, _ := GetCachePermissions(ctx, cacheClaims, s)
				if isValid {
					cacheClaims["permissions"] = cachePermissions
					ctx := context.WithValue(r.Context(), "props", cacheClaims)
					next.ServeHTTP(w, r.WithContext(ctx))
					return
				} else {
					permissions, err := ValidateRequestPermissions(ctx, cacheClaims, s)
					if err != nil {
						fmt.Println(err)
						w.WriteHeader(http.StatusUnauthorized)
						w.Write([]byte("Unable to parse cookie"))
						return
					}
					cachePerm, err := CachePermissions(ctx, cacheClaims, permissions, s)
					if err != nil {
						fmt.Println(err)
						w.WriteHeader(http.StatusUnauthorized)
						w.Write([]byte("Unable to parse cookie"))
						return
					}
					if cachePerm {
						cacheClaims["permissions"] = permissions
						// No issue has been found
						fmt.Print("2.JWT ok pem.", permissions)
						ctx := context.WithValue(r.Context(), "props", cacheClaims)
						next.ServeHTTP(w, r.WithContext(ctx))
						return
					}
				}
			}
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
					fmt.Println(err)
					w.WriteHeader(http.StatusUnauthorized)
					w.Write([]byte("Unable to parse cookie"))
					return
				}

				cacheClaims, err := CacheTokenClaims(ctx, cookie_claims, new_access_token, s)
				if err != nil || !cacheClaims {
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

				cachePerm, err := CachePermissions(ctx, cookie_claims, permissions, s)
				if err != nil || !cachePerm {
					fmt.Println(err)
					w.WriteHeader(http.StatusUnauthorized)
					w.Write([]byte("Unable to parse cookie"))
					return
				}

				fmt.Println("Sending new access token ", new_access_token)
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
				cachePerm, err := CachePermissions(ctx, claims, permissions, s)
				if err != nil || !cachePerm {
					fmt.Println(err)
					w.WriteHeader(http.StatusUnauthorized)
					w.Write([]byte("Unable to parse cookie"))
					return
				}

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

func generateJWTToken(claims jwt.MapClaims) (string, error) {
	env_config, err := config.LoadConfig()
	if err != nil {
		return fmt.Sprintf("unable to load config env %v", err), err
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
		return fmt.Sprintf("Unable to create JWT token %v", err), err
	}
	return token_string, nil
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

func CreateClaims(value string) (jwt.MapClaims, error) {
	delimiter := ":"
	values := strings.Split(value, delimiter)
	if len(values) <= 0 {
		return nil, fmt.Errorf("parsing value from cache results in empty set")
	}
	return jwt.MapClaims{}, nil
}

func GetCacheClaims(ctx context.Context, token string, s *services.AuthService) (jwt.MapClaims, bool, error) {
	userKey := fmt.Sprintf("auth:token:%s", token)
	result := s.Valkey().Do(ctx, s.Valkey().B().Get().Key(userKey).Build())
	claimsStr, err := result.ToString()
	if err != nil {
		return nil, false, fmt.Errorf("unable to parse toString")
	}
	var claims jwt.MapClaims
	if err := json.Unmarshal([]byte(claimsStr), &claims); err != nil {
		return nil, false, fmt.Errorf("unable to unmarshal claims")
	}
	return claims, true, nil
}

func GetCachePermissions(ctx context.Context, claims jwt.MapClaims, s *services.AuthService) (map[string]bool, bool, error) {
	id, ok := claims["id"].(float64)
	if !ok {
		return nil, false, fmt.Errorf("unable to parse id")
	}
	userType, ok := claims["type"].(string)
	if !ok {
		return nil, false, fmt.Errorf("unable to parse user type")
	}
	userKey := fmt.Sprintf("auth:perm:%d:%s", int64(id), userType)
	result := s.Valkey().Do(ctx, s.Valkey().B().Get().Key(userKey).Build())
	claimsStr, err := result.ToString()
	if err != nil {
		return nil, false, fmt.Errorf("unable to parse toString")
	}
	var permissions map[string]bool
	if err := json.Unmarshal([]byte(claimsStr), &permissions); err != nil {
		return nil, false, fmt.Errorf("unable to unmarshal claims")
	}
	return permissions, true, nil
}

func ValidateRequestPermissions(c context.Context, claims jwt.MapClaims, s *services.AuthService) (map[string]bool, error) {
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

func GetPermissions(c context.Context, id int64, userType string, orgid int64, s *services.AuthService) (map[string]bool, error) {
	rows, err := s.GetPermissionsById(c, orgid, userType, id)
	if err != nil {
		return nil, fmt.Errorf("error querying permissions: %w", err)
	}
	perm := generatePermissionsMap(rows)
	return perm, nil

}

func generatePermissionsMap(permissions []models.PermissionsList) map[string]bool {
	permissionsClaims := make(map[string]bool)
	for i := 0; i < len(permissions); i++ {
		permissionsClaims[permissions[i].Name] = true
	}
	return permissionsClaims
}

func CacheTokenClaims(ctx context.Context, claims jwt.MapClaims, token string, s *services.AuthService) (bool, error) {
	exp := time.Now().Add(1 * time.Hour)
	userKey := fmt.Sprintf("auth:token:%s", token)
	claimsJson, err := json.Marshal(claims)
	if err != nil {
		return false, fmt.Errorf("unable to marshal claims for cache %v", err)
	}
	res := s.Valkey().Do(ctx, s.Valkey().B().Set().Key(userKey).Value(string(claimsJson)).Exat(exp).Build()).Error()
	if res != nil {
		return false, fmt.Errorf("unable to cache login %v", err)
	}
	return true, nil
}

func CachePermissions(ctx context.Context, claims jwt.MapClaims, permissions map[string]bool, s *services.AuthService) (bool, error) {
	exp := time.Now().Add(1 * time.Hour)
	id, ok := claims["id"].(float64) // JWT stores numbers as float64
	if !ok {
		return false, fmt.Errorf("unable to parse id")
	}
	userType, ok := claims["type"].(string)
	if !ok {
		return false, fmt.Errorf("unable to parse user type")
	}
	permissionsKey := fmt.Sprintf("auth:perm:%d:%s", int64(id), userType)
	permClaims, err := json.Marshal(permissions)
	if err != nil {
		return false, fmt.Errorf("unable to marshal claims for cache %v", err)
	}
	err = s.Valkey().Do(ctx, s.Valkey().B().Set().Key(permissionsKey).Value(string(permClaims)).Exat(exp).Build()).Error()
	if err != nil {
		return false, fmt.Errorf("unable to cache permissions: %w", err)
	}
	return true, nil
}
