package service

import (
	"loquegasto-telegram/internal/domain"
	"loquegasto-telegram/internal/repository"
)

type TransactionsService interface {
	AddPayment(msgID int, amount int64, description string, source string) error
}

type transactionsService struct {
	repo repository.TransactionsRepository
}

func NewTransactionsService(repo repository.TransactionsRepository) TransactionsService {
	return &transactionsService{
		repo: repo,
	}
}

func (srv *transactionsService) AddPayment(msgID int, amount int64, description string, source string) error {
	// TODO Usar el msgID para actualizar la info en la DB en caso de que se edite el mensaje
	transactionDTO := domain.TransactionDTO{
		MsgID:       msgID,
		Amount:      amount,
		Description: description,
		Source:      source,
	}

	return srv.repo.Create(&transactionDTO)
}
