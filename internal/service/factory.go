package service

import (
	"github.com/fazendapro/FazendaPro-api/internal/repository"
)

type ServiceFactory struct {
	repoFactory *repository.RepositoryFactory
}

func NewServiceFactory(repoFactory *repository.RepositoryFactory) *ServiceFactory {
	return &ServiceFactory{repoFactory: repoFactory}
}

func (f *ServiceFactory) CreateAnimalService() *AnimalService {
	animalRepo := f.repoFactory.CreateAnimalRepository()
	return NewAnimalService(animalRepo)
}

func (f *ServiceFactory) CreateUserService() *UserService {
	userRepo := f.repoFactory.CreateUserRepository()
	return NewUserService(userRepo)
}

func (f *ServiceFactory) CreateMilkCollectionService() *MilkCollectionService {
	milkCollectionRepo := f.repoFactory.CreateMilkCollectionRepository()
	return NewMilkCollectionService(milkCollectionRepo)
}

func (f *ServiceFactory) CreateReproductionService() *ReproductionService {
	reproductionRepo := f.repoFactory.CreateReproductionRepository()
	return NewReproductionService(reproductionRepo)
}

func (f *ServiceFactory) CreateFarmService() *FarmService {
	farmRepo := f.repoFactory.CreateFarmRepository()
	return NewFarmService(farmRepo)
}
