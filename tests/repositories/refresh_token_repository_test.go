package repositories

import (
	"testing"
	"time"

	"github.com/fazendapro/FazendaPro-api/internal/models"
	"github.com/fazendapro/FazendaPro-api/internal/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupRefreshTokenTestDB(t *testing.T) (*repository.Database, *gorm.DB) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	err = db.AutoMigrate(&models.Company{}, &models.Farm{}, &models.Person{}, &models.User{}, &models.RefreshToken{})
	require.NoError(t, err)

	database := &repository.Database{DB: db}
	return database, db
}

func createTestUserForRefreshToken(t *testing.T, db *gorm.DB) *models.User {
	company := &models.Company{
		CompanyName: "Test Company",
		Location:    "Test Location",
		FarmCNPJ:    "12345678901234",
	}
	require.NoError(t, db.Create(company).Error)

	farm := &models.Farm{
		CompanyID: company.ID,
		Logo:      "",
	}
	require.NoError(t, db.Create(farm).Error)

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

	return user
}

func TestRefreshTokenRepository_Create(t *testing.T) {
	database, db := setupRefreshTokenTestDB(t)
	repo := repository.NewRefreshTokenRepository(database)

	user := createTestUserForRefreshToken(t, db)
	expiresAt := time.Now().Add(24 * time.Hour)

	token, err := repo.Create(user.ID, expiresAt)
	assert.NoError(t, err)
	assert.NotNil(t, token)
	assert.NotEmpty(t, token.Token)
	assert.Equal(t, user.ID, token.UserID)
	assert.True(t, token.ExpiresAt.Equal(expiresAt))
}

func TestRefreshTokenRepository_FindByToken(t *testing.T) {
	database, db := setupRefreshTokenTestDB(t)
	repo := repository.NewRefreshTokenRepository(database)

	user := createTestUserForRefreshToken(t, db)
	expiresAt := time.Now().Add(24 * time.Hour)

	createdToken, err := repo.Create(user.ID, expiresAt)
	require.NoError(t, err)

	found, err := repo.FindByToken(createdToken.Token)
	assert.NoError(t, err)
	assert.NotNil(t, found)
	assert.Equal(t, createdToken.Token, found.Token)
	assert.Equal(t, user.ID, found.UserID)
	assert.NotNil(t, found.User)
	assert.NotNil(t, found.User.Person)
}

func TestRefreshTokenRepository_FindByToken_NotFound(t *testing.T) {
	database, _ := setupRefreshTokenTestDB(t)
	repo := repository.NewRefreshTokenRepository(database)

	found, err := repo.FindByToken("non-existent-token")
	assert.NoError(t, err)
	assert.Nil(t, found)
}

func TestRefreshTokenRepository_FindByToken_Expired(t *testing.T) {
	database, db := setupRefreshTokenTestDB(t)
	repo := repository.NewRefreshTokenRepository(database)

	user := createTestUserForRefreshToken(t, db)
	expiresAt := time.Now().Add(-1 * time.Hour)

	createdToken, err := repo.Create(user.ID, expiresAt)
	require.NoError(t, err)

	require.NoError(t, db.Model(&models.RefreshToken{}).Where("id = ?", createdToken.ID).Update("expires_at", expiresAt).Error)

	found, err := repo.FindByToken(createdToken.Token)
	assert.NoError(t, err)
	assert.Nil(t, found)
}

func TestRefreshTokenRepository_DeleteByToken(t *testing.T) {
	database, db := setupRefreshTokenTestDB(t)
	repo := repository.NewRefreshTokenRepository(database)

	user := createTestUserForRefreshToken(t, db)
	expiresAt := time.Now().Add(24 * time.Hour)

	createdToken, err := repo.Create(user.ID, expiresAt)
	require.NoError(t, err)

	err = repo.DeleteByToken(createdToken.Token)
	assert.NoError(t, err)

	found, err := repo.FindByToken(createdToken.Token)
	assert.NoError(t, err)
	assert.Nil(t, found)
}

