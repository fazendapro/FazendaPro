package repository

import (
	"fmt"
	"os"

	"github.com/fazendapro/FazendaPro-api/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Database struct {
	DB *gorm.DB
}

func NewDatabase(cfg *config.Config) (*Database, error) {
	sslMode := "disable"

	// Check if DB_SSL_MODE is explicitly set
	if os.Getenv("DB_SSL_MODE") != "" {
		sslMode = os.Getenv("DB_SSL_MODE")
	} else {
		env := os.Getenv("ENV")
		if env == "production" {
			sslMode = "require"
		} else {
			sslMode = "disable"
		}
	}

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		cfg.DBHost, cfg.User, cfg.Password, cfg.Name, cfg.DBPort, sslMode,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("erro ao conectar ao banco: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("erro ao configurar pool: %w", err)
	}
	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetMaxIdleConns(25)

	return &Database{DB: db}, nil
}

func (d *Database) Close() error {
	sqlDB, err := d.DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}
