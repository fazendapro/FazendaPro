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
	animalRepo := f.repoFactory.CreateAnimalRepository()
	batchService := NewBatchService(animalRepo, milkCollectionRepo)
	return NewMilkCollectionService(milkCollectionRepo, batchService)
}

func (f *ServiceFactory) CreateReproductionService() *ReproductionService {
	reproductionRepo := f.repoFactory.CreateReproductionRepository()
	return NewReproductionService(reproductionRepo)
}

func (f *ServiceFactory) CreateFarmService() *FarmService {
	farmRepo := f.repoFactory.CreateFarmRepository()
	return NewFarmService(farmRepo)
}

func (f *ServiceFactory) CreateSaleService() SaleService {
	saleRepo := f.repoFactory.CreateSaleRepository()
	animalRepo := f.repoFactory.CreateAnimalRepository()
	return NewSaleService(saleRepo, animalRepo)
}

func (f *ServiceFactory) CreateDebtService() *DebtService {
	debtRepo := f.repoFactory.CreateDebtRepository()
	return NewDebtService(debtRepo)
}
