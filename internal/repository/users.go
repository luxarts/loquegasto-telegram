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

type UsersRepository interface {
	Create(transactionDTO *domain.UserDTO, token string) error
}

type usersRepository struct {
	client  *resty.Client
	baseURL string
}

func NewUsersRepository(client *resty.Client) UsersRepository {
	return &usersRepository{
		client:  client,
		baseURL: os.Getenv(defines.EnvBackendBaseURL),
	}
}

func (r *usersRepository) Create(userDTO *domain.UserDTO, token string) error {
	req := r.client.R()
	req = req.SetBody(userDTO)
	req = req.SetAuthScheme("Bearer")
	req = req.SetAuthToken(token)
	resp, err := req.Post(fmt.Sprintf("%s%s", r.baseURL, defines.APIUsersCreateURL))
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
