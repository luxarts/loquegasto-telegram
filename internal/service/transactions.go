package service

import (
	"loquegasto-telegram/internal/domain"
	"loquegasto-telegram/internal/repository"
	"loquegasto-telegram/internal/utils/jwt"
	"strconv"
	"time"
)

type TransactionsService interface {
	AddPayment(userID int, msgID int, amount float64, description string, source string, timestamp int64) error
	UpdatePayment(userID int, msgID int, amount float64, description string, source string, timestamp int64) error
	GetTotal(userID int) (float64, error)
}

type transactionsService struct {
	repo repository.TransactionsRepository
}

func NewTransactionsService(repo repository.TransactionsRepository) TransactionsService {
	return &transactionsService{
		repo: repo,
	}
}

func (srv *transactionsService) AddPayment(userID int, msgID int, amount float64, description string, source string, timestamp int64) error {
	token := jwt.GenerateToken(nil, &jwt.Payload{
		Subject: userID,
	})

	transactionDTO := domain.TransactionDTO{
		MsgID:       msgID,
		Amount:      amount,
		Description: description,
		Source:      source,
		Timestamp:   time.Unix(timestamp, 0),
	}

	return srv.repo.Create(&transactionDTO, token)
}

func (srv *transactionsService) UpdatePayment(userID int, msgID int, amount float64, description string, source string, timestamp int64) error {
	token := jwt.GenerateToken(nil, &jwt.Payload{
		Subject: userID,
	})

	dto := domain.TransactionDTO{
		Amount:      amount,
		Description: description,
		Source:      source,
	}

	msgIDStr := strconv.Itoa(msgID)

	return srv.repo.UpdateByMsgID(msgIDStr, &dto, token)
}

func (srv *transactionsService) GetTotal(userID int) (float64, error) {
	token := jwt.GenerateToken(nil, &jwt.Payload{
		Subject: userID,
	})

	r, err := srv.repo.GetTotal(token)
	if err != nil {
		return 0, err
	}

	return r.Total, nil
}
