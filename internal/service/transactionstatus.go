package service

import (
	"loquegasto-telegram/internal/domain"
	"loquegasto-telegram/internal/repository"
	"time"
)

type TransactionStatusService interface {
	Create(userID int64, amount float64, description string, createdAt time.Time, msgID int, status string) error
	GetByUserID(userID int64) (*domain.TransactionStatusDTO, error)
	UpdateByUserID(dto *domain.TransactionStatusDTO) error
	DeleteByUserID(userID int64) error
}

type transactionStatusService struct {
	repo repository.TransactionStatusRepository
}

func NewTransactionStatusService(repo repository.TransactionStatusRepository) TransactionStatusService {
	return &transactionStatusService{repo: repo}
}

func (s *transactionStatusService) Create(userID int64, amount float64, description string, createdAt time.Time, msgID int, status string) error {
	dto := &domain.TransactionStatusDTO{
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
func (s *transactionStatusService) GetByUserID(userID int64) (*domain.TransactionStatusDTO, error) {
	return s.repo.GetByUserID(userID)
}
func (s *transactionStatusService) UpdateByUserID(dto *domain.TransactionStatusDTO) error {
	return s.repo.UpdateByUserID(dto)
}
func (s *transactionStatusService) DeleteByUserID(userID int64) error {
	return s.repo.DeleteByUserID(userID)
}
