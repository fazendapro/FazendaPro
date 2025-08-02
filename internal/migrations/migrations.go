package migrations

import (
	"fmt"
	"log"
	"time"

	"github.com/fazendapro/FazendaPro-api/internal/models"
	"gorm.io/gorm"
)

// Migration representa uma migração
type Migration struct {
	ID        uint   `gorm:"primaryKey"`
	Name      string `gorm:"uniqueIndex"`
	CreatedAt time.Time
}

func RunMigrations(db *gorm.DB) error {
	if err := db.AutoMigrate(&Migration{}); err != nil {
		return fmt.Errorf("erro ao criar tabela de migrações: %w", err)
	}

	migrations := []struct {
		name string
		fn   func(*gorm.DB) error
	}{
		{"001_create_users_table", createUsersTable},
		{"002_create_companies_table", createCompaniesTable},
		{"003_create_farms_table", createFarmsTable},
		{"004_create_animals_table", createAnimalsTable},
		{"005_create_milk_collections_table", createMilkCollectionsTable},
		{"006_create_reproductions_table", createReproductionsTable},
		{"007_create_weights_table", createWeightsTable},
	}

	for _, migration := range migrations {
		var existingMigration Migration
		if err := db.Where("name = ?", migration.name).First(&existingMigration).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				log.Printf("Executando migração: %s", migration.name)

				if err := migration.fn(db); err != nil {
					return fmt.Errorf("erro ao executar migração %s: %w", migration.name, err)
				}

				if err := db.Create(&Migration{Name: migration.name}).Error; err != nil {
					return fmt.Errorf("erro ao registrar migração %s: %w", migration.name, err)
				}

				log.Printf("Migração %s executada com sucesso", migration.name)
			} else {
				return fmt.Errorf("erro ao verificar migração %s: %w", migration.name, err)
			}
		} else {
			log.Printf("Migração %s já foi executada", migration.name)
		}
	}

	return nil
}

func createUsersTable(db *gorm.DB) error {
	return db.AutoMigrate(&models.User{})
}

func createCompaniesTable(db *gorm.DB) error {
	return db.AutoMigrate(&models.Company{})
}

func createFarmsTable(db *gorm.DB) error {
	return db.AutoMigrate(&models.Farm{})
}

func createAnimalsTable(db *gorm.DB) error {
	return db.AutoMigrate(&models.Animal{})
}

func createMilkCollectionsTable(db *gorm.DB) error {
	return db.AutoMigrate(&models.MilkCollection{})
}

func createReproductionsTable(db *gorm.DB) error {
	return db.AutoMigrate(&models.Reproduction{})
}

func createWeightsTable(db *gorm.DB) error {
	return db.AutoMigrate(&models.Weight{})
}

func RollbackMigrations(db *gorm.DB, steps int) error {
	var migrations []Migration
	if err := db.Order("id desc").Limit(steps).Find(&migrations).Error; err != nil {
		return fmt.Errorf("erro ao buscar migrações: %w", err)
	}

	for _, migration := range migrations {
		log.Printf("Revertendo migração: %s", migration.Name)

		switch migration.Name {
		case "001_create_users_table":
			if err := db.Migrator().DropTable(&models.User{}); err != nil {
				return fmt.Errorf("erro ao reverter migração %s: %w", migration.Name, err)
			}
		case "002_create_companies_table":
			if err := db.Migrator().DropTable(&models.Company{}); err != nil {
				return fmt.Errorf("erro ao reverter migração %s: %w", migration.Name, err)
			}
		case "003_create_farms_table":
			if err := db.Migrator().DropTable(&models.Farm{}); err != nil {
				return fmt.Errorf("erro ao reverter migração %s: %w", migration.Name, err)
			}
		case "004_create_animals_table":
			if err := db.Migrator().DropTable(&models.Animal{}); err != nil {
				return fmt.Errorf("erro ao reverter migração %s: %w", migration.Name, err)
			}
		case "005_create_milk_collections_table":
			if err := db.Migrator().DropTable(&models.MilkCollection{}); err != nil {
				return fmt.Errorf("erro ao reverter migração %s: %w", migration.Name, err)
			}
		case "006_create_reproductions_table":
			if err := db.Migrator().DropTable(&models.Reproduction{}); err != nil {
				return fmt.Errorf("erro ao reverter migração %s: %w", migration.Name, err)
			}
		case "007_create_weights_table":
			if err := db.Migrator().DropTable(&models.Weight{}); err != nil {
				return fmt.Errorf("erro ao reverter migração %s: %w", migration.Name, err)
			}
		}

		if err := db.Delete(&migration).Error; err != nil {
			return fmt.Errorf("erro ao remover registro da migração %s: %w", migration.Name, err)
		}

		log.Printf("Migração %s revertida com sucesso", migration.Name)
	}

	return nil
}
