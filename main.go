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
	_ "github.com/fazendapro/FazendaPro-api/docs" // swagger docs
	"github.com/fazendapro/FazendaPro-api/internal/migrations"
	"github.com/fazendapro/FazendaPro-api/internal/repository"
	"github.com/fazendapro/FazendaPro-api/internal/routes"
	"github.com/getsentry/sentry-go"
)

// @title           FazendaPro API
// @version         1.0
// @description     API REST para gerenciamento de fazendas leiteiras
// @description     Sistema completo para gestão de animais, coletas de leite, reproduções, vendas e dívidas

// @contact.name   Suporte FazendaPro
// @contact.email  suporte@fazendapro.com

// @license.name  MIT
// @license.url   https://opensource.org/licenses/MIT

// @host      localhost:8080
// @BasePath  /

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Digite "Bearer" seguido de um espaço e depois o token JWT. Exemplo: "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."

func main() {
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

	if cfg.SentryDSN != "" {
		if err := sentry.Init(sentry.ClientOptions{
			Dsn: cfg.SentryDSN,
		}); err != nil {
			app.Logger.Printf("Sentry initialization failed: %v\n", err)
		} else {
			app.Logger.Println("Sentry inicializado com sucesso")
		}
	} else {
		app.Logger.Println("Sentry DSN não configurado, continuando sem Sentry")
	}

	db, err := repository.NewDatabase(cfg)
	if err != nil {
		app.Logger.Printf("WARNING: Erro ao conectar ao banco: %v", err)
		app.Logger.Println("Continuando sem conexão com banco de dados...")
	} else {
		defer db.Close()

		if err := migrations.RunMigrations(db.DB); err != nil {
			app.Logger.Printf("WARNING: Erro ao executar migrações: %v", err)
			app.Logger.Println("Continuando sem migrações...")
		}
	}

	var dbInstance *repository.Database
	if db != nil {
		dbInstance = db
	}
	r := routes.SetupRoutes(app, dbInstance, cfg)
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
