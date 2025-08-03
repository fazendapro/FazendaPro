package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/fazendapro/FazendaPro-api/cmd/app"
	"github.com/fazendapro/FazendaPro-api/config"
	"github.com/fazendapro/FazendaPro-api/internal/migrations"
	"github.com/fazendapro/FazendaPro-api/internal/repository"
	"github.com/fazendapro/FazendaPro-api/internal/routes"
)

func main() {
	// Verificar se é um comando de migração
	if len(os.Args) > 1 && os.Args[1] == "migrate" {
		runMigrations()
		return
	}

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

	// Executar migrações automaticamente na inicialização
	if err := migrations.RunMigrations(db.DB); err != nil {
		log.Fatal("Erro ao executar migrações:", err)
	}

	r := routes.SetupRoutes(app, db)
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

func runMigrations() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Erro ao carregar configuração:", err)
	}

	db, err := repository.NewDatabase(cfg)
	if err != nil {
		log.Fatal("Erro ao conectar ao banco:", err)
	}
	defer db.Close()

	log.Println("Executando migrações...")
	if err := migrations.RunMigrations(db.DB); err != nil {
		log.Fatal("Erro ao executar migrações:", err)
	}
	log.Println("Migrações executadas com sucesso!")
}
