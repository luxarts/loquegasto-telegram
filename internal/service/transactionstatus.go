package service

import (
	"loquegasto-telegram/internal/domain"
	"loquegasto-telegram/internal/repository"
	"time"
)

type UserStateService interface {
	Create(userID int64, amount float64, description string, createdAt time.Time, msgID int, status string) error
	GetByUserID(userID int64) (*domain.UserStateDTO, error)
	UpdateByUserID(dto *domain.UserStateDTO) error
	DeleteByUserID(userID int64) error
}

type userStateService struct {
	repo repository.UserStateRepository
}

func NewUserStateService(repo repository.UserStateRepository) UserStateService {
	return &userStateService{repo: repo}
}

func (s *userStateService) Create(userID int64, amount float64, description string, createdAt time.Time, msgID int, status string) error {
	dto := &domain.UserStateDTO{
		Status: status,
		Data: domain.TransactionDTO{
			MsgID:       msgID,
			UserID:      userID,
			Amount:      amount,
			Description: description,
			CreatedAt:   &createdAt,
		},
	}
	return s.repo.Create(dto)
}
func (s *userStateService) GetByUserID(userID int64) (*domain.UserStateDTO, error) {
	return s.repo.GetByUserID(userID)
}
func (s *userStateService) UpdateByUserID(dto *domain.UserStateDTO) error {
	return s.repo.UpdateByUserID(dto)
}
func (s *userStateService) DeleteByUserID(userID int64) error {
	return s.repo.DeleteByUserID(userID)
}
