package service

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/fazendapro/FazendaPro-api/internal/cache"
	"github.com/fazendapro/FazendaPro-api/internal/models"
	"github.com/fazendapro/FazendaPro-api/internal/repository"
)

type AnimalService struct {
	repository repository.AnimalRepositoryInterface
	cache      cache.CacheInterface
}

func NewAnimalService(repository repository.AnimalRepositoryInterface, cacheClient cache.CacheInterface) *AnimalService {
	return &AnimalService{
		repository: repository,
		cache:      cacheClient,
	}
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

	err = s.repository.Create(animal)
	if err != nil {
		return err
	}

	cacheKey := fmt.Sprintf(CacheKeyAnimalsFarm, animal.FarmID)
	if err := s.cache.Delete(cacheKey); err != nil {
		log.Printf(ErrInvalidateCache, err)
	}

	return nil
}

func (s *AnimalService) GetAnimalByID(id uint) (*models.Animal, error) {
	return s.repository.FindByID(id)
}

func (s *AnimalService) GetAnimalsByFarmID(farmID uint) ([]models.Animal, error) {
	cacheKey := fmt.Sprintf(CacheKeyAnimalsFarm, farmID)
	var cachedAnimals []models.Animal

	err := s.cache.Get(cacheKey, &cachedAnimals)
	if err == nil {
		log.Printf("Cache HIT para animais da fazenda %d", farmID)
		return cachedAnimals, nil
	}

	log.Printf("Cache MISS para animais da fazenda %d", farmID)
	animals, err := s.repository.FindByFarmID(farmID)
	if err != nil {
		return nil, err
	}

	if err := s.cache.Set(cacheKey, animals, 300); err != nil {
		log.Printf("Erro ao salvar no cache (não crítico): %v", err)
	}

	return animals, nil
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

	err = s.repository.Update(animal)
	if err != nil {
		return err
	}

	cacheKey := fmt.Sprintf(CacheKeyAnimalsFarm, animal.FarmID)
	if err := s.cache.Delete(cacheKey); err != nil {
		log.Printf(ErrInvalidateCache, err)
	}

	return nil
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

	farmID := existingAnimal.FarmID
	err = s.repository.Delete(id)
	if err != nil {
		return err
	}

	cacheKey := fmt.Sprintf(CacheKeyAnimalsFarm, farmID)
	if err := s.cache.Delete(cacheKey); err != nil {
		log.Printf(ErrInvalidateCache, err)
	}

	return nil
}

func (s *AnimalService) GetAnimalsByFarmIDAndSex(farmID uint, sex int) ([]models.Animal, error) {
	return s.repository.FindByFarmIDAndSex(farmID, sex)
}

func (s *AnimalService) GetAnimalsByFarmIDWithPagination(farmID uint, page, limit int) ([]models.Animal, int64, error) {
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}

	return s.repository.FindByFarmIDWithPagination(farmID, page, limit)
}
