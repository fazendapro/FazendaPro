package middleware

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/fazendapro/FazendaPro-api/config"
)

func CORSMiddleware(cfg *config.Config) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")

			if len(cfg.CORS.AllowedOrigins) > 0 {
				allowAll := false
				for _, allowedOrigin := range cfg.CORS.AllowedOrigins {
					if allowedOrigin == "*" {
						allowAll = true
						break
					}
				}

				if allowAll {
					w.Header().Set("Access-Control-Allow-Origin", "*")
				} else if origin != "" {
					originAllowed := false
					for _, allowedOrigin := range cfg.CORS.AllowedOrigins {
						if allowedOrigin == origin {
							w.Header().Set("Access-Control-Allow-Origin", origin)
							originAllowed = true
							break
						}
					}
					if !originAllowed {
						log.Printf("CORS: Origin %s not allowed", origin)
					}
				} else {
					if len(cfg.CORS.AllowedOrigins) > 0 {
						w.Header().Set("Access-Control-Allow-Origin", cfg.CORS.AllowedOrigins[0])
					}
				}
			}

			if r.Method == "OPTIONS" && len(cfg.CORS.AllowedMethods) > 0 {
				methods := strings.Join(cfg.CORS.AllowedMethods, ", ")
				w.Header().Set("Access-Control-Allow-Methods", methods)
			}

			if r.Method == "OPTIONS" && len(cfg.CORS.AllowedHeaders) > 0 {
				headers := strings.Join(cfg.CORS.AllowedHeaders, ", ")
				w.Header().Set("Access-Control-Allow-Headers", headers)
			}

			if len(cfg.CORS.ExposedHeaders) > 0 {
				exposedHeaders := strings.Join(cfg.CORS.ExposedHeaders, ", ")
				w.Header().Set("Access-Control-Expose-Headers", exposedHeaders)
			}

			if cfg.CORS.AllowCredentials && origin != "" {
				w.Header().Set("Access-Control-Allow-Credentials", "true")
			}

			if r.Method == "OPTIONS" && cfg.CORS.MaxAge > 0 {
				maxAge := fmt.Sprintf("%d", cfg.CORS.MaxAge)
				w.Header().Set("Access-Control-Max-Age", maxAge)
			}

			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
