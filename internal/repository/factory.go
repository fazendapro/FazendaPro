package repository

type RepositoryFactory struct {
	db *Database
}

func NewRepositoryFactory(db *Database) *RepositoryFactory {
	return &RepositoryFactory{db: db}
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
