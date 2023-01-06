package service

import (
	"loquegasto-telegram/internal/domain"
	"loquegasto-telegram/internal/repository"
	"loquegasto-telegram/internal/utils/jwt"
	"time"
)

type WalletsService interface {
	Create(userID int64, name string, balance float64, timestamp *time.Time, token string) (*domain.WalletDTO, error)
	GetAll(userID int64) (*[]domain.WalletDTO, error)
	GetByName(name string, userID int64) (*domain.WalletDTO, error)
}
type walletsService struct {
	repo repository.WalletsRepository
}

func NewWalletsService(repo repository.WalletsRepository) WalletsService {
	return &walletsService{
		repo: repo,
	}
}
func (s *walletsService) Create(userID int64, name string, balance float64, timestamp *time.Time, token string) (*domain.WalletDTO, error) {
	walletDTO := domain.WalletDTO{
		UserID:    userID,
		Name:      name,
		Balance:   balance,
		CreatedAt: timestamp,
	}
	return s.repo.Create(&walletDTO, token)
}
func (s *walletsService) GetAll(userID int64) (*[]domain.WalletDTO, error) {
	token := jwt.GenerateToken(nil, &jwt.Payload{
		Subject: userID,
	})

	return s.repo.GetAll(token)
}
func (s *walletsService) GetByName(name string, userID int64) (*domain.WalletDTO, error) {
	token := jwt.GenerateToken(nil, &jwt.Payload{
		Subject: userID,
	})

	return s.repo.GetByName(name, token)
}
