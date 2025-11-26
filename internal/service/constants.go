package service

import "github.com/fazendapro/FazendaPro-api/internal/repository"

const (
	CacheKeyAnimalsFarm = "animals:farm:%d"

	ErrInvalidateCache = "Erro ao invalidar cache (não crítico): %v"
	ErrAnimalNotFound  = "animal not found"
)

var ErrSaleNotFoundOrNotBelongsToFarm = repository.ErrSaleNotFoundOrNotBelongsToFarm
