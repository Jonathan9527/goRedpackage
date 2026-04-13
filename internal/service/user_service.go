package service

import (
	"context"

	"learnGO/internal/model"
	"learnGO/internal/repository"
)

type UserService struct {
	userRepository *repository.UserRepository
}

func NewUserService(userRepository *repository.UserRepository) *UserService {
	return &UserService{
		userRepository: userRepository,
	}
}

func (s *UserService) FindByAccount(ctx context.Context, account string) (model.User, error) {
	return s.userRepository.FindByAccount(ctx, account)
}

func (s *UserService) List(ctx context.Context, limit int, offset int) ([]model.User, error) {
	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}
	if offset < 0 {
		offset = 0
	}

	return s.userRepository.List(ctx, limit, offset)
}
