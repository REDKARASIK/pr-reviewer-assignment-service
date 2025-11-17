package service

import (
	"context"
	"pr-reviewer-assigment-service/internal/domain"
)

type UserRepository interface {
	GetByID(ctx context.Context, id string) (*domain.User, error)
	UpdateActive(ctx context.Context, user *domain.User) error
	Create(ctx context.Context, userID, name string, isActive *bool) (*domain.User, error)
}

type UserService struct {
	repo UserRepository
}

func NewUserService(repo UserRepository) *UserService {
	return &UserService{
		repo: repo,
	}
}

func (service *UserService) SetIsActive(ctx context.Context, userID string, isActive bool) (*domain.User, error) {
	user, err := service.repo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	user.IsActive = isActive

	if err = service.repo.UpdateActive(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}
