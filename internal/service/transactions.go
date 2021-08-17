package service

import (
	"loquegasto-telegram/internal/domain"
	"loquegasto-telegram/internal/repository"
	"loquegasto-telegram/internal/utils/jwt"
	"time"
)

type TransactionsService interface {
	AddPayment(userID int, msgID int, amount int64, description string, source string, timestamp int64) error
}

type transactionsService struct {
	repo repository.TransactionsRepository
}

func NewTransactionsService(repo repository.TransactionsRepository) TransactionsService {
	return &transactionsService{
		repo: repo,
	}
}

func (srv *transactionsService) AddPayment(userID int, msgID int, amount int64, description string, source string, timestamp int64) error {
	// TODO Usar el msgID para actualizar la info en la DB en caso de que se edite el mensaje
	transactionDTO := domain.TransactionDTO{
		MsgID:       msgID,
		Amount:      amount,
		Description: description,
		Source:      source,
		Timestamp:   time.Unix(timestamp, 0),
	}

	token := jwt.GenerateToken(nil, &jwt.Payload{
		Subject: userID,
	})

	return srv.repo.Create(&transactionDTO, token)
}