func TestRefreshTokenRepository_DeleteByUserID(t *testing.T) {
	database, db := setupRefreshTokenTestDB(t)
	repo := repository.NewRefreshTokenRepository(database)

	user1 := createTestUserForRefreshToken(t, db)

	company2 := &models.Company{
		CompanyName: "Test Company 2",
		Location:    "Test Location 2",
		FarmCNPJ:    "98765432109876",
	}
	require.NoError(t, db.Create(company2).Error)

	farm2 := &models.Farm{
		CompanyID: company2.ID,
		Logo:      "",
	}
	require.NoError(t, db.Create(farm2).Error)

	person2 := &models.Person{
		FirstName: "Maria",
		LastName:  "Santos",
		Email:     "maria@test.com",
		Password:  "password123",
		CPF:       "98765432100",
	}
	require.NoError(t, db.Create(person2).Error)

	user2 := &models.User{
		PersonID: &person2.ID,
		FarmID:   farm2.ID,
	}
	require.NoError(t, db.Create(user2).Error)

	expiresAt := time.Now().Add(24 * time.Hour)

	token1, err := repo.Create(user1.ID, expiresAt)
	require.NoError(t, err)
	token2, err := repo.Create(user1.ID, expiresAt)
	require.NoError(t, err)

	token3, err := repo.Create(user2.ID, expiresAt)
	require.NoError(t, err)

	err = repo.DeleteByUserID(user1.ID)
	assert.NoError(t, err)

	found1, err := repo.FindByToken(token1.Token)
	assert.NoError(t, err)
	assert.Nil(t, found1)

	found2, err := repo.FindByToken(token2.Token)
	assert.NoError(t, err)
	assert.Nil(t, found2)

	found3, err := repo.FindByToken(token3.Token)
	assert.NoError(t, err)
	assert.NotNil(t, found3)
}

func TestRefreshTokenRepository_DeleteExpired(t *testing.T) {
	database, db := setupRefreshTokenTestDB(t)
	repo := repository.NewRefreshTokenRepository(database)

	user := createTestUserForRefreshToken(t, db)

	expiredAt := time.Now().Add(-1 * time.Hour)
	expiredToken, err := repo.Create(user.ID, expiredAt)
	require.NoError(t, err)
	require.NoError(t, db.Model(&models.RefreshToken{}).Where("id = ?", expiredToken.ID).Update("expires_at", expiredAt).Error)

	validAt := time.Now().Add(24 * time.Hour)
	validToken, err := repo.Create(user.ID, validAt)
	require.NoError(t, err)

	err = repo.DeleteExpired()
	assert.NoError(t, err)

	var expiredCount int64
	db.Model(&models.RefreshToken{}).Where("id = ?", expiredToken.ID).Count(&expiredCount)
	assert.Equal(t, int64(0), expiredCount)

	found, err := repo.FindByToken(validToken.Token)
	assert.NoError(t, err)
	assert.NotNil(t, found)
}

func TestRefreshTokenRepository_DeleteExpired_Multiple(t *testing.T) {
	database, db := setupRefreshTokenTestDB(t)
	repo := repository.NewRefreshTokenRepository(database)

	user := createTestUserForRefreshToken(t, db)

	for i := 0; i < 3; i++ {
		expiredAt := time.Now().Add(-time.Duration(i+1) * time.Hour)
		token, err := repo.Create(user.ID, expiredAt)
		require.NoError(t, err)
		require.NoError(t, db.Model(&models.RefreshToken{}).Where("id = ?", token.ID).Update("expires_at", expiredAt).Error)
	}

	validAt := time.Now().Add(24 * time.Hour)
	validToken, err := repo.Create(user.ID, validAt)
	require.NoError(t, err)

	err = repo.DeleteExpired()
	assert.NoError(t, err)

	var count int64
	db.Model(&models.RefreshToken{}).Count(&count)
	assert.Equal(t, int64(1), count)

	found, err := repo.FindByToken(validToken.Token)
	assert.NoError(t, err)
	assert.NotNil(t, found)
}
