package middleware

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

const (
	HeaderContentType = "Content-Type"
	ContentTypeJSON   = "application/json"
)

type ErrorResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
	Message string `json:"message"`
	Code    int    `json:"code"`
}

func SendErrorResponse(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set(HeaderContentType, ContentTypeJSON)
	w.WriteHeader(statusCode)

	response := ErrorResponse{
		Success: false,
		Error:   http.StatusText(statusCode),
		Message: message,
		Code:    statusCode,
	}

	json.NewEncoder(w).Encode(response)
}

func validateToken(tokenString, jwtSecret string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtSecret), nil
	})
	if err != nil || !token.Valid {
		return nil, err
	}
	return token, nil
}

func extractFarmID(token *jwt.Token) (uint, bool) {
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, false
	}
	farmID, exists := claims["farm_id"]
	if !exists {
		return 0, false
	}
	farmIDFloat, ok := farmID.(float64)
	if !ok {
		return 0, false
	}
	return uint(farmIDFloat), true
}

func Auth(jwtSecret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tokenString := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
			if tokenString == "" {
				SendErrorResponse(w, "Token não fornecido", http.StatusUnauthorized)
				return
			}

			token, err := validateToken(tokenString, jwtSecret)
			if err != nil {
				SendErrorResponse(w, "Token inválido", http.StatusUnauthorized)
				return
			}

			farmID, ok := extractFarmID(token)
			if !ok || farmID == 0 {
				SendErrorResponse(w, "Token não contém farm_id válido", http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), "farm_id", farmID)
			r = r.WithContext(ctx)

			next.ServeHTTP(w, r)
		})
	}
}
