package repositories

import (
	"testing"

	"github.com/fazendapro/FazendaPro-api/internal/models"
	"github.com/fazendapro/FazendaPro-api/internal/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupUserTestDB(t *testing.T) (*repository.Database, *gorm.DB) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	err = db.AutoMigrate(&models.Company{}, &models.Farm{}, &models.Person{}, &models.User{}, &models.UserFarm{})
	require.NoError(t, err)

	database := &repository.Database{DB: db}
	return database, db
}

func createTestCompanyForUser(t *testing.T, db *gorm.DB) *models.Company {
	company := &models.Company{
		CompanyName: "Test Company",
		Location:    "Test Location",
		FarmCNPJ:    "12345678901234",
	}
	require.NoError(t, db.Create(company).Error)
	return company
}

func createTestFarmForUser(t *testing.T, db *gorm.DB, companyID uint) *models.Farm {
	farm := &models.Farm{
		CompanyID: companyID,
		Logo:      "",
	}
	require.NoError(t, db.Create(farm).Error)
	return farm
}

func TestUserRepository_FindByPersonEmail(t *testing.T) {
	database, db := setupUserTestDB(t)
	repo := repository.NewUserRepository(database)

	company := createTestCompanyForUser(t, db)
	farm := createTestFarmForUser(t, db, company.ID)

	person := &models.Person{
		FirstName: "João",
		LastName:  "Silva",
		Email:     "joao@test.com",
		Password:  "password123",
		CPF:       "12345678901",
	}
	require.NoError(t, db.Create(person).Error)

	user := &models.User{
		PersonID: &person.ID,
		FarmID:   farm.ID,
	}
	require.NoError(t, db.Create(user).Error)

	found, err := repo.FindByPersonEmail("joao@test.com")
	assert.NoError(t, err)
	assert.NotNil(t, found)
	assert.Equal(t, user.ID, found.ID)
	assert.NotNil(t, found.Person)
	assert.Equal(t, "João", found.Person.FirstName)
}

func TestUserRepository_FindByPersonEmail_NotFound(t *testing.T) {
	database, _ := setupUserTestDB(t)
	repo := repository.NewUserRepository(database)

	found, err := repo.FindByPersonEmail("inexistente@test.com")
	assert.NoError(t, err)
	assert.Nil(t, found)
}

func TestUserRepository_CreateWithPerson(t *testing.T) {
	database, db := setupUserTestDB(t)
	repo := repository.NewUserRepository(database)

	company := createTestCompanyForUser(t, db)
	farm := createTestFarmForUser(t, db, company.ID)

	person := &models.Person{
		FirstName: "Maria",
		LastName:  "Santos",
		Email:     "maria@test.com",
		Password:  "password123",
		CPF:       "98765432100",
	}

	user := &models.User{
		FarmID: farm.ID,
	}

	err := repo.CreateWithPerson(user, person)
	assert.NoError(t, err)
	assert.NotZero(t, user.ID)
	assert.NotNil(t, user.PersonID)
	assert.Equal(t, person.ID, *user.PersonID)

	var foundUser models.User
	require.NoError(t, db.Preload("Person").First(&foundUser, user.ID).Error)
	assert.Equal(t, "Maria", foundUser.Person.FirstName)
}

func TestUserRepository_FindByIDWithPerson(t *testing.T) {
	database, db := setupUserTestDB(t)
	repo := repository.NewUserRepository(database)

	company := createTestCompanyForUser(t, db)
	farm := createTestFarmForUser(t, db, company.ID)

	person := &models.Person{
		FirstName: "João",
		LastName:  "Silva",
		Email:     "joao@test.com",
		Password:  "password123",
		CPF:       "12345678901",
	}
	require.NoError(t, db.Create(person).Error)

	user := &models.User{
		PersonID: &person.ID,
		FarmID:   farm.ID,
	}
	require.NoError(t, db.Create(user).Error)

	found, err := repo.FindByIDWithPerson(user.ID)
	assert.NoError(t, err)
	assert.NotNil(t, found)
	assert.Equal(t, user.ID, found.ID)
	assert.NotNil(t, found.Person)
	assert.Equal(t, "João", found.Person.FirstName)
	assert.NotNil(t, found.Farm)
}

func TestUserRepository_FindByIDWithPerson_NotFound(t *testing.T) {
	database, _ := setupUserTestDB(t)
	repo := repository.NewUserRepository(database)

	found, err := repo.FindByIDWithPerson(999)
	assert.NoError(t, err)
	assert.Nil(t, found)
}

