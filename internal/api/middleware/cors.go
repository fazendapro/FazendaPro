package middleware

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/fazendapro/FazendaPro-api/config"
)

const HeaderAccessControlAllowOrigin = "Access-Control-Allow-Origin"

func handleOrigin(w http.ResponseWriter, origin string, allowedOrigins []string) {
	if len(allowedOrigins) == 0 {
		return
	}

	allowAll := false
	for _, allowedOrigin := range allowedOrigins {
		if allowedOrigin == "*" {
			allowAll = true
			break
		}
	}

	if allowAll {
		w.Header().Set(HeaderAccessControlAllowOrigin, "*")
		return
	}

	if origin == "" {
		if len(allowedOrigins) > 0 {
			w.Header().Set(HeaderAccessControlAllowOrigin, allowedOrigins[0])
		}
		return
	}

	originAllowed := false
	for _, allowedOrigin := range allowedOrigins {
		if allowedOrigin == origin {
			w.Header().Set(HeaderAccessControlAllowOrigin, origin)
			originAllowed = true
			break
		}
	}

	if !originAllowed {
		log.Printf("CORS: Origin %s not allowed", origin)
	}
}

func handleOptionsRequest(w http.ResponseWriter, cfg *config.Config) {
	if len(cfg.CORS.AllowedMethods) > 0 {
		methods := strings.Join(cfg.CORS.AllowedMethods, ", ")
		w.Header().Set("Access-Control-Allow-Methods", methods)
	}

	if len(cfg.CORS.AllowedHeaders) > 0 {
		headers := strings.Join(cfg.CORS.AllowedHeaders, ", ")
		w.Header().Set("Access-Control-Allow-Headers", headers)
	}

	if cfg.CORS.MaxAge > 0 {
		maxAge := fmt.Sprintf("%d", cfg.CORS.MaxAge)
		w.Header().Set("Access-Control-Max-Age", maxAge)
	}
}

func CORSMiddleware(cfg *config.Config) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")

			handleOrigin(w, origin, cfg.CORS.AllowedOrigins)

			if len(cfg.CORS.ExposedHeaders) > 0 {
				exposedHeaders := strings.Join(cfg.CORS.ExposedHeaders, ", ")
				w.Header().Set("Access-Control-Expose-Headers", exposedHeaders)
			}

			if cfg.CORS.AllowCredentials && origin != "" {
				w.Header().Set("Access-Control-Allow-Credentials", "true")
			}

			if r.Method == "OPTIONS" {
				handleOptionsRequest(w, cfg)
				w.WriteHeader(http.StatusOK)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
