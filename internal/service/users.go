package service

import (
	"loquegasto-telegram/internal/defines"
	"loquegasto-telegram/internal/domain"
	"loquegasto-telegram/internal/repository"
	"loquegasto-telegram/internal/utils/jwt"
	"time"
)

type UsersService interface {
	Create(userID int64, timestamp *time.Time, chatID int64, groupIDs ...int64) error
	GetByID(userID int64) (*domain.UserDTO, error)
}
type usersService struct {
	repo repository.UsersRepository
}

func NewUsersService(repo repository.UsersRepository) UsersService {
	return &usersService{
		repo: repo,
	}
}
func (s *usersService) Create(userID int64, timestamp *time.Time, chatID int64, groupIDs ...int64) error {
	token := jwt.GenerateToken(nil, &jwt.Payload{
		Subject: userID,
	})
	userDTO := domain.UserDTO{
		ID:             userID,
		CreatedAt:      timestamp,
		ChatID:         chatID,
		UpdatedAt:      timestamp,
		TimezoneOffset: defines.DefaultUserTimeZone,
		GroupsIDs:      groupIDs,
	}
	return s.repo.Create(&userDTO, token)
}
func (s *usersService) GetByID(userID int64) (*domain.UserDTO, error) {
	token := jwt.GenerateToken(nil, &jwt.Payload{
		Subject: userID,
	})
	return s.repo.GetByID(token)
}
