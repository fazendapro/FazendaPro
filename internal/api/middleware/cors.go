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

			// Log para debug
			log.Printf("CORS: Origin=%s, Method=%s, Path=%s", origin, r.Method, r.URL.Path)
			log.Printf("CORS: AllowedOrigins=%v", cfg.CORS.AllowedOrigins)

			// Sempre definir headers de CORS para todas as requisições
			if len(cfg.CORS.AllowedOrigins) > 0 {
				// Se "*" está na lista, permitir qualquer origem
				allowAll := false
				for _, allowedOrigin := range cfg.CORS.AllowedOrigins {
					if allowedOrigin == "*" {
						allowAll = true
						break
					}
				}

				if allowAll {
					w.Header().Set("Access-Control-Allow-Origin", "*")
					log.Printf("CORS: Setting Access-Control-Allow-Origin to *")
				} else if origin != "" {
					// Verificar se a origem está na lista de permitidas
					originAllowed := false
					for _, allowedOrigin := range cfg.CORS.AllowedOrigins {
						if allowedOrigin == origin {
							w.Header().Set("Access-Control-Allow-Origin", origin)
							originAllowed = true
							log.Printf("CORS: Setting Access-Control-Allow-Origin to %s", origin)
							break
						}
					}
					if !originAllowed {
						log.Printf("CORS: Origin %s not allowed", origin)
					}
				}
			}

			// Configurar métodos permitidos (apenas para preflight)
			if r.Method == "OPTIONS" && len(cfg.CORS.AllowedMethods) > 0 {
				methods := strings.Join(cfg.CORS.AllowedMethods, ", ")
				w.Header().Set("Access-Control-Allow-Methods", methods)
				log.Printf("CORS: Setting Access-Control-Allow-Methods to %s", methods)
			}

			// Configurar headers permitidos (apenas para preflight)
			if r.Method == "OPTIONS" && len(cfg.CORS.AllowedHeaders) > 0 {
				headers := strings.Join(cfg.CORS.AllowedHeaders, ", ")
				w.Header().Set("Access-Control-Allow-Headers", headers)
				log.Printf("CORS: Setting Access-Control-Allow-Headers to %s", headers)
			}

			// Configurar headers expostos
			if len(cfg.CORS.ExposedHeaders) > 0 {
				exposedHeaders := strings.Join(cfg.CORS.ExposedHeaders, ", ")
				w.Header().Set("Access-Control-Expose-Headers", exposedHeaders)
				log.Printf("CORS: Setting Access-Control-Expose-Headers to %s", exposedHeaders)
			}

			// Configurar credenciais (não pode ser usado com "*")
			if cfg.CORS.AllowCredentials && origin != "" {
				w.Header().Set("Access-Control-Allow-Credentials", "true")
				log.Printf("CORS: Setting Access-Control-Allow-Credentials to true")
			}

			// Configurar max age (apenas para preflight)
			if r.Method == "OPTIONS" && cfg.CORS.MaxAge > 0 {
				maxAge := fmt.Sprintf("%d", cfg.CORS.MaxAge)
				w.Header().Set("Access-Control-Max-Age", maxAge)
				log.Printf("CORS: Setting Access-Control-Max-Age to %s", maxAge)
			}

			// Responder imediatamente para requisições OPTIONS (preflight)
			if r.Method == "OPTIONS" {
				log.Printf("CORS: Handling OPTIONS preflight request")
				w.WriteHeader(http.StatusOK)
				return
			}

			// Continuar para o próximo handler
			next.ServeHTTP(w, r)
		})
	}
}
