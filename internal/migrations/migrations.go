package migrations

import (
	"fmt"
	"log"
	"time"

	"github.com/fazendapro/FazendaPro-api/internal/models"
	"gorm.io/gorm"
)

const ErrRevertingMigration = "error reverting migration %s: %w"

type Migration struct {
	ID        uint   `gorm:"primaryKey"`
	Name      string `gorm:"uniqueIndex"`
	CreatedAt time.Time
}

func executeMigration(db *gorm.DB, name string, fn func(*gorm.DB) error) error {
	var existingMigration Migration
	if err := db.Where("name = ?", name).First(&existingMigration).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			log.Printf("Executando migração: %s", name)

			if err := fn(db); err != nil {
				return fmt.Errorf("error executing migration %s: %w", name, err)
			}

			if err := db.Create(&Migration{Name: name}).Error; err != nil {
				return fmt.Errorf("error registering migration %s: %w", name, err)
			}

			log.Printf("Migration %s executed successfully", name)
		} else {
			return fmt.Errorf("error checking migration %s: %w", name, err)
		}
	} else {
		log.Printf("Migration %s already executed", name)
	}
	return nil
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
		{"021_create_debts_table", createDebtsTable},
		{"022_create_vaccines_table", createVaccinesTable},
		{"023_create_vaccine_applications_table", createVaccineApplicationsTable},
		{"024_add_farm_language", addFarmLanguage},
	}

	for _, migration := range migrations {
		if err := executeMigration(db, migration.name, migration.fn); err != nil {
			return err
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

func dropColumnIfExists(db *gorm.DB, model interface{}, columnName string) error {
	if db.Migrator().HasColumn(model, columnName) {
		if err := db.Migrator().DropColumn(model, columnName); err != nil {
			return fmt.Errorf("error dropping %s column: %w", columnName, err)
		}
	}
	return nil
}

func updateAnimalsTable(db *gorm.DB) error {
	columnsToDrop := []string{"ear_tag_number", "age", "fertilization", "status", "purpose", "animal_type"}
	for _, column := range columnsToDrop {
		if err := dropColumnIfExists(db, &models.Animal{}, column); err != nil {
			return err
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

type rollbackFunc func(*gorm.DB, string) error

func revertDropTable(db *gorm.DB, model interface{}, migrationName string) error {
	if err := db.Migrator().DropTable(model); err != nil {
		return fmt.Errorf(ErrRevertingMigration, migrationName, err)
	}
	return nil
}

func revertDropColumn(db *gorm.DB, model interface{}, columnName, migrationName string) error {
	if err := db.Migrator().DropColumn(model, columnName); err != nil {
		return fmt.Errorf(ErrRevertingMigration, migrationName, err)
	}
	return nil
}

func revertRecreateUsersTable(db *gorm.DB, migrationName string) error {
	if err := db.Migrator().DropTable(&models.User{}); err != nil {
		return fmt.Errorf(ErrRevertingMigration, migrationName, err)
	}
	if err := db.AutoMigrate(&models.User{}); err != nil {
		return fmt.Errorf("error recreating users table: %w", err)
	}
	return nil
}

func revertAutoMigrate(db *gorm.DB, model interface{}, migrationName string) error {
	if err := db.AutoMigrate(model); err != nil {
		return fmt.Errorf("error reverting %s table: %w", migrationName, err)
	}
	return nil
}

func revertMigration(db *gorm.DB, migration Migration, rollbackFuncs map[string]rollbackFunc) error {
	log.Printf("Reverting migration: %s", migration.Name)

	rollback, exists := rollbackFuncs[migration.Name]
	if !exists {
		log.Printf("No rollback function found for migration: %s", migration.Name)
		return nil
	}

	if err := rollback(db, migration.Name); err != nil {
		return err
	}

	if err := db.Delete(&migration).Error; err != nil {
		return fmt.Errorf("error removing migration record %s: %w", migration.Name, err)
	}

	log.Printf("Migration %s reverted successfully", migration.Name)
	return nil
}

func RollbackMigrations(db *gorm.DB, steps int) error {
	var migrations []Migration
	if err := db.Order("id desc").Limit(steps).Find(&migrations).Error; err != nil {
		return fmt.Errorf("error searching migrations: %w", err)
	}

	rollbackFuncs := map[string]rollbackFunc{
		"001_create_users_table": func(db *gorm.DB, name string) error {
			return revertDropTable(db, &models.User{}, name)
		},
		"002_create_companies_table": func(db *gorm.DB, name string) error {
			return revertDropTable(db, &models.Company{}, name)
		},
		"003_create_farms_table": func(db *gorm.DB, name string) error {
			return revertDropTable(db, &models.Farm{}, name)
		},
		"004_create_animals_table": func(db *gorm.DB, name string) error {
			return revertDropTable(db, &models.Animal{}, name)
		},
		"005_create_milk_collections_table": func(db *gorm.DB, name string) error {
			return revertDropTable(db, &models.MilkCollection{}, name)
		},
		"006_create_reproductions_table": func(db *gorm.DB, name string) error {
			return revertDropTable(db, &models.Reproduction{}, name)
		},
		"007_create_weights_table": func(db *gorm.DB, name string) error {
			return revertDropTable(db, &models.Weight{}, name)
		},
		"008_create_persons_table": func(db *gorm.DB, name string) error {
			return revertDropTable(db, &models.Person{}, name)
		},
		"009_update_users_table": func(db *gorm.DB, name string) error {
			return revertRecreateUsersTable(db, name)
		},
		"010_create_expenses_table": func(db *gorm.DB, name string) error {
			return revertDropTable(db, &models.Expense{}, name)
		},
		"011_update_users_with_person": func(db *gorm.DB, name string) error {
			return revertRecreateUsersTable(db, name)
		},
		"012_add_company_name": func(db *gorm.DB, name string) error {
			return revertDropColumn(db, &models.Company{}, "company_name", name)
		},
		"013_add_farm_logo": func(db *gorm.DB, name string) error {
			return revertDropColumn(db, &models.Farm{}, "logo", name)
		},
		"014_add_animal_photo": func(db *gorm.DB, name string) error {
			return revertDropColumn(db, &models.Animal{}, "photo", name)
		},
		"015_update_animals_table": func(db *gorm.DB, name string) error {
			return revertAutoMigrate(db, &models.Animal{}, name)
		},
		"013_create_refresh_tokens_table": func(db *gorm.DB, name string) error {
			return revertDropTable(db, &models.RefreshToken{}, name)
		},
		"016_update_reproductions_table": func(db *gorm.DB, name string) error {
			return revertAutoMigrate(db, &models.Reproduction{}, name)
		},
	}

	for _, migration := range migrations {
		if err := revertMigration(db, migration, rollbackFuncs); err != nil {
			return err
		}
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

func createDebtsTable(db *gorm.DB) error {
	log.Printf("Creating debts table...")

	if err := db.AutoMigrate(&models.Debt{}); err != nil {
		return fmt.Errorf("error creating debts table: %w", err)
	}

	log.Printf("Debts table created successfully")
	return nil
}

func createVaccinesTable(db *gorm.DB) error {
	log.Printf("Creating vaccines table...")

	if err := db.AutoMigrate(&models.Vaccine{}); err != nil {
		return fmt.Errorf("error creating vaccines table: %w", err)
	}

	log.Printf("Vaccines table created successfully")
	return nil
}

func createVaccineApplicationsTable(db *gorm.DB) error {
	log.Printf("Creating vaccine_applications table...")

	if err := db.AutoMigrate(&models.VaccineApplication{}); err != nil {
		return fmt.Errorf("error creating vaccine_applications table: %w", err)
	}

	log.Printf("Vaccine applications table created successfully")
	return nil
}

func addFarmLanguage(db *gorm.DB) error {
	log.Printf("Adding language column to farms table...")

	if err := db.AutoMigrate(&models.Farm{}); err != nil {
		return fmt.Errorf("error adding language column to farms table: %w", err)
	}

	if err := db.Model(&models.Farm{}).Where("language = '' OR language IS NULL").Update("language", "pt").Error; err != nil {
		log.Printf("Warning: Could not update existing farms with default language: %v", err)
	}

	log.Printf("Language column added to farms table successfully")
	return nil
}
