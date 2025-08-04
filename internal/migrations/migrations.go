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
		return fmt.Errorf("error creating migrations table: %w", err)
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
		{"008_create_persons_table", createPersonsTable},
		{"009_update_users_table", updateUsersTable},
		{"011_update_users_with_person", updateUsersWithPerson},
		{"010_create_expenses_table", createExpensesTable},
		{"012_add_company_name", addCompanyName},
		{"013_add_farm_logo", addFarmLogo},
		{"014_add_animal_photo", addAnimalPhoto},
		{"015_update_animals_table", updateAnimalsTable},
	}

	for _, migration := range migrations {
		var existingMigration Migration
		if err := db.Where("name = ?", migration.name).First(&existingMigration).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				log.Printf("Executando migração: %s", migration.name)

				if err := migration.fn(db); err != nil {
					return fmt.Errorf("error executing migration %s: %w", migration.name, err)
				}

				if err := db.Create(&Migration{Name: migration.name}).Error; err != nil {
					return fmt.Errorf("error registering migration %s: %w", migration.name, err)
				}

				log.Printf("Migration %s executed successfully", migration.name)
			} else {
				return fmt.Errorf("error checking migration %s: %w", migration.name, err)
			}
		} else {
			log.Printf("Migration %s already executed", migration.name)
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

func createPersonsTable(db *gorm.DB) error {
	return db.AutoMigrate(&models.Person{})
}

func updateUsersTable(db *gorm.DB) error {
	var count int64
	db.Model(&models.User{}).Count(&count)

	if count > 0 {
		log.Printf("Users table has %d records. Executing auto migrate...", count)
	}

	return db.AutoMigrate(&models.User{})
}

func createExpensesTable(db *gorm.DB) error {
	return db.AutoMigrate(&models.Expense{})
}

func updateUsersWithPerson(db *gorm.DB) error {
	return db.AutoMigrate(&models.User{})
}

func addCompanyName(db *gorm.DB) error {
	return db.AutoMigrate(&models.Company{})
}

func addFarmLogo(db *gorm.DB) error {
	return db.AutoMigrate(&models.Farm{})
}

func addAnimalPhoto(db *gorm.DB) error {
	return db.AutoMigrate(&models.Animal{})
}

func updateAnimalsTable(db *gorm.DB) error {
	// Remove a coluna ear_tag_number antiga
	if db.Migrator().HasColumn(&models.Animal{}, "ear_tag_number") {
		if err := db.Migrator().DropColumn(&models.Animal{}, "ear_tag_number"); err != nil {
			return fmt.Errorf("error dropping ear_tag_number column: %w", err)
		}
	}

	// Remove a coluna age antiga
	if db.Migrator().HasColumn(&models.Animal{}, "age") {
		if err := db.Migrator().DropColumn(&models.Animal{}, "age"); err != nil {
			return fmt.Errorf("error dropping age column: %w", err)
		}
	}

	// Remove a coluna fertilization antiga (string)
	if db.Migrator().HasColumn(&models.Animal{}, "fertilization") {
		if err := db.Migrator().DropColumn(&models.Animal{}, "fertilization"); err != nil {
			return fmt.Errorf("error dropping fertilization column: %w", err)
		}
	}

	// Remove a coluna status antiga (string)
	if db.Migrator().HasColumn(&models.Animal{}, "status") {
		if err := db.Migrator().DropColumn(&models.Animal{}, "status"); err != nil {
			return fmt.Errorf("error dropping status column: %w", err)
		}
	}

	// Remove a coluna purpose antiga (string)
	if db.Migrator().HasColumn(&models.Animal{}, "purpose") {
		if err := db.Migrator().DropColumn(&models.Animal{}, "purpose"); err != nil {
			return fmt.Errorf("error dropping purpose column: %w", err)
		}
	}

	// Remove a coluna animal_type antiga
	if db.Migrator().HasColumn(&models.Animal{}, "animal_type") {
		if err := db.Migrator().DropColumn(&models.Animal{}, "animal_type"); err != nil {
			return fmt.Errorf("error dropping animal_type column: %w", err)
		}
	}

	// Executa o AutoMigrate para adicionar as novas colunas
	return db.AutoMigrate(&models.Animal{})
}

func RollbackMigrations(db *gorm.DB, steps int) error {
	var migrations []Migration
	if err := db.Order("id desc").Limit(steps).Find(&migrations).Error; err != nil {
		return fmt.Errorf("error searching migrations: %w", err)
	}

	for _, migration := range migrations {
		log.Printf("Reverting migration: %s", migration.Name)

		switch migration.Name {
		case "001_create_users_table":
			if err := db.Migrator().DropTable(&models.User{}); err != nil {
				return fmt.Errorf("error reverting migration %s: %w", migration.Name, err)
			}
		case "002_create_companies_table":
			if err := db.Migrator().DropTable(&models.Company{}); err != nil {
				return fmt.Errorf("error reverting migration %s: %w", migration.Name, err)
			}
		case "003_create_farms_table":
			if err := db.Migrator().DropTable(&models.Farm{}); err != nil {
				return fmt.Errorf("error reverting migration %s: %w", migration.Name, err)
			}
		case "004_create_animals_table":
			if err := db.Migrator().DropTable(&models.Animal{}); err != nil {
				return fmt.Errorf("error reverting migration %s: %w", migration.Name, err)
			}
		case "005_create_milk_collections_table":
			if err := db.Migrator().DropTable(&models.MilkCollection{}); err != nil {
				return fmt.Errorf("error reverting migration %s: %w", migration.Name, err)
			}
		case "006_create_reproductions_table":
			if err := db.Migrator().DropTable(&models.Reproduction{}); err != nil {
				return fmt.Errorf("error reverting migration %s: %w", migration.Name, err)
			}
		case "007_create_weights_table":
			if err := db.Migrator().DropTable(&models.Weight{}); err != nil {
				return fmt.Errorf("error reverting migration %s: %w", migration.Name, err)
			}
		case "008_create_persons_table":
			if err := db.Migrator().DropTable(&models.Person{}); err != nil {
				return fmt.Errorf("error reverting migration %s: %w", migration.Name, err)
			}
		case "009_update_users_table":
			if err := db.Migrator().DropTable(&models.User{}); err != nil {
				return fmt.Errorf("error reverting migration %s: %w", migration.Name, err)
			}
			if err := db.AutoMigrate(&models.User{}); err != nil {
				return fmt.Errorf("error recreating users table: %w", err)
			}
		case "010_create_expenses_table":
			if err := db.Migrator().DropTable(&models.Expense{}); err != nil {
				return fmt.Errorf("error reverting migration %s: %w", migration.Name, err)
			}
		case "011_update_users_with_person":
			if err := db.Migrator().DropTable(&models.User{}); err != nil {
				return fmt.Errorf("error reverting migration %s: %w", migration.Name, err)
			}
			if err := db.AutoMigrate(&models.User{}); err != nil {
				return fmt.Errorf("error recreating users table: %w", err)
			}
		case "012_add_company_name":
			if err := db.Migrator().DropColumn(&models.Company{}, "company_name"); err != nil {
				return fmt.Errorf("error reverting migration %s: %w", migration.Name, err)
			}
		case "013_add_farm_logo":
			if err := db.Migrator().DropColumn(&models.Farm{}, "logo"); err != nil {
				return fmt.Errorf("error reverting migration %s: %w", migration.Name, err)
			}
		case "014_add_animal_photo":
			if err := db.Migrator().DropColumn(&models.Animal{}, "photo"); err != nil {
				return fmt.Errorf("error reverting migration %s: %w", migration.Name, err)
			}
		case "015_update_animals_table":
			if err := db.AutoMigrate(&models.Animal{}); err != nil {
				return fmt.Errorf("error reverting animals table: %w", err)
			}
		}

		if err := db.Delete(&migration).Error; err != nil {
			return fmt.Errorf("error removing migration record %s: %w", migration.Name, err)
		}

		log.Printf("Migration %s reverted successfully", migration.Name)
	}

	return nil
}
