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
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/luxarts/jsend-go"
)

type TransactionsRepository interface {
	Create(transactionDTO *domain.APITransactionCreateRequest, token string) error
	GetAll(token string, from *time.Time, to *time.Time) (*[]domain.APITransactionCreateRequest, error)
	UpdateByMsgID(msgID int64, transactionDTO *domain.APITransactionCreateRequest, token string) error
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

func (r *transactionsRepository) Create(req *domain.APITransactionCreateRequest, token string) error {
	resp, err := r.client.R().
		SetAuthScheme("Bearer").
		SetAuthToken(token).
		SetBody(req).
		Post(fmt.Sprintf("%s%s", r.baseURL, defines.APITransactionAddURL))
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
func (r *transactionsRepository) GetAll(token string, from *time.Time, to *time.Time) (*[]domain.APITransactionCreateRequest, error) {
	req := r.client.R()
	req = req.SetAuthScheme("Bearer")
	req = req.SetAuthToken(token)

	if from != nil {
		req = req.SetQueryParam("from", strconv.FormatInt(from.UTC().Unix(), 10))
	}
	if to != nil {
		req = req.SetQueryParam("to", strconv.FormatInt(to.UTC().Unix(), 10))
	}

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
	var response []domain.APITransactionCreateRequest
	err = json.Unmarshal(jsonBody, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}
func (r *transactionsRepository) UpdateByMsgID(msgID int64, transactionDTO *domain.APITransactionCreateRequest, token string) error {
	req := r.client.R()
	req = req.SetBody(transactionDTO)
	req = req.SetAuthScheme("Bearer")
	req = req.SetAuthToken(token)
	req = req.SetPathParam(defines.ParamMsgID, strconv.FormatInt(msgID, 10))
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
