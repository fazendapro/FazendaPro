package routes

import (
	"fmt"
	"net/http"

	"github.com/fazendapro/FazendaPro-api/cmd/app"
	"github.com/fazendapro/FazendaPro-api/config"
	"github.com/fazendapro/FazendaPro-api/internal/api/handlers"
	"github.com/fazendapro/FazendaPro-api/internal/api/middleware"
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

	if db != nil && db.DB != nil {
		app.Logger.Println("Database conectado - configurando rotas de dados")
		repoFactory := repository.NewRepositoryFactory(db)
		serviceFactory := service.NewServiceFactory(repoFactory)

		userService := serviceFactory.CreateUserService()
		userHandler := handlers.NewUserHandler(userService)

		r.Route("/users", func(r chi.Router) {
			r.Post("/", userHandler.CreateUser)
			r.Get("/", userHandler.GetUser)
		})

		animalService := serviceFactory.CreateAnimalService()
		animalHandler := handlers.NewAnimalHandler(animalService)

		r.Route("/animals", func(r chi.Router) {
			r.Post("/", animalHandler.CreateAnimal)
			r.Get("/", animalHandler.GetAnimal)
			r.Get("/farm", animalHandler.GetAnimalsByFarm)
			r.Put("/", animalHandler.UpdateAnimal)
			r.Delete("/", animalHandler.DeleteAnimal)
		})

		app.Logger.Println("Rotas de animais configuradas: /animals/farm")
	} else {
		app.Logger.Println("Database não disponível - rotas de dados desabilitadas")
	}

	fmt.Println("Todas as rotas configuradas com sucesso")
	return r
}
