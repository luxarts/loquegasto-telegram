package service

import (
	"loquegasto-telegram/internal/domain"
	"loquegasto-telegram/internal/repository"
	"loquegasto-telegram/internal/utils/jwt"
	"time"
)

type TransactionsService interface {
	AddTransaction(userID int, msgID int, amount float64, description string, walletID int, timestamp int64) error
	UpdateTransaction(userID int, msgID int, amount float64, description string, walletID int) error
	GetAll(userID int) (*[]domain.TransactionDTO, error)
}

type transactionsService struct {
	repo repository.TransactionsRepository
}

func NewTransactionsService(repo repository.TransactionsRepository) TransactionsService {
	return &transactionsService{
		repo: repo,
	}
}

func (srv *transactionsService) AddTransaction(userID int, msgID int, amount float64, description string, walletID int, timestamp int64) error {
	token := jwt.GenerateToken(nil, &jwt.Payload{
		Subject: userID,
	})

	ts := time.Unix(timestamp, 0)
	transactionDTO := domain.TransactionDTO{
		MsgID:       msgID,
		UserID:      userID,
		Amount:      -amount,
		Description: description,
		WalletID:    walletID,
		CreatedAt:   &ts,
	}

	return srv.repo.Create(&transactionDTO, token)
}

func (srv *transactionsService) UpdateTransaction(userID int, msgID int, amount float64, description string, walletID int) error {
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

func (srv *transactionsService) GetAll(userID int) (*[]domain.TransactionDTO, error) {
	token := jwt.GenerateToken(nil, &jwt.Payload{
		Subject: userID,
	})

	r, err := srv.repo.GetAll(token)
	if err != nil {
		return nil, err
	}

	return r, nil
}
