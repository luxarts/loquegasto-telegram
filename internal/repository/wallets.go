package repository

import (
	"encoding/json"
	"errors"
	"fmt"
	"loquegasto-telegram/internal/defines"
	"loquegasto-telegram/internal/domain"
	"loquegasto-telegram/internal/utils/maptostruct"
	"net/http"
	"os"

	"github.com/go-resty/resty/v2"
	"github.com/luxarts/jsend-go"
)

var ErrNotFound = errors.New("not found")

type WalletsRepository interface {
	Create(transactionDTO *domain.WalletDTO, token string) (*domain.WalletDTO, error)
	GetAll(token string) (*[]domain.WalletDTO, error)
	GetByName(name string, token string) (*domain.WalletDTO, error)
}

type walletsRepository struct {
	client  *resty.Client
	baseURL string
}

func NewWalletsRepository(client *resty.Client) WalletsRepository {
	return &walletsRepository{
		client:  client,
		baseURL: os.Getenv(defines.EnvBackendBaseURL),
	}
}

func (r *walletsRepository) Create(userDTO *domain.WalletDTO, token string) (*domain.WalletDTO, error) {
	req := r.client.R()
	req = req.SetBody(userDTO)
	req = req.SetAuthScheme("Bearer")
	req = req.SetAuthToken(token)
	resp, err := req.Post(fmt.Sprintf("%s%s", r.baseURL, defines.APIWalletsCreateURL))
	if err != nil {
		return nil, err
	}

	var body jsend.Body
	err = json.Unmarshal(resp.Body(), &body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode() != http.StatusCreated {
		return nil, jsend.NewError(body.Error(), err, resp.StatusCode())
	}

	var response domain.WalletDTO
	err = maptostruct.Convert(body.Data, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}
func (r *walletsRepository) GetAll(token string) (*[]domain.WalletDTO, error) {
	req := r.client.R()
	req = req.SetAuthScheme("Bearer")
	req = req.SetAuthToken(token)
	resp, err := req.Get(fmt.Sprintf("%s%s", r.baseURL, defines.APIWalletsGetAllURL))
	if err != nil {
		return nil, err
	}

	var body jsend.Body
	if err := json.Unmarshal(resp.Body(), &body); err != nil {
		return nil, err
	}

	// Convert map into struct
	jsonBody, err := json.Marshal(body.Data)
	if err != nil {
		return nil, err
	}
	var response []domain.WalletDTO
	err = json.Unmarshal(jsonBody, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}
func (r *walletsRepository) GetByName(name string, token string) (*domain.WalletDTO, error) {
	req := r.client.R()
	req = req.SetAuthScheme("Bearer")
	req = req.SetAuthToken(token)
	req = req.SetQueryParam(defines.ParamName, name)
	resp, err := req.Get(fmt.Sprintf("%s%s", r.baseURL, defines.APIWalletsGetAllURL))
	if err != nil {
		return nil, err
	}

	var body jsend.Body
	if err := json.Unmarshal(resp.Body(), &body); err != nil {
		return nil, err
	}

	if resp.StatusCode() == http.StatusNotFound {
		return nil, ErrNotFound
	}
	if resp.StatusCode() != http.StatusOK {
		return nil, errors.New(body.Error())
	}

	// Convert map into struct
	var response domain.WalletDTO
	err = maptostruct.Convert(body.Data, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}
