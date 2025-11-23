package models

import (
	"testing"
	"time"

	"github.com/fazendapro/FazendaPro-api/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestModels(t *testing.T) {
	t.Run("User_Model_Validation", func(t *testing.T) {
		user := &models.User{
			ID:     1,
			FarmID: 1,
		}

		assert.Equal(t, uint(1), user.ID)
		assert.Equal(t, uint(1), user.FarmID)
	})

	t.Run("Person_Model_Validation", func(t *testing.T) {
		person := &models.Person{
			ID:        1,
			FirstName: "João",
			LastName:  "Silva",
			Email:     "joao@fazenda.com",
			Password:  "senha123",
			CPF:       "12345678901",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		assert.Equal(t, uint(1), person.ID)
		assert.Equal(t, "João", person.FirstName)
		assert.Equal(t, "Silva", person.LastName)
		assert.Equal(t, "joao@fazenda.com", person.Email)
		assert.Equal(t, "senha123", person.Password)
		assert.Equal(t, "12345678901", person.CPF)
		assert.NotNil(t, person.CreatedAt)
		assert.NotNil(t, person.UpdatedAt)
	})

	t.Run("Animal_Model_Validation", func(t *testing.T) {
		birthDate := time.Now().AddDate(-2, 0, 0)
		animal := &models.Animal{
			ID:                1,
			FarmID:            1,
			EarTagNumberLocal: 123,
			AnimalName:        "Boi João",
			Sex:               1,
			Breed:             "Nelore",
			Type:              "Bovino",
			BirthDate:         &birthDate,
			AnimalType:        0,
			Status:            0,
			Purpose:           0,
			CreatedAt:         time.Now(),
			UpdatedAt:         time.Now(),
		}

		assert.Equal(t, uint(1), animal.ID)
		assert.Equal(t, uint(1), animal.FarmID)
		assert.Equal(t, 123, animal.EarTagNumberLocal)
		assert.Equal(t, "Boi João", animal.AnimalName)
		assert.Equal(t, 1, animal.Sex)
		assert.Equal(t, "Nelore", animal.Breed)
		assert.Equal(t, "Bovino", animal.Type)
		assert.Equal(t, &birthDate, animal.BirthDate)
		assert.Equal(t, 0, animal.AnimalType)
		assert.Equal(t, 0, animal.Status)
		assert.Equal(t, 0, animal.Purpose)
		assert.NotNil(t, animal.CreatedAt)
		assert.NotNil(t, animal.UpdatedAt)
	})

	t.Run("MilkCollection_Model_Validation", func(t *testing.T) {
		collectionDate := time.Now()
		milkCollection := &models.MilkCollection{
			ID:        1,
			AnimalID:  1,
			Liters:    25.5,
			Date:      collectionDate,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		assert.Equal(t, uint(1), milkCollection.ID)
		assert.Equal(t, uint(1), milkCollection.AnimalID)
		assert.Equal(t, 25.5, milkCollection.Liters)
		assert.Equal(t, collectionDate, milkCollection.Date)
		assert.NotNil(t, milkCollection.CreatedAt)
		assert.NotNil(t, milkCollection.UpdatedAt)
	})

	t.Run("Reproduction_Model_Validation", func(t *testing.T) {
		reproductionDate := time.Now()
		reproduction := &models.Reproduction{
			ID:                     1,
			AnimalID:               1,
			CurrentPhase:           models.PhaseVazias,
			InseminationDate:       &reproductionDate,
			InseminationType:       "Monta Natural",
			VeterinaryConfirmation: true,
			Observations:           "Reprodução bem-sucedida",
			CreatedAt:              time.Now(),
			UpdatedAt:              time.Now(),
		}

		assert.Equal(t, uint(1), reproduction.ID)
		assert.Equal(t, uint(1), reproduction.AnimalID)
		assert.Equal(t, models.PhaseVazias, reproduction.CurrentPhase)
		assert.Equal(t, &reproductionDate, reproduction.InseminationDate)
		assert.Equal(t, "Monta Natural", reproduction.InseminationType)
		assert.True(t, reproduction.VeterinaryConfirmation)
		assert.Equal(t, "Reprodução bem-sucedida", reproduction.Observations)
		assert.NotNil(t, reproduction.CreatedAt)
		assert.NotNil(t, reproduction.UpdatedAt)
	})

	t.Run("RefreshToken_Model_Validation", func(t *testing.T) {
		expiresAt := time.Now().Add(24 * time.Hour)
		refreshToken := &models.RefreshToken{
			ID:        1,
			UserID:    1,
			Token:     "refresh_token_123",
			ExpiresAt: expiresAt,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		assert.Equal(t, uint(1), refreshToken.ID)
		assert.Equal(t, uint(1), refreshToken.UserID)
		assert.Equal(t, "refresh_token_123", refreshToken.Token)
		assert.Equal(t, expiresAt, refreshToken.ExpiresAt)
		assert.NotNil(t, refreshToken.CreatedAt)
		assert.NotNil(t, refreshToken.UpdatedAt)
	})

	t.Run("Model_Relationships", func(t *testing.T) {
		user := &models.User{
			ID:     1,
			FarmID: 1,
		}

		person := &models.Person{
			ID:        1,
			FirstName: "João",
			LastName:  "Silva",
			Email:     "joao@fazenda.com",
			Password:  "senha123",
			CPF:       "12345678901",
		}

		animal := &models.Animal{
			ID:                1,
			FarmID:            1,
			EarTagNumberLocal: 123,
			AnimalName:        "Boi João",
			Breed:             "Nelore",
		}

		user.Person = person

		assert.NotNil(t, user.Person)
		assert.Equal(t, "João", user.Person.FirstName)
		assert.Equal(t, "Boi João", animal.AnimalName)
	})

	t.Run("Model_Edge_Cases", func(t *testing.T) {
		t.Run("Empty_Fields", func(t *testing.T) {
			user := &models.User{}
			assert.Equal(t, uint(0), user.ID)
			assert.Equal(t, uint(0), user.FarmID)
		})

		t.Run("Zero_Values", func(t *testing.T) {
			animal := &models.Animal{
				EarTagNumberLocal: 0,
			}
			assert.Equal(t, 0, animal.EarTagNumberLocal)
		})

		t.Run("Negative_Values", func(t *testing.T) {
			animal := &models.Animal{
				EarTagNumberLocal: -1,
			}
			assert.Equal(t, -1, animal.EarTagNumberLocal)
		})
	})

	t.Run("Sale_TableName", func(t *testing.T) {
		sale := models.Sale{}
		assert.Equal(t, "sales", sale.TableName())
	})

	t.Run("Company_Model_Validation", func(t *testing.T) {
		company := &models.Company{
			ID:          1,
			CompanyName: "Fazenda Teste",
			Location:    "São Paulo",
			FarmCNPJ:    "12345678901234",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}

		assert.Equal(t, uint(1), company.ID)
		assert.Equal(t, "Fazenda Teste", company.CompanyName)
		assert.Equal(t, "São Paulo", company.Location)
		assert.Equal(t, "12345678901234", company.FarmCNPJ)
		assert.NotNil(t, company.CreatedAt)
		assert.NotNil(t, company.UpdatedAt)
	})

	t.Run("Farm_Model_Validation", func(t *testing.T) {
		farm := &models.Farm{
			ID:        1,
			CompanyID: 1,
			Logo:      "logo.png",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		assert.Equal(t, uint(1), farm.ID)
		assert.Equal(t, uint(1), farm.CompanyID)
		assert.Equal(t, "logo.png", farm.Logo)
		assert.NotNil(t, farm.CreatedAt)
		assert.NotNil(t, farm.UpdatedAt)
	})

	t.Run("Debt_Model_Validation", func(t *testing.T) {
		debt := &models.Debt{
			ID:        1,
			Person:    "João Silva",
			Value:     1500.50,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		assert.Equal(t, uint(1), debt.ID)
		assert.Equal(t, "João Silva", debt.Person)
		assert.Equal(t, 1500.50, debt.Value)
		assert.NotNil(t, debt.CreatedAt)
		assert.NotNil(t, debt.UpdatedAt)
	})

	t.Run("Expense_Model_Validation", func(t *testing.T) {
		expenseDate := time.Now()
		expense := &models.Expense{
			ID:          1,
			FarmID:      1,
			Description: "Compra de ração",
			Amount:      500.00,
			Category:    "Alimentação",
			Date:        expenseDate,
			Notes:       "Ração para o mês",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}

		assert.Equal(t, uint(1), expense.ID)
		assert.Equal(t, uint(1), expense.FarmID)
		assert.Equal(t, "Compra de ração", expense.Description)
		assert.Equal(t, 500.00, expense.Amount)
		assert.Equal(t, "Alimentação", expense.Category)
		assert.Equal(t, expenseDate, expense.Date)
		assert.Equal(t, "Ração para o mês", expense.Notes)
		assert.NotNil(t, expense.CreatedAt)
		assert.NotNil(t, expense.UpdatedAt)
	})

	t.Run("Weight_Model_Validation", func(t *testing.T) {
		weightDate := time.Now()
		weight := &models.Weight{
			ID:           1,
			AnimalID:     1,
			Date:         weightDate,
			AnimalWeight: 450.5,
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		}

		assert.Equal(t, uint(1), weight.ID)
		assert.Equal(t, uint(1), weight.AnimalID)
		assert.Equal(t, weightDate, weight.Date)
		assert.Equal(t, 450.5, weight.AnimalWeight)
		assert.NotNil(t, weight.CreatedAt)
		assert.NotNil(t, weight.UpdatedAt)
	})

	t.Run("UserFarm_Model_Validation", func(t *testing.T) {
		userFarm := &models.UserFarm{
			ID:        1,
			UserID:    1,
			FarmID:    1,
			IsPrimary: true,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		assert.Equal(t, uint(1), userFarm.ID)
		assert.Equal(t, uint(1), userFarm.UserID)
		assert.Equal(t, uint(1), userFarm.FarmID)
		assert.True(t, userFarm.IsPrimary)
		assert.NotNil(t, userFarm.CreatedAt)
		assert.NotNil(t, userFarm.UpdatedAt)
	})

	t.Run("Batch_GetBatchByLiters", func(t *testing.T) {
		assert.Equal(t, models.Batch1, models.GetBatchByLiters(35.0))
		assert.Equal(t, models.Batch2, models.GetBatchByLiters(25.0))
		assert.Equal(t, models.Batch3, models.GetBatchByLiters(15.0))
		assert.Equal(t, models.Batch1, models.GetBatchByLiters(30.1))
		assert.Equal(t, models.Batch2, models.GetBatchByLiters(20.0))
		assert.Equal(t, models.Batch2, models.GetBatchByLiters(30.0))
	})

	t.Run("AnimalStatus_GetStatusName", func(t *testing.T) {
		assert.Equal(t, "Ativo", models.GetStatusName(models.AnimalStatusActive))
		assert.Equal(t, "Vendido", models.GetStatusName(models.AnimalStatusSold))
		assert.Equal(t, "Falecido", models.GetStatusName(models.AnimalStatusDeceased))
		assert.Equal(t, "Desconhecido", models.GetStatusName(999))
	})

	t.Run("ReproductionPhase_String", func(t *testing.T) {
		assert.Equal(t, "Lactação", models.PhaseLactacao.String())
		assert.Equal(t, "Secando", models.PhaseSecando.String())
		assert.Equal(t, "Vazias", models.PhaseVazias.String())
		assert.Equal(t, "Prenhas", models.PhasePrenhas.String())
		assert.Equal(t, "Desconhecida", models.ReproductionPhase(999).String())
	})
}
