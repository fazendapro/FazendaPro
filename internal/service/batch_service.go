package service

import (
	"github.com/fazendapro/FazendaPro-api/internal/models"
	"github.com/fazendapro/FazendaPro-api/internal/repository"
)

type BatchService struct {
	animalRepository repository.AnimalRepositoryInterface
	milkRepository   repository.MilkCollectionRepositoryInterface
}

func NewBatchService(animalRepository repository.AnimalRepositoryInterface, milkRepository repository.MilkCollectionRepositoryInterface) *BatchService {
	return &BatchService{
		animalRepository: animalRepository,
		milkRepository:   milkRepository,
	}
}

func (s *BatchService) UpdateAnimalBatch(animalID uint) error {
	animal, err := s.animalRepository.FindByID(animalID)
	if err != nil {
		return err
	}

	milkCollections, err := s.milkRepository.FindByAnimalID(animalID)
	if err != nil {
		return err
	}

	if len(milkCollections) == 0 {
		return nil
	}

	latestMilkCollection := milkCollections[0]
	for _, collection := range milkCollections {
		if collection.Date.After(latestMilkCollection.Date) {
			latestMilkCollection = collection
		}
	}

	newBatch := models.GetBatchByLiters(latestMilkCollection.Liters)
	if animal.CurrentBatch != newBatch {
		animal.CurrentBatch = newBatch
		return s.animalRepository.Update(animal)
	}

	return nil
}
