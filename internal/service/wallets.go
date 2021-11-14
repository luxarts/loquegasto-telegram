package service

import (
	"loquegasto-telegram/internal/domain"
	"loquegasto-telegram/internal/repository"
	"loquegasto-telegram/internal/utils/jwt"
	"time"
)

type WalletsService interface {
	Create(userID int, name string, balance float64, timestamp int64) error
	GetAll(userID int) (*[]domain.WalletDTO, error)
	GetByName(name string, userID int) (*domain.WalletDTO, error)
}
type walletsService struct {
	repo repository.WalletsRepository
}

func NewWalletsService(repo repository.WalletsRepository) WalletsService {
	return &walletsService{
		repo: repo,
	}
}
func (s *walletsService) Create(userID int, name string, balance float64, timestamp int64) error {
	token := jwt.GenerateToken(nil, &jwt.Payload{
		Subject: userID,
	})
	ts := time.Unix(timestamp, 0)
	walletDTO := domain.WalletDTO{
		UserID:    userID,
		Name:      name,
		Balance:   balance,
		CreatedAt: &ts,
	}
	return s.repo.Create(&walletDTO, token)
}
func (s *walletsService) GetAll(userID int) (*[]domain.WalletDTO, error) {
	token := jwt.GenerateToken(nil, &jwt.Payload{
		Subject: userID,
	})

	return s.repo.GetAll(token)
}
func (s *walletsService) GetByName(name string, userID int) (*domain.WalletDTO, error) {
	token := jwt.GenerateToken(nil, &jwt.Payload{
		Subject: userID,
	})

	return s.repo.GetByName(name, token)
}
