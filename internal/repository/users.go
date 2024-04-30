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

type UsersRepository interface {
	Create(req *domain.APIUsersCreateRequest) (*domain.APIUsersCreateResponse, error)
	Auth(req *domain.APIUsersAuthRequest) (*domain.APIUsersAuthResponse, error)
}

type usersRepository struct {
	clientID     string
	clientSecret string
	client       *resty.Client
	baseURL      string
}

func NewUsersRepository(client *resty.Client) UsersRepository {
	return &usersRepository{
		clientID:     os.Getenv(defines.EnvAPIClientID),
		clientSecret: os.Getenv(defines.EnvAPIClientSecret),
		client:       client,
		baseURL:      os.Getenv(defines.EnvBackendBaseURL),
	}
}

func (r *usersRepository) Create(req *domain.APIUsersCreateRequest) (*domain.APIUsersCreateResponse, error) {
	resp, err := r.client.R().
		SetHeader("client_id", r.clientID).
		SetHeader("client_secret", r.clientSecret).
		SetBody(req).
		Post(fmt.Sprintf("%s%s", r.baseURL, defines.APIUsersCreateURL))
	if err != nil {
		return nil, err
	}

	var jsendBody jsend.Body
	err = json.Unmarshal(resp.Body(), &jsendBody)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode() != http.StatusCreated {
		return nil, errors.New(jsendBody.Error())
	}

	// Convert map into struct
	var response domain.APIUsersCreateResponse
	err = maptostruct.Convert(jsendBody.Data, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}
func (r *usersRepository) Auth(req *domain.APIUsersAuthRequest) (*domain.APIUsersAuthResponse, error) {
	resp, err := r.client.R().
		SetHeader("client_id", r.clientID).
		SetHeader("client_secret", r.clientSecret).
		SetBody(req).
		Post(fmt.Sprintf("%s%s", r.baseURL, defines.APIUsersAuthURL))
	if err != nil {
		return nil, err
	}

	var jsendBody jsend.Body
	err = json.Unmarshal(resp.Body(), &jsendBody)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode() != http.StatusOK {
		return nil, errors.New(jsendBody.Error())
	}

	// Convert map into struct
	var response domain.APIUsersAuthResponse
	err = maptostruct.Convert(jsendBody.Data, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}
