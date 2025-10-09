package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/fazendapro/FazendaPro-api/cmd/app"
	"github.com/fazendapro/FazendaPro-api/config"
	"github.com/fazendapro/FazendaPro-api/internal/migrations"
	"github.com/fazendapro/FazendaPro-api/internal/models"
	"github.com/fazendapro/FazendaPro-api/internal/repository"
	"github.com/fazendapro/FazendaPro-api/internal/routes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *repository.Database {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	err = migrations.RunMigrations(db)
	require.NoError(t, err)

	company := &models.Company{
		CompanyName: "Test Company",
	}
	err = db.Create(company).Error
	require.NoError(t, err)

	farm := &models.Farm{
		CompanyID: company.ID,
		Logo:      "",
	}
	err = db.Create(farm).Error
	require.NoError(t, err)

	return &repository.Database{DB: db}
}

func setupTestApp(t *testing.T) (*app.Application, *repository.Database, *config.Config) {
	app, err := app.NewApplication()
	require.NoError(t, err)

	db := setupTestDB(t)

	cfg := &config.Config{
		Port:      "8080",
		JWTSecret: "test-secret-key",
		CORS: config.CORSConfig{
			AllowedOrigins:   []string{"*"},
			AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
			AllowedHeaders:   []string{"Content-Type", "Authorization"},
			ExposedHeaders:   []string{},
			AllowCredentials: true,
			MaxAge:           86400,
		},
	}

	return app, db, cfg
}

func TestIntegration(t *testing.T) {
	app, db, cfg := setupTestApp(t)
	defer db.Close()

	router := routes.SetupRoutes(app, db, cfg)

	t.Run("Health_Check", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/health", nil)
		require.NoError(t, err)

		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Contains(t, rr.Body.String(), "OK")
	})

	t.Run("Root_Endpoint", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/", nil)
		require.NoError(t, err)

		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Contains(t, rr.Body.String(), "FazendaPro API is running!")
	})

	t.Run("CORS_Headers", func(t *testing.T) {
		req, err := http.NewRequest("OPTIONS", "/api/v1/auth/register", nil)
		require.NoError(t, err)
		req.Header.Set("Origin", "http://localhost:3000")

		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Equal(t, "http://localhost:3000", rr.Header().Get("Access-Control-Allow-Origin"))
	})
}

func TestAuthIntegration(t *testing.T) {
	app, db, cfg := setupTestApp(t)
	defer db.Close()

	router := routes.SetupRoutes(app, db, cfg)

	t.Run("Register_User_Success", func(t *testing.T) {
		userData := map[string]interface{}{
			"user": map[string]interface{}{
				"farm_id": 1,
			},
			"person": map[string]interface{}{
				"first_name": "João",
				"last_name":  "Silva",
				"email":      "joao@test.com",
				"password":   "123456",
				"cpf":        "12345678901",
			},
		}

		jsonData, err := json.Marshal(userData)
		require.NoError(t, err)

		req, err := http.NewRequest("POST", "/api/v1/auth/register", bytes.NewBuffer(jsonData))
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusCreated, rr.Code)

		var response map[string]interface{}
		err = json.Unmarshal(rr.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.True(t, response["success"].(bool))
		assert.Contains(t, response, "access_token")
		assert.Contains(t, response, "refresh_token")
		assert.Contains(t, response, "user")
	})

	t.Run("Register_User_Invalid_Data", func(t *testing.T) {
		userData := map[string]interface{}{
			"user": map[string]interface{}{
				"farm_id": 1,
			},
			"person": map[string]interface{}{
				"first_name": "",
				"last_name":  "Silva",
				"email":      "invalid-email",
				"password":   "123",
				"cpf":        "",
			},
		}

		jsonData, err := json.Marshal(userData)
		require.NoError(t, err)

		req, err := http.NewRequest("POST", "/api/v1/auth/register", bytes.NewBuffer(jsonData))
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)

		var response map[string]interface{}
		err = json.Unmarshal(rr.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.False(t, response["success"].(bool))
		assert.Contains(t, response, "error")
	})

	t.Run("Login_User_Success", func(t *testing.T) {
		userData := map[string]interface{}{
			"user": map[string]interface{}{
				"farm_id": 1,
			},
			"person": map[string]interface{}{
				"first_name": "Maria",
				"last_name":  "Santos",
				"email":      "maria@test.com",
				"password":   "123456",
				"cpf":        "98765432100",
			},
		}

		jsonData, err := json.Marshal(userData)
		require.NoError(t, err)

		req, err := http.NewRequest("POST", "/api/v1/auth/register", bytes.NewBuffer(jsonData))
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusCreated, rr.Code)

		loginData := map[string]interface{}{
			"email":    "maria@test.com",
			"password": "123456",
		}

		loginJson, err := json.Marshal(loginData)
		require.NoError(t, err)

		req, err = http.NewRequest("POST", "/api/v1/auth/login", bytes.NewBuffer(loginJson))
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		rr = httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)

		var response map[string]interface{}
		err = json.Unmarshal(rr.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.True(t, response["success"].(bool))
		assert.Contains(t, response, "access_token")
		assert.Contains(t, response, "refresh_token")
	})

	t.Run("Login_User_Invalid_Credentials", func(t *testing.T) {
		loginData := map[string]interface{}{
			"email":    "nonexistent@test.com",
			"password": "wrongpassword",
		}

		jsonData, err := json.Marshal(loginData)
		require.NoError(t, err)

		req, err := http.NewRequest("POST", "/api/v1/auth/login", bytes.NewBuffer(jsonData))
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusUnauthorized, rr.Code)

		var response map[string]interface{}
		err = json.Unmarshal(rr.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.False(t, response["success"].(bool))
		assert.Contains(t, response["message"], "Credenciais inválidas")
	})
}

