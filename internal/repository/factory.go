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
