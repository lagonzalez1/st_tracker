package middleware

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
	"tracker/app/cache"
	"tracker/app/models"
	"tracker/app/services"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
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

func Middleware(s *services.AuthService, c *cache.CacheHandler, cfg *models.Config) func(http.Handler) http.Handler {

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
			defer cancel()
			meteredPaths := map[string]string{
				"POST /api/create_student":     "max_students_per_location",
				"POST /api/create_tutor":       "max_tutors_per_location",
				"POST /api/create_location":    "max_locations_per_district",
				"POST /api/create_admin_staff": "max_admin_per_district",
				"POST /api/max_districts":      "max_districts",
				"POST /api/micro_generate":     "max_llm_tokens",
			}
			routeKey := fmt.Sprintf("%s %s", r.Method, r.URL.Path)
			entitlementKey, exist := meteredPaths[routeKey]

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
			// First check the cache and plan limits
			if isHit && exist {
				planLimitsOk, err := CheckPlanLimits(ctx, cacheClaims, entitlementKey, s, c, r.Header)
				fmt.Printf("plan limits ok %v", planLimitsOk)
				if err != nil {
					fmt.Println(err)
					http.Error(w, err.Error(), http.StatusUpgradeRequired)
					return
				}
				if !planLimitsOk {
					fmt.Println(err)
					http.Error(w, err.Error(), http.StatusUpgradeRequired)
					return
				}
			}

			// If a cache hit process the request
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
						http.Error(w, "unable to validate permissions", http.StatusUnauthorized)
						return
					}
					cachePerm, err := CachePermissions(ctx, cacheClaims, permissions, s)
					if err != nil {
						fmt.Println(err)
						http.Error(w, "unable to cache permissions", http.StatusUnauthorized)
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
			// Unable to hit cache validate using db
			claims, jwtError := validateJWT(jwtToken, cfg.JWT)
			switch jwtError.Code {
			case 500:
				cookie_auth, err := r.Cookie("_auth")
				if err != nil {
					http.Error(w, "invalid cookie", http.StatusUnauthorized)
					return
				}
				// Validate cookie to ensure new creation of token is valid
				cookie_claims, jwtError := validateJWT(cookie_auth.Value, cfg.JWT)
				// If cookie token is expired reject.
				if jwtError.Code == 500 || jwtError.Code == 501 {
					http.Error(w, "invalid cookie token", http.StatusUnauthorized)
					return
				}
				if cookie_claims != nil && exist {
					planLimitsOk, err := CheckPlanLimits(ctx, cookie_claims, entitlementKey, s, c, r.Header)
					fmt.Printf("plan limits ok: %v", planLimitsOk)
					if err != nil {
						fmt.Println(err)
						http.Error(w, err.Error(), http.StatusUpgradeRequired)
						return
					}
					if !planLimitsOk {
						fmt.Println(err)
						http.Error(w, err.Error(), http.StatusUpgradeRequired)
						return
					}
				}
				new_access_token, err := generateJWTToken(cookie_claims, cfg.JWT)
				if err != nil {
					http.Error(w, "unable to create new token", http.StatusUnauthorized)

					return
				}
				cacheClaims, err := CacheTokenClaims(ctx, cookie_claims, new_access_token, s)
				if err != nil || !cacheClaims {
					http.Error(w, "unable to cache new token", http.StatusUnauthorized)
					return
				}
				permissions, err := ValidateRequestPermissions(ctx, cookie_claims, s)
				if err != nil {
					http.Error(w, "unable to validate permissions", http.StatusUnauthorized)
					return
				}
				cachePerm, err := CachePermissions(ctx, cookie_claims, permissions, s)
				if err != nil || !cachePerm {
					http.Error(w, "unable to cache permissions", http.StatusUnauthorized)
					return
				}

				cookie_claims["permissions"] = permissions
				ctx := context.WithValue(r.Context(), "props", cookie_claims)
				w.Header().Set("X-Access-Token", new_access_token)
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			case 0:
				if claims != nil && exist {
					planLimitsOk, err := CheckPlanLimits(ctx, claims, entitlementKey, s, c, r.Header)
					fmt.Printf("plan limits ok %v", planLimitsOk)
					if err != nil {
						fmt.Println(err)
						http.Error(w, err.Error(), http.StatusUpgradeRequired)
						return
					}
					if !planLimitsOk {
						fmt.Println(err)
						http.Error(w, err.Error(), http.StatusUpgradeRequired)
						return
					}
				}
				permissions, err := ValidateRequestPermissions(ctx, claims, s)
				if err != nil {
					fmt.Println(err)
					http.Error(w, "unable to validate permissions", http.StatusUnauthorized)
					return
				}
				claims["permissions"] = permissions
				cachePerm, err := CachePermissions(ctx, claims, permissions, s)
				if err != nil || !cachePerm {
					fmt.Println(err)
					http.Error(w, "unable to cache permissions", http.StatusUnauthorized)
					return
				}
				ctx := context.WithValue(r.Context(), "props", claims)
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}

			http.Error(w, "unable to verify request", http.StatusUnauthorized)
		})
	}
}

