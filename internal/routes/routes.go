package routes

import (
	"fmt"
	"net/http"

	"github.com/fazendapro/FazendaPro-api/cmd/app"
	"github.com/fazendapro/FazendaPro-api/config"
	"github.com/fazendapro/FazendaPro-api/internal/api/handlers"
	"github.com/fazendapro/FazendaPro-api/internal/api/middleware"
	"github.com/fazendapro/FazendaPro-api/internal/models"
	"github.com/fazendapro/FazendaPro-api/internal/repository"
	"github.com/fazendapro/FazendaPro-api/internal/service"
	"github.com/go-chi/chi/v5"
)

func SetupRoutes(app *app.Application, db *repository.Database, cfg *config.Config) *chi.Mux {
	r := chi.NewRouter()

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

	r.Post("/init-data", func(w http.ResponseWriter, r *http.Request) {
		if db == nil || db.DB == nil {
			http.Error(w, "Database not available", http.StatusInternalServerError)
			return
		}

		var companyCount int64
		db.DB.Model(&models.Company{}).Count(&companyCount)

		if companyCount > 0 {
			w.Header().Set("Content-Type", "application/json")
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

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"success": true, "message": "Dados iniciais criados", "company_id": ` + fmt.Sprintf("%d", company.ID) + `, "farm_id": ` + fmt.Sprintf("%d", farm.ID) + `}`))
	})

	if db != nil && db.DB != nil {
		app.Logger.Println("Database conectado - configurando rotas de dados")
		repoFactory := repository.NewRepositoryFactory(db)
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
			animalHandler := handlers.NewAnimalHandler(animalService)

			r.Route("/animals", func(r chi.Router) {
				r.Use(middleware.Auth(cfg.JWTSecret))
				r.Post("/", animalHandler.CreateAnimal)
				r.Get("/", animalHandler.GetAnimal)
				r.Get("/farm", animalHandler.GetAnimalsByFarm)
				r.Put("/", animalHandler.UpdateAnimal)
				r.Delete("/", animalHandler.DeleteAnimal)
			})

			milkCollectionService := serviceFactory.CreateMilkCollectionService()
			milkCollectionHandler := handlers.NewMilkCollectionHandler(milkCollectionService)

			r.Route("/milk-collections", func(r chi.Router) {
				r.Use(middleware.Auth(cfg.JWTSecret))
				r.Post("/", milkCollectionHandler.CreateMilkCollection)
				r.Put("/{id}", milkCollectionHandler.UpdateMilkCollection)
				r.Get("/farm/{farmId}", milkCollectionHandler.GetMilkCollectionsByFarmID)
				r.Get("/animal/{animalId}", milkCollectionHandler.GetMilkCollectionsByAnimalID)
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
				r.Put("/", reproductionHandler.UpdateReproduction)
				r.Put("/phase", reproductionHandler.UpdateReproductionPhase)
				r.Delete("/", reproductionHandler.DeleteReproduction)
			})
		})

		app.Logger.Println("Rotas de animais configuradas: /api/v1/animals/farm")
		app.Logger.Println("Rotas de milk collections configuradas: /api/v1/milk-collections")
		app.Logger.Println("Rotas de reprodução configuradas: /api/v1/reproductions")
	} else {
		app.Logger.Println("Database não disponível - rotas de dados desabilitadas")
	}

	fmt.Println("Todas as rotas configuradas com sucesso")
	return r
}
