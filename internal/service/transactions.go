package service

import (
	"loquegasto-telegram/internal/domain"
	"loquegasto-telegram/internal/repository"
	"loquegasto-telegram/internal/utils/jwt"
	"time"
)

type TransactionsService interface {
	AddTransaction(msgID int64, amount float64, description string, walletID string, categoryID string, timestamp *time.Time, token string) error
	UpdateTransaction(userID int64, msgID int64, amount float64, description string, walletID string) error
	GetAll(userID int64, from *time.Time, to *time.Time) (*[]domain.APITransactionCreateRequest, error)
}

type transactionsService struct {
	repo repository.TransactionsRepository
}

func NewTransactionsService(repo repository.TransactionsRepository) TransactionsService {
	return &transactionsService{
		repo: repo,
	}
}

func (srv *transactionsService) AddTransaction(msgID int64, amount float64, description string, walletID string, categoryID string, timestamp *time.Time, token string) error {
	req := domain.APITransactionCreateRequest{
		MsgID:       msgID,
		Amount:      amount,
		Description: description,
		WalletID:    walletID,
		CategoryID:  categoryID,
		CreatedAt:   timestamp,
	}

	return srv.repo.Create(&req, token)
}
func (srv *transactionsService) UpdateTransaction(userID int64, msgID int64, amount float64, description string, walletID string) error {
	token := jwt.GenerateToken(nil, &jwt.Payload{
		Subject: userID,
	})

	dto := domain.APITransactionCreateRequest{
		MsgID:       msgID,
		Amount:      amount,
		Description: description,
		WalletID:    walletID,
	}

	return srv.repo.UpdateByMsgID(msgID, &dto, token)
}
func (srv *transactionsService) GetAll(userID int64, from *time.Time, to *time.Time) (*[]domain.APITransactionCreateRequest, error) {
	token := jwt.GenerateToken(nil, &jwt.Payload{
		Subject: userID,
	})

	r, err := srv.repo.GetAll(token, from, to)
	if err != nil {
		return nil, err
	}

	return r, nil
}
