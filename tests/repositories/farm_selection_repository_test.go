package repositories

import (
	"testing"

	"github.com/fazendapro/FazendaPro-api/internal/models"
	"github.com/fazendapro/FazendaPro-api/internal/repository"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupFarmSelectionTestDB() *gorm.DB {
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	db.AutoMigrate(&models.User{}, &models.Person{}, &models.Farm{}, &models.Company{})
	return db
}

func TestUserRepository_GetUserFarms(t *testing.T) {
	db := setupFarmSelectionTestDB()
	userRepo := repository.NewUserRepository(&repository.Database{DB: db})

	company := &models.Company{CompanyName: "Test Company"}
	db.Create(company)

	farm1 := &models.Farm{CompanyID: company.ID, Logo: "logo1.png"}
	farm2 := &models.Farm{CompanyID: company.ID, Logo: "logo2.png"}
	db.Create(farm1)
	db.Create(farm2)

	person := &models.Person{
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john@test.com",
		Password:  "password123",
		CPF:       "12345678901",
	}
	db.Create(person)

	user := &models.User{
		PersonID: &person.ID,
		FarmID:   farm1.ID,
	}
	db.Create(user)

	farms, err := userRepo.GetUserFarms(user.ID)

	assert.NoError(t, err)
	assert.Len(t, farms, 1)
	assert.Equal(t, farm1.ID, farms[0].ID)
}

func TestUserRepository_GetUserFarmCount(t *testing.T) {
	db := setupFarmSelectionTestDB()
	userRepo := repository.NewUserRepository(&repository.Database{DB: db})

	company := &models.Company{CompanyName: "Test Company"}
	db.Create(company)

	farm := &models.Farm{CompanyID: company.ID, Logo: "logo1.png"}
	db.Create(farm)

	person := &models.Person{
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john@test.com",
		Password:  "password123",
		CPF:       "12345678901",
	}
	db.Create(person)

	user := &models.User{
		PersonID: &person.ID,
		FarmID:   farm.ID,
	}
	db.Create(user)

	count, err := userRepo.GetUserFarmCount(user.ID)

	assert.NoError(t, err)
	assert.Equal(t, int64(1), count)
}

func TestUserRepository_GetUserFarmByID(t *testing.T) {
	db := setupFarmSelectionTestDB()
	userRepo := repository.NewUserRepository(&repository.Database{DB: db})

	company := &models.Company{CompanyName: "Test Company"}
	db.Create(company)

	farm := &models.Farm{CompanyID: company.ID, Logo: "logo1.png"}
	db.Create(farm)

	person := &models.Person{
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john@test.com",
		Password:  "password123",
		CPF:       "12345678901",
	}
	db.Create(person)

	user := &models.User{
		PersonID: &person.ID,
		FarmID:   farm.ID,
	}
	db.Create(user)

	foundFarm, err := userRepo.GetUserFarmByID(user.ID, farm.ID)

	assert.NoError(t, err)
	assert.NotNil(t, foundFarm)
	assert.Equal(t, farm.ID, foundFarm.ID)
}

func TestUserRepository_GetUserFarmByID_NotFound(t *testing.T) {
	db := setupFarmSelectionTestDB()
	userRepo := repository.NewUserRepository(&repository.Database{DB: db})

	company := &models.Company{CompanyName: "Test Company"}
	db.Create(company)

	farm := &models.Farm{CompanyID: company.ID, Logo: "logo1.png"}
	db.Create(farm)

	person := &models.Person{
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john@test.com",
		Password:  "password123",
		CPF:       "12345678901",
	}
	db.Create(person)

	user := &models.User{
		PersonID: &person.ID,
		FarmID:   farm.ID,
	}
	db.Create(user)

	foundFarm, err := userRepo.GetUserFarmByID(user.ID, 999)

	assert.Error(t, err)
	assert.Nil(t, foundFarm)
	assert.Equal(t, gorm.ErrRecordNotFound, err)
}
