package repository

import (
	"encoding/json"
	"errors"
	"fmt"
	"loquegasto-telegram/internal/defines"
	"loquegasto-telegram/internal/domain"
	"net/http"
	"os"
	"strconv"

	"github.com/go-resty/resty/v2"
	"github.com/luxarts/jsend-go"
)

type TransactionsRepository interface {
	Create(transactionDTO *domain.TransactionDTO, token string) error
	GetAll(token string) (*[]domain.TransactionDTO, error)
	UpdateByMsgID(msgID int, transactionDTO *domain.TransactionDTO, token string) error
}

type transactionsRepository struct {
	client  *resty.Client
	baseURL string
}

func NewTransactionsRepository(client *resty.Client) TransactionsRepository {
	return &transactionsRepository{
		client:  client,
		baseURL: os.Getenv(defines.EnvBackendBaseURL),
	}
}

func (r *transactionsRepository) Create(transactionDTO *domain.TransactionDTO, token string) error {
	req := r.client.R()
	req = req.SetBody(transactionDTO)
	req = req.SetAuthScheme("Bearer")
	req = req.SetAuthToken(token)
	resp, err := req.Post(fmt.Sprintf("%s%s", r.baseURL, defines.APITransactionAddURL))
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
func (r *transactionsRepository) GetAll(token string) (*[]domain.TransactionDTO, error) {
	req := r.client.R()
	req = req.SetAuthScheme("Bearer")
	req = req.SetAuthToken(token)
	resp, err := req.Get(fmt.Sprintf("%s%s", r.baseURL, defines.APITransactionsGetAllURL))
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
	var response []domain.TransactionDTO
	err = json.Unmarshal(jsonBody, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}
func (r *transactionsRepository) UpdateByMsgID(msgID int, transactionDTO *domain.TransactionDTO, token string) error {
	req := r.client.R()
	req = req.SetBody(transactionDTO)
	req = req.SetAuthScheme("Bearer")
	req = req.SetAuthToken(token)
	req = req.SetPathParam(defines.ParamMsgID, strconv.Itoa(msgID))
	resp, err := req.Put(fmt.Sprintf("%s%s", r.baseURL, defines.APITransactionsUpdateURL))
	if err != nil {
		return err
	}

	var body jsend.Body
	err = json.Unmarshal(resp.Body(), &body)
	if err != nil {
		return err
	}

	if resp.StatusCode() != http.StatusOK {
		return errors.New(body.Error())
	}

	return nil
}
