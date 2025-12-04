package repository

import (
	"github.com/fazendapro/FazendaPro-api/internal/cache"
)

type RepositoryFactory struct {
	db    *Database
	cache cache.CacheInterface
}

func NewRepositoryFactory(db *Database, cacheClient cache.CacheInterface) *RepositoryFactory {
	return &RepositoryFactory{
		db:    db,
		cache: cacheClient,
	}
}

func (f *RepositoryFactory) CreateAnimalRepository() AnimalRepositoryInterface {
	return NewAnimalRepository(f.db)
}

func (f *RepositoryFactory) CreateUserRepository() UserRepositoryInterface {
	return NewUserRepository(f.db)
}

func (f *RepositoryFactory) CreateMilkCollectionRepository() MilkCollectionRepositoryInterface {
	return NewMilkCollectionRepository(f.db.DB)
}

func (f *RepositoryFactory) CreateReproductionRepository() ReproductionRepositoryInterface {
	return NewReproductionRepository(f.db.DB)
}

func (f *RepositoryFactory) CreateRefreshTokenRepository() RefreshTokenRepositoryInterface {
	return NewRefreshTokenRepository(f.db)
}

func (f *RepositoryFactory) CreateFarmRepository() FarmRepositoryInterface {
	return NewFarmRepository(f.db)
}

func (f *RepositoryFactory) CreateSaleRepository() SaleRepository {
	return NewSaleRepository(f.db.DB)
}

func (f *RepositoryFactory) CreateDebtRepository() DebtRepositoryInterface {
	return NewDebtRepository(f.db.DB)
}

func (f *RepositoryFactory) CreateVaccineRepository() VaccineRepositoryInterface {
	return NewVaccineRepository(f.db.DB)
}

func (f *RepositoryFactory) CreateVaccineApplicationRepository() VaccineApplicationRepositoryInterface {
	return NewVaccineApplicationRepository(f.db.DB)
}

func (f *RepositoryFactory) GetCache() cache.CacheInterface {
	return f.cache
}