func TestUserRepository_UpdatePersonData(t *testing.T) {
	database, db := setupUserTestDB(t)
	repo := repository.NewUserRepository(database)

	company := createTestCompanyForUser(t, db)
	farm := createTestFarmForUser(t, db, company.ID)

	person := &models.Person{
		FirstName: "João",
		LastName:  "Silva",
		Email:     "joao@test.com",
		Password:  "password123",
		CPF:       "12345678901",
	}
	require.NoError(t, db.Create(person).Error)

	user := &models.User{
		PersonID: &person.ID,
		FarmID:   farm.ID,
	}
	require.NoError(t, db.Create(user).Error)

	updatedPerson := &models.Person{
		FirstName: "João Atualizado",
		LastName:  "Silva Santos",
		Email:     "joao.novo@test.com",
	}

	err := repo.UpdatePersonData(user.ID, updatedPerson)
	assert.NoError(t, err)

	var foundPerson models.Person
	require.NoError(t, db.First(&foundPerson, person.ID).Error)
	assert.Equal(t, "João Atualizado", foundPerson.FirstName)
	assert.Equal(t, "joao.novo@test.com", foundPerson.Email)
}

func TestUserRepository_ValidatePassword(t *testing.T) {
	database, db := setupUserTestDB(t)
	repo := repository.NewUserRepository(database)

	company := createTestCompanyForUser(t, db)
	farm := createTestFarmForUser(t, db, company.ID)

	person := &models.Person{
		FirstName: "João",
		LastName:  "Silva",
		Email:     "joao@test.com",
		Password:  "senha123",
		CPF:       "12345678901",
	}
	require.NoError(t, db.Create(person).Error)

	user := &models.User{
		PersonID: &person.ID,
		FarmID:   farm.ID,
	}
	require.NoError(t, db.Create(user).Error)

	valid, err := repo.ValidatePassword(user.ID, "senha123")
	assert.NoError(t, err)
	assert.True(t, valid)

	invalid, err := repo.ValidatePassword(user.ID, "senha_errada")
	assert.NoError(t, err)
	assert.False(t, invalid)
}

func TestUserRepository_FindByEmail(t *testing.T) {
	database, db := setupUserTestDB(t)
	repo := repository.NewUserRepository(database)

	company := createTestCompanyForUser(t, db)
	farm := createTestFarmForUser(t, db, company.ID)

	person := &models.Person{
		FirstName: "João",
		LastName:  "Silva",
		Email:     "joao@test.com",
		Password:  "password123",
		CPF:       "12345678901",
	}
	require.NoError(t, db.Create(person).Error)

	user := &models.User{
		PersonID: &person.ID,
		FarmID:   farm.ID,
	}
	require.NoError(t, db.Create(user).Error)

	found, err := repo.FindByEmail("joao@test.com")
	assert.NoError(t, err)
	assert.NotNil(t, found)
	assert.Equal(t, user.ID, found.ID)
}

func TestUserRepository_FarmExists(t *testing.T) {
	database, db := setupUserTestDB(t)
	repo := repository.NewUserRepository(database)

	company := createTestCompanyForUser(t, db)
	farm := createTestFarmForUser(t, db, company.ID)

	exists, err := repo.FarmExists(farm.ID)
	assert.NoError(t, err)
	assert.True(t, exists)

	notExists, err := repo.FarmExists(999)
	assert.NoError(t, err)
	assert.False(t, notExists)
}

func TestUserRepository_CreateDefaultFarm(t *testing.T) {
	database, db := setupUserTestDB(t)
	repo := repository.NewUserRepository(database)

	farmID := uint(100)

	err := repo.CreateDefaultFarm(farmID)
	assert.NoError(t, err)

	var farm models.Farm
	require.NoError(t, db.Preload("Company").First(&farm, farmID).Error)
	assert.Equal(t, farmID, farm.ID)
	assert.Equal(t, "FazendaPro Demo", farm.Company.CompanyName)
}

func TestUserRepository_GetUserFarms_Real(t *testing.T) {
	database, db := setupUserTestDB(t)
	repo := repository.NewUserRepository(database)

	company1 := createTestCompanyForUser(t, db)
	farm1 := createTestFarmForUser(t, db, company1.ID)

	company2 := &models.Company{
		CompanyName: "Test Company 2",
		Location:    "Test Location 2",
		FarmCNPJ:    "98765432109876",
	}
	require.NoError(t, db.Create(company2).Error)
	farm2 := createTestFarmForUser(t, db, company2.ID)

	person := &models.Person{
		FirstName: "João",
		LastName:  "Silva",
		Email:     "joao@test.com",
		Password:  "password123",
		CPF:       "12345678901",
	}
	require.NoError(t, db.Create(person).Error)

	user := &models.User{
		PersonID: &person.ID,
		FarmID:   farm1.ID,
	}
	require.NoError(t, db.Create(user).Error)

	userFarm1 := &models.UserFarm{
		UserID: user.ID,
		FarmID: farm1.ID,
	}
	require.NoError(t, db.Create(userFarm1).Error)

	userFarm2 := &models.UserFarm{
		UserID: user.ID,
		FarmID: farm2.ID,
	}
	require.NoError(t, db.Create(userFarm2).Error)

	farms, err := repo.GetUserFarms(user.ID)
	assert.NoError(t, err)
	assert.Len(t, farms, 2)
}

