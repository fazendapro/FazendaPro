package migrations

import (
	"fmt"
	"log"
	"time"

	"github.com/fazendapro/FazendaPro-api/internal/models"
	"gorm.io/gorm"
)

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
		{"013_create_refresh_tokens_table", createRefreshTokensTable},
		{"013_add_farm_logo", addFarmLogo},
		{"014_add_animal_photo", addAnimalPhoto},
		{"015_update_animals_table", updateAnimalsTable},
		{"016_update_reproductions_table", updateReproductionsTable},
		{"017_seed_initial_data", seedInitialData},
		{"018_create_user_farms_table", createUserFarmsTable},
		{"019_migrate_users_to_user_farms", migrateUsersToUserFarms},
		{"020_create_sales_table", createSalesTable},
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
	if db.Migrator().HasColumn(&models.Animal{}, "ear_tag_number") {
		if err := db.Migrator().DropColumn(&models.Animal{}, "ear_tag_number"); err != nil {
			return fmt.Errorf("error dropping ear_tag_number column: %w", err)
		}
	}

	if db.Migrator().HasColumn(&models.Animal{}, "age") {
		if err := db.Migrator().DropColumn(&models.Animal{}, "age"); err != nil {
			return fmt.Errorf("error dropping age column: %w", err)
		}
	}

	if db.Migrator().HasColumn(&models.Animal{}, "fertilization") {
		if err := db.Migrator().DropColumn(&models.Animal{}, "fertilization"); err != nil {
			return fmt.Errorf("error dropping fertilization column: %w", err)
		}
	}

	if db.Migrator().HasColumn(&models.Animal{}, "status") {
		if err := db.Migrator().DropColumn(&models.Animal{}, "status"); err != nil {
			return fmt.Errorf("error dropping status column: %w", err)
		}
	}

	if db.Migrator().HasColumn(&models.Animal{}, "purpose") {
		if err := db.Migrator().DropColumn(&models.Animal{}, "purpose"); err != nil {
			return fmt.Errorf("error dropping purpose column: %w", err)
		}
	}

	if db.Migrator().HasColumn(&models.Animal{}, "animal_type") {
		if err := db.Migrator().DropColumn(&models.Animal{}, "animal_type"); err != nil {
			return fmt.Errorf("error dropping animal_type column: %w", err)
		}
	}

	return db.AutoMigrate(&models.Animal{})
}

func updateReproductionsTable(db *gorm.DB) error {
	if db.Migrator().HasColumn(&models.Reproduction{}, "type") {
		if err := db.Migrator().DropColumn(&models.Reproduction{}, "type"); err != nil {
			return fmt.Errorf("error dropping type column: %w", err)
		}
	}

	if db.Migrator().HasColumn(&models.Reproduction{}, "pregnancy_date") {
		if err := db.Migrator().DropColumn(&models.Reproduction{}, "pregnancy_date"); err != nil {
			return fmt.Errorf("error dropping pregnancy_date column: %w", err)
		}
	}

	if db.Migrator().HasColumn(&models.Reproduction{}, "situation") {
		if err := db.Migrator().DropColumn(&models.Reproduction{}, "situation"); err != nil {
			return fmt.Errorf("error dropping situation column: %w", err)
		}
	}

	return db.AutoMigrate(&models.Reproduction{})
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
		case "013_create_refresh_tokens_table":
			if err := db.Migrator().DropTable(&models.RefreshToken{}); err != nil {
				return fmt.Errorf("error reverting migration %s: %w", migration.Name, err)
			}
		case "016_update_reproductions_table":
			if err := db.AutoMigrate(&models.Reproduction{}); err != nil {
				return fmt.Errorf("error reverting reproductions table: %w", err)
			}
		}

		if err := db.Delete(&migration).Error; err != nil {
			return fmt.Errorf("error removing migration record %s: %w", migration.Name, err)
		}

		log.Printf("Migration %s reverted successfully", migration.Name)
	}

	return nil
}

func createRefreshTokensTable(db *gorm.DB) error {
	return db.AutoMigrate(&models.RefreshToken{})
}

func seedInitialData(db *gorm.DB) error {
	var companyCount int64
	db.Model(&models.Company{}).Count(&companyCount)

	if companyCount > 0 {
		log.Printf("Dados iniciais já existem, pulando seed")
		return nil
	}

	company := &models.Company{
		CompanyName: "FazendaPro Demo",
	}
	if err := db.Create(company).Error; err != nil {
		return fmt.Errorf("error creating company: %w", err)
	}

	farm := &models.Farm{
		CompanyID: company.ID,
		Logo:      "",
	}
	if err := db.Create(farm).Error; err != nil {
		return fmt.Errorf("error creating farm: %w", err)
	}

	log.Printf("Dados iniciais criados: Company ID=%d, Farm ID=%d", company.ID, farm.ID)
	return nil
}

func createUserFarmsTable(db *gorm.DB) error {
	return db.AutoMigrate(&models.UserFarm{})
}

func migrateUsersToUserFarms(db *gorm.DB) error {
	var users []models.User
	if err := db.Find(&users).Error; err != nil {
		return fmt.Errorf("error finding users: %w", err)
	}

	for _, user := range users {
		userFarm := &models.UserFarm{
			UserID:    user.ID,
			FarmID:    user.FarmID,
			IsPrimary: true,
		}

		var existingUserFarm models.UserFarm
		if err := db.Where("user_id = ? AND farm_id = ?", user.ID, user.FarmID).First(&existingUserFarm).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				if err := db.Create(userFarm).Error; err != nil {
					return fmt.Errorf("error creating user farm for user %d: %w", user.ID, err)
				}
				log.Printf("Migrated user %d to farm %d", user.ID, user.FarmID)
			} else {
				return fmt.Errorf("error checking existing user farm: %w", err)
			}
		} else {
			log.Printf("User %d already has farm %d, skipping", user.ID, user.FarmID)
		}
	}

	log.Printf("Migration completed: %d users processed", len(users))
	return nil
}

func createSalesTable(db *gorm.DB) error {
	log.Printf("Creating sales table...")

	if err := db.AutoMigrate(&models.Sale{}); err != nil {
		return fmt.Errorf("error creating sales table: %w", err)
	}

	log.Printf("Sales table created successfully")
	return nil
}
