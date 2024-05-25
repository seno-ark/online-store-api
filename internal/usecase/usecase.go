package usecase

import (
	"online-store/internal/repository"
	"online-store/pkg/cache"
)

type Usecase struct {
	repo  *repository.Repository
	cache *cache.Cache
}

func NewUsecase(repo *repository.Repository, cache *cache.Cache) *Usecase {
	return &Usecase{
		repo:  repo,
		cache: cache,
	}
}