func generateJWTToken(claims jwt.MapClaims, token string) (string, error) {
	secret_key := []byte(token)
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

// Dangerous to assume this will always return
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

func GetCachePlanEntitlements(c context.Context, orgid int64, s *services.AuthService, cache *cache.CacheHandler) ([]models.OrganizationPlanEntitlement, error) {
	key := fmt.Sprintf("auth:plan-entitlement:%d", orgid)
	lkey := cache.LockKey(key)
	data, isHit := cache.CheckCache(c, key)
	if isHit {
		var organizationEntitlement []models.OrganizationPlanEntitlement
		if err := json.Unmarshal([]byte(data), &organizationEntitlement); err != nil {
			return nil, fmt.Errorf("unable to unmarshal plan entitlements")
		}
		return organizationEntitlement, nil
	}
	token := uuid.NewString()
	lock, _ := cache.TryAcquireLock(c, lkey, token)
	if lock {
		rows, err := s.GetPlanEntitlements(c, &orgid)
		if err != nil {
			if errors.Is(err, context.DeadlineExceeded) {
				return nil, fmt.Errorf("deadline exceeded")
			}
			if errors.Is(err, context.Canceled) {
				return nil, fmt.Errorf("context canceled")
			}
			return nil, fmt.Errorf("deadline exceeded")
		}
		cache.SetCacheByString(c, key, string(rows))
		cache.SafeUnlock(c, lkey, token)
	}
	status := cache.WaitForCacheUpdate(c, key)
	if status {
		cdata, isHit := cache.CheckCache(c, key)
		if isHit {
			var organizationEntitlement []models.OrganizationPlanEntitlement
			if err := json.Unmarshal([]byte(cdata), &organizationEntitlement); err != nil {
			}
			return organizationEntitlement, nil
		}
	}
	return nil, fmt.Errorf("unable to get data")

}

func CheckPlanLimits(ctx context.Context, claims jwt.MapClaims, entitlementKey string,
	s *services.AuthService, c *cache.CacheHandler, header http.Header) (bool, error) {
	orgid, ok := claims["orgid"].(float64)
	if !ok {
		return false, fmt.Errorf("unable to get claims value")
	}
	var locationID *int64
	locationStr := header.Get("X-Location-Id")
	if locationStr != "" {
		id, err := strconv.ParseInt(locationStr, 10, 64)
		if err != nil {
			return false, fmt.Errorf("unable to prase location to int")
		}
		locationID = &id
	}
	var districtID *int64
	districtStr := header.Get("X-District-Id")
	if districtStr != "" {
		id, err := strconv.ParseInt(districtStr, 10, 64)
		if err != nil {
			return false, fmt.Errorf("unable to prase District to int")
		}
		districtID = &id
	}
	planEntitlements, err := GetCachePlanEntitlements(ctx, int64(orgid), s, c)
	if err != nil {
		return false, err
	}
	usage, err := CheckPlanEntitlements(ctx, int64(orgid), entitlementKey, locationID, districtID, s, planEntitlements)
	if err != nil || !usage {
		return false, err
	}
	return usage, nil
}

func CheckPlanEntitlements(c context.Context, orgid int64, key string, locationID *int64, districtID *int64, s *services.AuthService, plan []models.OrganizationPlanEntitlement) (bool, error) {
	for _, value := range plan {
		if key == *value.ActionKey {
			usage, err := s.CheckUsage(c, int(orgid), key, locationID, districtID)
			if err != nil {
				return false, fmt.Errorf("unable to check usage.")
			}
			if usage == nil || value.LimitValue == nil {
				return false, fmt.Errorf("encountered error")
			}
			if *usage == *value.LimitValue {
				return false, fmt.Errorf("Plan limits reached for %s", *value.ActionKey)
			}
			if *usage == *value.LimitValue {
				return false, fmt.Errorf("Plan limits reached for %s", *value.ActionKey)
			}
			fmt.Printf("plan limits check-> usage: %v, planValue: %v", *usage, *value.LimitValue)
			if *usage < *value.LimitValue {
				return true, nil
			}
		}
	}

	return false, fmt.Errorf("unable to determine usage")
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
