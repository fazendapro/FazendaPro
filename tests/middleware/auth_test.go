package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/fazendapro/FazendaPro-api/internal/api/middleware"
	"github.com/fazendapro/FazendaPro-api/tests"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

func generateValidToken(secret string, farmID uint) string {
	claims := jwt.MapClaims{
		"sub":     1,
		"email":   tests.TestEmailExample,
		"farm_id": farmID,
		"iat":     time.Now().Unix(),
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString([]byte(secret))
	return tokenString
}

func TestAuthMiddleware_ValidToken(t *testing.T) {
	jwtSecret := tests.TestSecret
	authMiddleware := middleware.Auth(jwtSecret)

	handler := authMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		farmID := r.Context().Value("farm_id")
		assert.NotNil(t, farmID)
		assert.Equal(t, uint(1), farmID.(uint))
		w.WriteHeader(http.StatusOK)
	}))

	token := generateValidToken(jwtSecret, 1)
	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set(tests.HeaderAuthorization, tests.BearerPrefix+token)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestAuthMiddleware_MissingToken(t *testing.T) {
	jwtSecret := tests.TestSecret
	authMiddleware := middleware.Auth(jwtSecret)

	handler := authMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestAuthMiddleware_InvalidToken(t *testing.T) {
	jwtSecret := tests.TestSecret
	authMiddleware := middleware.Auth(jwtSecret)

	handler := authMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer invalid-token")
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestAuthMiddleware_WrongSecret(t *testing.T) {
	jwtSecret := tests.TestSecret
	authMiddleware := middleware.Auth(jwtSecret)

	handler := authMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	token := generateValidToken("wrong-secret", 1)
	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set(tests.HeaderAuthorization, tests.BearerPrefix+token)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestAuthMiddleware_ExpiredToken(t *testing.T) {
	jwtSecret := tests.TestSecret
	authMiddleware := middleware.Auth(jwtSecret)

	handler := authMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	claims := jwt.MapClaims{
		"sub":     1,
		"email":   tests.TestEmailExample,
		"farm_id": 1,
		"iat":     time.Now().Add(-time.Hour * 25).Unix(),
		"exp":     time.Now().Add(-time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString([]byte(jwtSecret))

	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer "+tokenString)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestAuthMiddleware_NoBearerPrefix(t *testing.T) {
	jwtSecret := tests.TestSecret
	authMiddleware := middleware.Auth(jwtSecret)

	handler := authMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "not-bearer-token")
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}