func TestProtectedEndpoints(t *testing.T) {
	app, db, cfg := setupTestApp(t)
	defer db.Close()

	router := routes.SetupRoutes(app, db, cfg)

	userData := map[string]interface{}{
		"user": map[string]interface{}{
			"farm_id": 1,
		},
		"person": map[string]interface{}{
			"first_name": "Test",
			"last_name":  "User",
			"email":      "test@example.com",
			"password":   "123456",
			"cpf":        "11122233344",
		},
	}

	jsonData, err := json.Marshal(userData)
	require.NoError(t, err)

	req, err := http.NewRequest("POST", "/api/v1/auth/register", bytes.NewBuffer(jsonData))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusCreated, rr.Code)

	var authResponse map[string]interface{}
	err = json.Unmarshal(rr.Body.Bytes(), &authResponse)
	require.NoError(t, err)

	accessToken := authResponse["access_token"].(string)

	t.Run("Get_User_With_Valid_Token", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/api/v1/users", nil)
		require.NoError(t, err)
		req.Header.Set("Authorization", "Bearer "+accessToken)

		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
	})

	t.Run("Get_User_Without_Token", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/api/v1/users", nil)
		require.NoError(t, err)

		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusUnauthorized, rr.Code)

		var response map[string]interface{}
		err = json.Unmarshal(rr.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.Contains(t, response["message"], "Token não fornecido")
	})

	t.Run("Get_User_With_Invalid_Token", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/api/v1/users", nil)
		require.NoError(t, err)
		req.Header.Set("Authorization", "Bearer invalid-token")

		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusUnauthorized, rr.Code)
	})
}

func TestDatabaseIntegration(t *testing.T) {
	app, db, cfg := setupTestApp(t)
	defer db.Close()

	t.Run("Database_Connection", func(t *testing.T) {
		assert.NotNil(t, db)
		assert.NotNil(t, db.DB)
	})

	t.Run("Database_Migrations", func(t *testing.T) {
		var tables []string
		err := db.DB.Raw("SELECT name FROM sqlite_master WHERE type='table'").Scan(&tables).Error
		require.NoError(t, err)

		expectedTables := []string{"users", "people", "companies", "farms", "animals", "milk_collections", "reproductions", "weights", "expenses", "refresh_tokens"}

		for _, expectedTable := range expectedTables {
			assert.Contains(t, tables, expectedTable, "Table %s should exist", expectedTable)
		}
	})

	t.Run("Database_Seed_Data", func(t *testing.T) {
		var companyCount int64
		err := db.DB.Model(&models.Company{}).Count(&companyCount).Error
		require.NoError(t, err)
		assert.Greater(t, companyCount, int64(0))

		var farmCount int64
		err = db.DB.Model(&models.Farm{}).Count(&farmCount).Error
		require.NoError(t, err)
		assert.Greater(t, farmCount, int64(0))
	})
}

func TestErrorHandling(t *testing.T) {
	app, db, cfg := setupTestApp(t)
	defer db.Close()

	router := routes.SetupRoutes(app, db, cfg)

	t.Run("Invalid_JSON", func(t *testing.T) {
		req, err := http.NewRequest("POST", "/api/v1/auth/register", bytes.NewBufferString("invalid json"))
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
	})

	t.Run("Method_Not_Allowed", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/api/v1/auth/register", nil)
		require.NoError(t, err)

		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusMethodNotAllowed, rr.Code)
	})

	t.Run("Non_Existent_Endpoint", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/api/v1/nonexistent", nil)
		require.NoError(t, err)

		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusNotFound, rr.Code)
	})
}

func TestPerformance(t *testing.T) {
	app, db, cfg := setupTestApp(t)
	defer db.Close()

	router := routes.SetupRoutes(app, db, cfg)

	t.Run("Health_Check_Performance", func(t *testing.T) {
		start := time.Now()

		req, err := http.NewRequest("GET", "/health", nil)
		require.NoError(t, err)

		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		duration := time.Since(start)

		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Less(t, duration, 100*time.Millisecond, "Health check should be fast")
	})

	t.Run("Concurrent_Requests", func(t *testing.T) {
		concurrency := 10
		done := make(chan bool, concurrency)

		for i := 0; i < concurrency; i++ {
			go func() {
				req, err := http.NewRequest("GET", "/health", nil)
				require.NoError(t, err)

				rr := httptest.NewRecorder()
				router.ServeHTTP(rr, req)

				assert.Equal(t, http.StatusOK, rr.Code)
				done <- true
			}()
		}

		for i := 0; i < concurrency; i++ {
			<-done
		}
	})
}
