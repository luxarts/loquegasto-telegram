package repository

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	"loquegasto-telegram/internal/defines"
	"loquegasto-telegram/internal/domain"
	"os"
)

type TransactionsRepository interface {
	Create(transactionDTO *domain.TransactionDTO) error
}

type transactionsRepository struct {
	client *resty.Client
	baseURL string
}

func NewTransactionsRepository(client *resty.Client) TransactionsRepository {
	return &transactionsRepository{
		client: client,
		baseURL: os.Getenv(defines.EnvTransactionsBaseURL),
	}
}

func (r *transactionsRepository) Create(transactionDTO *domain.TransactionDTO) error{
	req := r.client.R()
	req = req.SetBody(transactionDTO)
	_, err := req.Post(fmt.Sprintf("%s%s",r.baseURL, defines.APITransactionPostURL))

	if err != nil {
		return err
	}

	return nil
}