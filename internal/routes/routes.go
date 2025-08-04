package routes

import (
	"fmt"

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

	r.Get("/health", app.HealthCheck)

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

	fmt.Println("Todas as rotas configuradas com sucesso")
	return r
}
