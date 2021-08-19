package repository

import (
	"encoding/json"
	"errors"
	"fmt"
	"loquegasto-telegram/internal/defines"
	"loquegasto-telegram/internal/domain"
	"net/http"
	"os"

	"github.com/go-resty/resty/v2"
	"github.com/luxarts/jsend-go"
)

type TransactionsRepository interface {
	Create(transactionDTO *domain.TransactionDTO, token string) error
	GetTotal(token string) (*domain.TotalDTO, error)
}

type transactionsRepository struct {
	client  *resty.Client
	baseURL string
}

func NewTransactionsRepository(client *resty.Client) TransactionsRepository {
	return &transactionsRepository{
		client:  client,
		baseURL: os.Getenv(defines.EnvTransactionsBaseURL),
	}
}

func (r *transactionsRepository) Create(transactionDTO *domain.TransactionDTO, token string) error {
	req := r.client.R()
	req = req.SetBody(transactionDTO)
	req = req.SetAuthScheme("Bearer")
	req = req.SetAuthToken(token)
	resp, err := req.Post(fmt.Sprintf("%s%s", r.baseURL, defines.APITransactionPostURL))
	if err != nil {
		return err
	}

	var body jsend.Body
	err = json.Unmarshal(resp.Body(), &body)
	if err != nil {
		return err
	}

	if resp.StatusCode() != http.StatusCreated {
		return errors.New(body.Error())
	}

	return nil
}
func (r *transactionsRepository) GetTotal(token string) (*domain.TotalDTO, error) {
	req := r.client.R()
	req = req.SetAuthScheme("Bearer")
	req = req.SetAuthToken(token)
	resp, err := req.Get(fmt.Sprintf("%s%s", r.baseURL, defines.APITransactionsGetTotal))
	if err != nil {
		return nil, err
	}

	var body jsend.Body
	err = json.Unmarshal(resp.Body(), &body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode() != http.StatusOK {
		return nil, errors.New(body.Error())
	}

	// Convert map into struct
	jsonBody, err := json.Marshal(body.Data)
	if err != nil {
		return nil, err
	}
	var totalDTO domain.TotalDTO
	err = json.Unmarshal(jsonBody, &totalDTO)
	if err != nil {
		return nil, err
	}

	return &totalDTO, nil
}
