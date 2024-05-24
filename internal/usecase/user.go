package usecase

import (
	"context"
	"errors"
	"online-store/internal/entity"
	"online-store/internal/repository"
	"online-store/pkg/utils"
)

func (u *Usecase) Register(ctx context.Context, payload entity.InUserRegister) (*entity.User, error) {
	_, err := u.repo.GetUserByEmail(ctx, payload.Email)
	if err == nil {
		return nil, utils.NewErrInvalidRequest("Email has already registered")
	}
	if !errors.Is(err, utils.ErrNotFound) {
		return nil, err
	}

	hashedPassword, err := utils.HashPassword(payload.Password)
	if err != nil {
		return nil, utils.NewErrInternalServer("Failed to register user")
	}

	err = u.repo.WithTx(ctx, func(rtx *repository.Repository) error {
		var errTx error

		_, errTx = rtx.CreateUser(ctx, entity.User{
			Email:    payload.Email,
			Password: hashedPassword,
			FullName: payload.FullName,
		})
		if errTx != nil {
			return errTx
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return u.repo.GetUserByEmail(ctx, payload.Email)
}

func (u *Usecase) Login(ctx context.Context, payload entity.InUserLogin) (*entity.User, error) {
	user, err := u.repo.GetUserByEmail(ctx, payload.Email)
	if err != nil {
		if errors.Is(err, utils.ErrNotFound) {
			return nil, utils.NewErrUnauthorized("Invalid email or password")
		}
		return nil, err
	}

	err = utils.CheckPassword(payload.Password, user.Password)
	if err != nil {
		return nil, utils.NewErrUnauthorized("Invalid email or password")
	}

	return user, nil
}

func (u *Usecase) GetUser(ctx context.Context, userID string) (*entity.User, error) {
	user, err := u.repo.GetUser(ctx, userID)
	if err != nil {
		if errors.Is(err, utils.ErrNotFound) {
			return nil, utils.NewErrUnauthorized("User not found")
		}
		return nil, err
	}

	return user, nil
}

func (u *Usecase) GetUserByEmail(ctx context.Context, email string) (*entity.User, error) {
	user, err := u.repo.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, utils.ErrNotFound) {
			return nil, utils.NewErrUnauthorized("User not found")
		}
		return nil, err
	}

	return user, nil
}
