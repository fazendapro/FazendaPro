package routes

import (
	"fmt"

	"github.com/fazendapro/FazendaPro-api/cmd/app"
	"github.com/fazendapro/FazendaPro-api/internal/api/handlers"
	"github.com/fazendapro/FazendaPro-api/internal/repository"
	"github.com/fazendapro/FazendaPro-api/internal/service"
	"github.com/go-chi/chi/v5"
)

func SetupRoutes(app *app.Application, db *repository.Database) *chi.Mux {
	r := chi.NewRouter()

	fmt.Println("Configurando rotas...")

	r.Get("/health", app.HealthCheck)
	fmt.Println("Rota /health configurada")

	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo)
	userHandler := handlers.NewUserHandler(userService)

	r.Route("/users", func(r chi.Router) {
		r.Post("/", userHandler.CreateUser)
		r.Get("/", userHandler.GetUser)
		fmt.Println("Rotas /users configuradas")
	})

	fmt.Println("Todas as rotas configuradas com sucesso")
	return r
}
