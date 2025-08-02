package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/fazendapro/FazendaPro-api/cmd/app"
	"github.com/fazendapro/FazendaPro-api/config"
	"github.com/fazendapro/FazendaPro-api/internal/repository"
	"github.com/fazendapro/FazendaPro-api/internal/routes"
)

func main() {
	var port int
	flag.IntVar(&port, "port", 8080, "Porta do servidor")
	flag.Parse()

	app, err := app.NewApplication()
	if err != nil {
		panic(err)
	}

	app.Logger.Println("Iniciando aplicação...")

	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Erro ao carregar configuração:", err)
	}

	db, err := repository.NewDatabase(cfg)
	if err != nil {
		log.Fatal("Erro ao conectar ao banco:", err)
	}
	defer db.Close()

	// userRepo := repository.NewUserRepository(db)
	// userService := service.NewUserService(*userRepo)
	// userHandler := handlers.NewUserHandler(userService)

	r := routes.SetupRoutes(app)
	server := http.Server{
		Addr:         fmt.Sprintf(":%d", port),
		Handler:      r,
		IdleTimeout:  30 * time.Second,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	app.Logger.Printf("Servidor rodando na porta %d", port)

	err = server.ListenAndServe()
	if err != nil {
		app.Logger.Fatal("Erro ao iniciar servidor:", err)
	}
}
