package service

import (
	"loquegasto-telegram/internal/domain"
	"loquegasto-telegram/internal/repository"
	"loquegasto-telegram/internal/utils/jwt"
	"time"
)

type TransactionsService interface {
	AddTransaction(userID int64, msgID int, amount float64, description string, walletID int, categoryID int, timestamp *time.Time) error
	UpdateTransaction(userID int64, msgID int, amount float64, description string, walletID int) error
	GetAll(userID int64) (*[]domain.TransactionDTO, error)
}

type transactionsService struct {
	repo repository.TransactionsRepository
}

func NewTransactionsService(repo repository.TransactionsRepository) TransactionsService {
	return &transactionsService{
		repo: repo,
	}
}

func (srv *transactionsService) AddTransaction(userID int64, msgID int, amount float64, description string, walletID int, categoryID int, timestamp *time.Time) error {
	token := jwt.GenerateToken(nil, &jwt.Payload{
		Subject: userID,
	})

	transactionDTO := domain.TransactionDTO{
		MsgID:       msgID,
		UserID:      userID,
		Amount:      -amount,
		Description: description,
		WalletID:    walletID,
		CategoryID:  categoryID,
		CreatedAt:   timestamp,
	}

	return srv.repo.Create(&transactionDTO, token)
}
func (srv *transactionsService) UpdateTransaction(userID int64, msgID int, amount float64, description string, walletID int) error {
	token := jwt.GenerateToken(nil, &jwt.Payload{
		Subject: userID,
	})

	dto := domain.TransactionDTO{
		MsgID:       msgID,
		Amount:      -amount,
		Description: description,
		WalletID:    walletID,
	}

	return srv.repo.UpdateByMsgID(msgID, &dto, token)
}
func (srv *transactionsService) GetAll(userID int64) (*[]domain.TransactionDTO, error) {
	token := jwt.GenerateToken(nil, &jwt.Payload{
		Subject: userID,
	})

	r, err := srv.repo.GetAll(token)
	if err != nil {
		return nil, err
	}

	return r, nil
}
