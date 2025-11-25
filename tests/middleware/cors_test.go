package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/fazendapro/FazendaPro-api/config"
	"github.com/fazendapro/FazendaPro-api/internal/api/middleware"
	"github.com/fazendapro/FazendaPro-api/tests"
	"github.com/stretchr/testify/assert"
)

func TestCORSMiddleware_AllowAll(t *testing.T) {
	cfg := &config.Config{
		CORS: config.CORSConfig{
			AllowedOrigins: []string{"*"},
			AllowedMethods: []string{"GET", "POST", "PUT", "DELETE"},
			AllowedHeaders: []string{tests.HeaderContentType, tests.HeaderAuthorization},
		},
	}

	corsMiddleware := middleware.CORSMiddleware(cfg)
	handler := corsMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Origin", tests.TestOriginLocalhost)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, "*", w.Header().Get(tests.HeaderAccessControlAllowOrigin))
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestCORSMiddleware_SpecificOrigin(t *testing.T) {
	cfg := &config.Config{
		CORS: config.CORSConfig{
			AllowedOrigins: []string{"http://localhost:3000"},
			AllowedMethods: []string{"GET", "POST"},
			AllowedHeaders: []string{tests.HeaderContentType},
		},
	}

	corsMiddleware := middleware.CORSMiddleware(cfg)
	handler := corsMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Origin", tests.TestOriginLocalhost)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, tests.TestOriginLocalhost, w.Header().Get(tests.HeaderAccessControlAllowOrigin))
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestCORSMiddleware_NotAllowedOrigin(t *testing.T) {
	cfg := &config.Config{
		CORS: config.CORSConfig{
			AllowedOrigins: []string{"http://localhost:3000"},
			AllowedMethods: []string{"GET", "POST"},
			AllowedHeaders: []string{tests.HeaderContentType},
		},
	}

	corsMiddleware := middleware.CORSMiddleware(cfg)
	handler := corsMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Origin", "http://evil.com")
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Empty(t, w.Header().Get(tests.HeaderAccessControlAllowOrigin))
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestCORSMiddleware_OPTIONSRequest(t *testing.T) {
	cfg := &config.Config{
		CORS: config.CORSConfig{
			AllowedOrigins: []string{"http://localhost:3000"},
			AllowedMethods: []string{"GET", "POST", "PUT", "DELETE"},
			AllowedHeaders: []string{tests.HeaderContentType, tests.HeaderAuthorization},
		},
	}

	corsMiddleware := middleware.CORSMiddleware(cfg)
	handler := corsMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest("OPTIONS", "/test", nil)
	req.Header.Set("Origin", tests.TestOriginLocalhost)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, tests.TestOriginLocalhost, w.Header().Get(tests.HeaderAccessControlAllowOrigin))
	assert.Equal(t, "GET, POST, PUT, DELETE", w.Header().Get("Access-Control-Allow-Methods"))
	assert.Equal(t, "Content-Type, Authorization", w.Header().Get("Access-Control-Allow-Headers"))
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestCORSMiddleware_ExposedHeaders(t *testing.T) {
	cfg := &config.Config{
		CORS: config.CORSConfig{
			AllowedOrigins: []string{"http://localhost:3000"},
			ExposedHeaders: []string{"X-Total-Count"},
		},
	}

	corsMiddleware := middleware.CORSMiddleware(cfg)
	handler := corsMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Origin", tests.TestOriginLocalhost)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, "X-Total-Count", w.Header().Get("Access-Control-Expose-Headers"))
}

func TestCORSMiddleware_AllowCredentials(t *testing.T) {
	cfg := &config.Config{
		CORS: config.CORSConfig{
			AllowedOrigins:   []string{"http://localhost:3000"},
			AllowCredentials: true,
		},
	}

	corsMiddleware := middleware.CORSMiddleware(cfg)
	handler := corsMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Origin", tests.TestOriginLocalhost)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, "true", w.Header().Get("Access-Control-Allow-Credentials"))
}

func TestCORSMiddleware_MaxAge(t *testing.T) {
	cfg := &config.Config{
		CORS: config.CORSConfig{
			AllowedOrigins: []string{"http://localhost:3000"},
			MaxAge:         3600,
		},
	}

	corsMiddleware := middleware.CORSMiddleware(cfg)
	handler := corsMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest("OPTIONS", "/test", nil)
	req.Header.Set("Origin", tests.TestOriginLocalhost)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, "3600", w.Header().Get("Access-Control-Max-Age"))
}
