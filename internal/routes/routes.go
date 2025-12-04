package routes

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/fazendapro/FazendaPro-api/cmd/app"
	"github.com/fazendapro/FazendaPro-api/config"
	"github.com/fazendapro/FazendaPro-api/internal/api/handlers"
	"github.com/fazendapro/FazendaPro-api/internal/api/middleware"
	"github.com/fazendapro/FazendaPro-api/internal/cache"
	"github.com/fazendapro/FazendaPro-api/internal/models"
	"github.com/fazendapro/FazendaPro-api/internal/repository"
	"github.com/fazendapro/FazendaPro-api/internal/service"
	"github.com/getsentry/sentry-go"
	sentryhttp "github.com/getsentry/sentry-go/http"
	"github.com/go-chi/chi/v5"
	httpSwagger "github.com/swaggo/http-swagger"
)

func SetupRoutes(app *app.Application, db *repository.Database, cfg *config.Config) *chi.Mux {
	r := chi.NewRouter()

	if cfg.SentryDSN != "" {
		sentryHandler := sentryhttp.New(sentryhttp.Options{
			Repanic: true,
		})
		r.Use(func(next http.Handler) http.Handler {
			return sentryHandler.Handle(next)
		})
	}

	r.Use(middleware.CORSMiddleware(cfg))

	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			app.Logger.Printf("Request: %s %s", r.Method, r.URL.Path)
			next.ServeHTTP(w, r)
		})
	})

	r.Get("/health", app.HealthCheck)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("FazendaPro API is running!"))
	})

	r.Get("/swagger/*", func(w http.ResponseWriter, r *http.Request) {
		scheme := "http"
		if r.TLS != nil || r.Header.Get("X-Forwarded-Proto") == "https" {
			scheme = "https"
		}

		host := r.Host
		if host == "" {
			host = "localhost:8080"
		}

		docURL := fmt.Sprintf("%s://%s/swagger/doc.json", scheme, host)

		httpSwagger.Handler(httpSwagger.URL(docURL)).ServeHTTP(w, r)
	})

	r.Get("/test-error", func(w http.ResponseWriter, r *http.Request) {
		err := errors.New("erro de teste do Sentry - Rota /test-error")
		sentry.CaptureException(err)
		sentry.Flush(2)

		w.Header().Set(handlers.HeaderContentType, handlers.ContentTypeJSON)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": "Erro de teste enviado para o Sentry", "message": "Verifique o dashboard do Sentry para ver este erro"}`))
	})

	r.Get("/test-panic", func(w http.ResponseWriter, r *http.Request) {
		panic("Panic de teste do Sentry - Rota /test-panic")
	})

	var cacheClient cache.CacheInterface
	if db != nil && db.DB != nil {
		memcachedServer := fmt.Sprintf("%s:%s", cfg.MemcachedHost, cfg.MemcachedPort)
		cacheClient = cache.NewMemcacheClient(memcachedServer)
		app.Logger.Printf("Cache Memcached inicializado em %s", memcachedServer)
	}

	if db != nil && db.DB != nil {
		repoFactory := repository.NewRepositoryFactory(db, cacheClient)
		serviceFactory := service.NewServiceFactory(repoFactory)
		debtService := serviceFactory.CreateDebtService()
		debtHandler := handlers.NewDebtHandler(debtService)

		r.Route("/debts", func(r chi.Router) {
			r.Post("/", debtHandler.CreateDebt)
			r.Get("/", debtHandler.GetDebts)
			r.Delete("/{id}", debtHandler.DeleteDebt)
			r.Get("/total-by-person", debtHandler.GetTotalByPerson)
		})
	}

	r.Post("/init-data", func(w http.ResponseWriter, r *http.Request) {
		if db == nil || db.DB == nil {
			http.Error(w, "Database not available", http.StatusInternalServerError)
			return
		}

		var companyCount int64
		db.DB.Model(&models.Company{}).Count(&companyCount)

		if companyCount > 0 {
			w.Header().Set(handlers.HeaderContentType, handlers.ContentTypeJSON)
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"success": true, "message": "Dados já existem"}`))
			return
		}

		company := &models.Company{
			CompanyName: "FazendaPro Demo",
		}
		if err := db.DB.Create(company).Error; err != nil {
			http.Error(w, "Error creating company: "+err.Error(), http.StatusInternalServerError)
			return
		}

		farm := &models.Farm{
			CompanyID: company.ID,
			Logo:      "",
		}
		if err := db.DB.Create(farm).Error; err != nil {
			http.Error(w, "Error creating farm: "+err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set(handlers.HeaderContentType, handlers.ContentTypeJSON)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"success": true, "message": "Dados iniciais criados", "company_id": ` + fmt.Sprintf("%d", company.ID) + `, "farm_id": ` + fmt.Sprintf("%d", farm.ID) + `}`))
	})

	if db != nil && db.DB != nil {
		app.Logger.Println("Database conectado - configurando rotas de dados")
		repoFactory := repository.NewRepositoryFactory(db, cacheClient)
		serviceFactory := service.NewServiceFactory(repoFactory)

		r.Route("/api/v1", func(r chi.Router) {
			userService := serviceFactory.CreateUserService()
			userHandler := handlers.NewUserHandler(userService)
			refreshTokenRepo := repoFactory.CreateRefreshTokenRepository()
			authHandler := handlers.NewAuthHandler(userService, refreshTokenRepo, cfg.JWTSecret)

			r.Route("/auth", func(r chi.Router) {
				r.Post("/login", authHandler.Login)
				r.Post("/register", authHandler.Register)
				r.Post("/refresh", authHandler.RefreshToken)
				r.Post("/logout", authHandler.Logout)
			})

			r.Route("/users", func(r chi.Router) {
				r.Use(middleware.Auth(cfg.JWTSecret))
				r.Post("/", userHandler.CreateUser)
				r.Get("/", userHandler.GetUser)
			})

			farmSelectionHandler := handlers.NewFarmSelectionHandler(userService, cfg.JWTSecret)
			r.Route("/farms", func(r chi.Router) {
				r.Use(middleware.Auth(cfg.JWTSecret))
				r.Get("/user", farmSelectionHandler.GetUserFarms)
				r.Post("/select", farmSelectionHandler.SelectFarm)
			})

			animalService := serviceFactory.CreateAnimalService()
			weightService := serviceFactory.CreateWeightService()
			animalHandler := handlers.NewAnimalHandlerWithWeight(animalService, weightService)

			r.Route("/animals", func(r chi.Router) {
				r.Use(middleware.Auth(cfg.JWTSecret))
				r.Post("/", animalHandler.CreateAnimal)
				r.Get("/", animalHandler.GetAnimal)
				r.Get("/farm", animalHandler.GetAnimalsByFarm)
				r.Get("/sex", animalHandler.GetAnimalsBySex)
				r.Put("/", animalHandler.UpdateAnimal)
				r.Delete("/", animalHandler.DeleteAnimal)
				r.Post("/photo", animalHandler.UploadAnimalPhoto)
			})

			milkCollectionService := serviceFactory.CreateMilkCollectionService()
			milkCollectionHandler := handlers.NewMilkCollectionHandler(milkCollectionService)

			r.Route("/milk-collections", func(r chi.Router) {
				r.Use(middleware.Auth(cfg.JWTSecret))
				r.Post("/", milkCollectionHandler.CreateMilkCollection)
				r.Put("/{id}", milkCollectionHandler.UpdateMilkCollection)
				r.Get("/farm/{farmId}", milkCollectionHandler.GetMilkCollectionsByFarmID)
				r.Get("/animal/{animalId}", milkCollectionHandler.GetMilkCollectionsByAnimalID)
				r.Get("/top-producers", milkCollectionHandler.GetTopMilkProducers)
			})

			reproductionService := serviceFactory.CreateReproductionService()
			reproductionHandler := handlers.NewReproductionHandler(reproductionService)

			r.Route("/reproductions", func(r chi.Router) {
				r.Use(middleware.Auth(cfg.JWTSecret))
				r.Post("/", reproductionHandler.CreateReproduction)
				r.Get("/", reproductionHandler.GetReproduction)
				r.Get("/animal", reproductionHandler.GetReproductionByAnimal)
				r.Get("/farm", reproductionHandler.GetReproductionsByFarm)
				r.Get("/phase", reproductionHandler.GetReproductionsByPhase)
				r.Get("/next-to-calve", reproductionHandler.GetNextToCalve)
				r.Put("/", reproductionHandler.UpdateReproduction)
				r.Put("/phase", reproductionHandler.UpdateReproductionPhase)
				r.Delete("/", reproductionHandler.DeleteReproduction)
			})

			farmService := serviceFactory.CreateFarmService()
			farmHandler := handlers.NewFarmHandler(farmService)

			r.Route("/farm", func(r chi.Router) {
				r.Use(middleware.Auth(cfg.JWTSecret))
				r.Get("/", farmHandler.GetFarm)
				r.Put("/", farmHandler.UpdateFarm)
			})

			saleService := serviceFactory.CreateSaleService()
			saleHandler := handlers.NewSaleChiHandler(saleService)

			r.Route("/sales", func(r chi.Router) {
				r.Use(middleware.Auth(cfg.JWTSecret))
				r.Post("/", saleHandler.CreateSale)
				r.Get("/", saleHandler.GetSalesByFarm)
				r.Get("/history", saleHandler.GetSalesHistory)
				r.Get("/monthly-stats", saleHandler.GetMonthlySalesStats)
				r.Get("/monthly-data", saleHandler.GetMonthlySalesAndPurchases)
				r.Get("/overview", saleHandler.GetOverviewStats)
				r.Get("/date-range", saleHandler.GetSalesByDateRange)
				r.Get("/{id}", saleHandler.GetSaleByID)
				r.Put("/{id}", saleHandler.UpdateSale)
				r.Delete("/{id}", saleHandler.DeleteSale)
			})

			r.Route("/animals/{animal_id}/sales", func(r chi.Router) {
				r.Use(middleware.Auth(cfg.JWTSecret))
				r.Get("/", saleHandler.GetSalesByAnimal)
			})

			vaccineService := serviceFactory.CreateVaccineService()
			vaccineHandler := handlers.NewVaccineHandler(vaccineService)

			r.Route("/vaccines", func(r chi.Router) {
				r.Use(middleware.Auth(cfg.JWTSecret))
				r.Post("/", vaccineHandler.CreateVaccine)
				r.Get("/farm/{farmId}", vaccineHandler.GetVaccinesByFarmID)
				r.Get("/{id}", vaccineHandler.GetVaccineByID)
				r.Put("/{id}", vaccineHandler.UpdateVaccine)
				r.Delete("/{id}", vaccineHandler.DeleteVaccine)
			})

			vaccineApplicationService := serviceFactory.CreateVaccineApplicationService()
			vaccineApplicationHandler := handlers.NewVaccineApplicationHandler(vaccineApplicationService)

			r.Route("/vaccine-applications", func(r chi.Router) {
				r.Use(middleware.Auth(cfg.JWTSecret))
				r.Post("/", vaccineApplicationHandler.CreateVaccineApplication)
				r.Get("/farm/{farmId}", vaccineApplicationHandler.GetVaccineApplicationsByFarmID)
				r.Get("/animal/{animalId}", vaccineApplicationHandler.GetVaccineApplicationsByAnimalID)
				r.Get("/{id}", vaccineApplicationHandler.GetVaccineApplicationByID)
				r.Put("/{id}", vaccineApplicationHandler.UpdateVaccineApplication)
				r.Delete("/{id}", vaccineApplicationHandler.DeleteVaccineApplication)
			})

			weightHandler := handlers.NewWeightHandler(weightService)

			r.Route("/weights", func(r chi.Router) {
				r.Use(middleware.Auth(cfg.JWTSecret))
				r.Post("/", weightHandler.CreateOrUpdateWeight)
				r.Put("/", weightHandler.UpdateWeight)
				r.Get("/farm/{farmId}", weightHandler.GetWeightsByFarm)
				r.Get("/animal/{animalId}", weightHandler.GetWeightByAnimal)
			})
		})

		app.Logger.Println("Rotas de animais configuradas: /api/v1/animals/farm")
		app.Logger.Println("Rotas de milk collections configuradas: /api/v1/milk-collections")
		app.Logger.Println("Rotas de reprodução configuradas: /api/v1/reproductions")
		app.Logger.Println("Rotas de vacinas configuradas: /api/v1/vaccines")
		app.Logger.Println("Rotas de aplicações de vacinas configuradas: /api/v1/vaccine-applications")
		app.Logger.Println("Rotas de pesos configuradas: /api/v1/weights")
	} else {
		app.Logger.Println("Database não disponível - rotas de dados desabilitadas")
	}

	fmt.Println("Todas as rotas configuradas com sucesso")
	return r
}
