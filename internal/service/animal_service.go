package service

import (
	"errors"
	"log"
	"time"

	"github.com/fazendapro/FazendaPro-api/internal/models"
	"github.com/fazendapro/FazendaPro-api/internal/repository"
)

type AnimalService struct {
	repository repository.AnimalRepositoryInterface
}

func NewAnimalService(repository repository.AnimalRepositoryInterface) *AnimalService {
	return &AnimalService{repository: repository}
}

func (s *AnimalService) CreateAnimal(animal *models.Animal) error {
	log.Println("Creating animal", animal)
	if animal.FarmID == 0 {
		return errors.New("farm ID é obrigatório")
	}

	if animal.EarTagNumberLocal == 0 {
		return errors.New("número da brinca local é obrigatório")
	}

	if animal.AnimalName == "" {
		return errors.New("nome do animal é obrigatório")
	}

	if animal.Breed == "" {
		return errors.New("raça do animal é obrigatória")
	}

	if animal.Type == "" {
		return errors.New("tipo do animal é obrigatório")
	}

	if animal.Sex != 0 && animal.Sex != 1 {
		return errors.New("sexo deve ser 0 (Fêmea) ou 1 (Macho)")
	}

	if animal.AnimalType < 0 || animal.AnimalType > 10 {
		return errors.New("tipo de animal inválido")
	}

	if animal.Purpose < 0 || animal.Purpose > 2 {
		return errors.New("propósito deve ser 0 (Carne), 1 (Leite) ou 2 (Reprodução)")
	}

	existingAnimal, err := s.repository.FindByEarTagNumber(animal.FarmID, animal.EarTagNumberLocal)
	if err != nil {
		return err
	}

	if existingAnimal != nil {
		return errors.New("já existe um animal com este número de brinca nesta fazenda")
	}

	if animal.Status == 0 {
		animal.Status = 0
	}

	now := time.Now()
	animal.CreatedAt = now
	animal.UpdatedAt = now

	return s.repository.Create(animal)
}

func (s *AnimalService) GetAnimalByID(id uint) (*models.Animal, error) {
	return s.repository.FindByID(id)
}

func (s *AnimalService) GetAnimalsByFarmID(farmID uint) ([]models.Animal, error) {
	return s.repository.FindByFarmID(farmID)
}

func (s *AnimalService) UpdateAnimal(animal *models.Animal) error {
	if animal.ID == 0 {
		return errors.New("ID do animal é obrigatório")
	}

	existingAnimal, err := s.repository.FindByID(animal.ID)
	if err != nil {
		return err
	}

	if existingAnimal == nil {
		return errors.New("animal não encontrado")
	}

	if animal.Sex != 0 && animal.Sex != 1 {
		return errors.New("sexo deve ser 0 (Fêmea) ou 1 (Macho)")
	}

	if animal.AnimalType < 0 || animal.AnimalType > 10 {
		return errors.New("tipo de animal inválido")
	}

	if animal.Purpose < 0 || animal.Purpose > 2 {
		return errors.New("propósito deve ser 0 (Carne), 1 (Leite) ou 2 (Reprodução)")
	}

	animal.FarmID = existingAnimal.FarmID

	now := time.Now()
	animal.UpdatedAt = now

	return s.repository.Update(animal)
}

func (s *AnimalService) DeleteAnimal(id uint) error {
	if id == 0 {
		return errors.New("ID do animal é obrigatório")
	}

	existingAnimal, err := s.repository.FindByID(id)
	if err != nil {
		return err
	}

	if existingAnimal == nil {
		return errors.New("animal não encontrado")
	}

	return s.repository.Delete(id)
}

func (s *AnimalService) GetAnimalsByFarmIDAndSex(farmID uint, sex int) ([]models.Animal, error) {
	return s.repository.FindByFarmIDAndSex(farmID, sex)
}
