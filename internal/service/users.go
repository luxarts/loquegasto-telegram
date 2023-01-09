package service

import (
	"loquegasto-telegram/internal/domain"
	"loquegasto-telegram/internal/repository"
	"time"
)

type UsersService interface {
	Create(userID int64, timestamp *time.Time, chatID int64, token string) error
}
type usersService struct {
	repo repository.UsersRepository
}

func NewUsersService(repo repository.UsersRepository) UsersService {
	return &usersService{
		repo: repo,
	}
}
func (s *usersService) Create(userID int64, timestamp *time.Time, chatID int64, token string) error {
	userDTO := domain.UserDTO{
		ID:        userID,
		CreatedAt: timestamp,
		ChatID:    chatID,
		UpdatedAt: timestamp,
	}
	return s.repo.Create(&userDTO, token)
}
