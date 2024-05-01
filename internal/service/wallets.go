package service

import (
	"loquegasto-telegram/internal/domain"
	"loquegasto-telegram/internal/repository"
	"loquegasto-telegram/internal/utils/jwt"
	"time"
)

type WalletsService interface {
	Create(userID int64, name string, balance float64, timestamp *time.Time) (*domain.APIWalletCreateResponse, error)
	GetAll(token string) (*[]domain.APIWalletGetResponse, error)
	GetByName(name string, userID int64) (*domain.APIWalletGetResponse, error)
	GetByID(ID string, token string) (*domain.APIWalletGetResponse, error)
}
type walletsService struct {
	repo     repository.WalletsRepository
	sessRepo repository.SessionsRepository
}

func NewWalletsService(repo repository.WalletsRepository, sr repository.SessionsRepository) WalletsService {
	return &walletsService{
		repo:     repo,
		sessRepo: sr,
	}
}
func (s *walletsService) Create(userID int64, name string, balance float64, timestamp *time.Time) (*domain.APIWalletCreateResponse, error) {
	token := jwt.GenerateToken(nil, &jwt.Payload{
		Subject: userID,
	})
	walletDTO := domain.APIWalletCreateRequest{
		Name: name,
	}
	return s.repo.Create(&walletDTO, token)
}
func (s *walletsService) GetAll(token string) (*[]domain.APIWalletGetResponse, error) {
	return s.repo.GetAll(token)
}
func (s *walletsService) GetByName(name string, userID int64) (*domain.APIWalletGetResponse, error) {
	token := jwt.GenerateToken(nil, &jwt.Payload{
		Subject: userID,
	})

	return s.repo.GetByName(name, token)
}
func (s *walletsService) GetByID(ID string, token string) (*domain.APIWalletGetResponse, error) {
	return s.repo.GetByID(ID, token)
}
