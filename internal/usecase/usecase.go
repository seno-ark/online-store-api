package usecase

import "online-store/internal/repository"

type Usecase struct {
	repo *repository.Repository
}

func NewUsecase(repo *repository.Repository) *Usecase {
	return &Usecase{repo: repo}
}
