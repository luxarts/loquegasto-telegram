package service

import (
	"loquegasto-telegram/internal/domain"
	"loquegasto-telegram/internal/repository"
	"loquegasto-telegram/internal/utils/jwt"
	"time"
)

type UsersService interface {
	Create(userID int, timestamp int64, chatID int64) error
}
type usersService struct {
	repo repository.UsersRepository
}

func NewUsersService(repo repository.UsersRepository) UsersService {
	return &usersService{
		repo: repo,
	}
}
func (s *usersService) Create(userID int, timestamp int64, chatID int64) error {
	token := jwt.GenerateToken(nil, &jwt.Payload{
		Subject: userID,
	})
	ts := time.Unix(timestamp, 0)
	userDTO := domain.UserDTO{
		ID:        userID,
		CreatedAt: &ts,
		ChatID:    chatID,
	}
	return s.repo.Create(&userDTO, token)
}