func TestUserRepository_GetUserFarmCount_Real(t *testing.T) {
	database, db := setupUserTestDB(t)
	repo := repository.NewUserRepository(database)

	company := createTestCompanyForUser(t, db)
	farm1 := createTestFarmForUser(t, db, company.ID)
	farm2 := createTestFarmForUser(t, db, company.ID)

	person := &models.Person{
		FirstName: "João",
		LastName:  "Silva",
		Email:     "joao@test.com",
		Password:  "password123",
		CPF:       "12345678901",
	}
	require.NoError(t, db.Create(person).Error)

	user := &models.User{
		PersonID: &person.ID,
		FarmID:   farm1.ID,
	}
	require.NoError(t, db.Create(user).Error)

	userFarm1 := &models.UserFarm{
		UserID: user.ID,
		FarmID: farm1.ID,
	}
	require.NoError(t, db.Create(userFarm1).Error)

	userFarm2 := &models.UserFarm{
		UserID: user.ID,
		FarmID: farm2.ID,
	}
	require.NoError(t, db.Create(userFarm2).Error)

	count, err := repo.GetUserFarmCount(user.ID)
	assert.NoError(t, err)
	assert.Equal(t, int64(2), count)
}

func TestUserRepository_GetUserFarmByID_Real(t *testing.T) {
	database, db := setupUserTestDB(t)
	repo := repository.NewUserRepository(database)

	company := createTestCompanyForUser(t, db)
	farm := createTestFarmForUser(t, db, company.ID)

	person := &models.Person{
		FirstName: "João",
		LastName:  "Silva",
		Email:     "joao@test.com",
		Password:  "password123",
		CPF:       "12345678901",
	}
	require.NoError(t, db.Create(person).Error)

	user := &models.User{
		PersonID: &person.ID,
		FarmID:   farm.ID,
	}
	require.NoError(t, db.Create(user).Error)

	userFarm := &models.UserFarm{
		UserID: user.ID,
		FarmID: farm.ID,
	}
	require.NoError(t, db.Create(userFarm).Error)

	found, err := repo.GetUserFarmByID(user.ID, farm.ID)
	assert.NoError(t, err)
	assert.NotNil(t, found)
	assert.Equal(t, farm.ID, found.ID)
	assert.NotNil(t, found.Company)
}

func TestUserRepository_GetUserFarmByID_NotFound_Real(t *testing.T) {
	database, db := setupUserTestDB(t)
	repo := repository.NewUserRepository(database)

	company := createTestCompanyForUser(t, db)
	farm := createTestFarmForUser(t, db, company.ID)

	person := &models.Person{
		FirstName: "João",
		LastName:  "Silva",
		Email:     "joao@test.com",
		Password:  "password123",
		CPF:       "12345678901",
	}
	require.NoError(t, db.Create(person).Error)

	user := &models.User{
		PersonID: &person.ID,
		FarmID:   farm.ID,
	}
	require.NoError(t, db.Create(user).Error)

	found, err := repo.GetUserFarmByID(user.ID, 999)
	assert.Error(t, err)
	assert.Nil(t, found)
}

func TestUserRepository_CreateUserFarm(t *testing.T) {
	database, db := setupUserTestDB(t)
	repo := repository.NewUserRepository(database)

	company := createTestCompanyForUser(t, db)
	farm := createTestFarmForUser(t, db, company.ID)

	person := &models.Person{
		FirstName: "João",
		LastName:  "Silva",
		Email:     "joao@test.com",
		Password:  "password123",
		CPF:       "12345678901",
	}
	require.NoError(t, db.Create(person).Error)

	user := &models.User{
		PersonID: &person.ID,
		FarmID:   farm.ID,
	}
	require.NoError(t, db.Create(user).Error)

	userFarm := &models.UserFarm{
		UserID: user.ID,
		FarmID: farm.ID,
	}

	err := repo.CreateUserFarm(userFarm)
	assert.NoError(t, err)
	assert.NotZero(t, userFarm.ID)
}
