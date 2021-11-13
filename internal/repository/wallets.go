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

type WalletsRepository interface {
	Create(transactionDTO *domain.WalletDTO, token string) error
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

func (r *walletsRepository) Create(userDTO *domain.WalletDTO, token string) error {
	req := r.client.R()
	req = req.SetBody(userDTO)
	req = req.SetAuthScheme("Bearer")
	req = req.SetAuthToken(token)
	resp, err := req.Post(fmt.Sprintf("%s%s", r.baseURL, defines.APIWalletsCreateURL))
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
