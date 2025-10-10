package middleware

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

type ErrorResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
	Message string `json:"message"`
	Code    int    `json:"code"`
}

func SendErrorResponse(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	response := ErrorResponse{
		Success: false,
		Error:   http.StatusText(statusCode),
		Message: message,
		Code:    statusCode,
	}

	json.NewEncoder(w).Encode(response)
}

func Auth(jwtSecret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tokenString := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
			if tokenString == "" {
				SendErrorResponse(w, "Token não fornecido", http.StatusUnauthorized)
				return
			}

			token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
				return []byte(jwtSecret), nil
			})

			if err != nil || !token.Valid {
				SendErrorResponse(w, "Token inválido", http.StatusUnauthorized)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
