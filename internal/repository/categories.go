package repository

import (
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/luxarts/jsend-go"
	"loquegasto-telegram/internal/defines"
	"loquegasto-telegram/internal/domain"
	"loquegasto-telegram/internal/utils/maptostruct"
	"net/http"
	"os"
)

type CategoriesRepository interface {
	Create(req *domain.APICategoryCreateRequest, token string) (*domain.APICategoryCreateResponse, error)
	GetAll(token string) (*[]domain.APICategoryGetResponse, error)
	GetByID(ID string, token string) (*domain.APICategoryGetResponse, error)
}

type categoriesRepository struct {
	client  *resty.Client
	baseURL string
}

func NewCategoriesRepository(client *resty.Client) CategoriesRepository {
	return &categoriesRepository{
		client:  client,
		baseURL: os.Getenv(defines.EnvBackendBaseURL),
	}
}

func (r *categoriesRepository) Create(req *domain.APICategoryCreateRequest, token string) (*domain.APICategoryCreateResponse, error) {
	resp, err := r.client.R().
		SetAuthScheme("Bearer").
		SetAuthToken(token).
		SetBody(req).
		Post(fmt.Sprintf("%s%s", r.baseURL, defines.APICategoriesCreateURL))
	if err != nil {
		return nil, err
	}

	var body jsend.Body
	if err := json.Unmarshal(resp.Body(), &body); err != nil {
		return nil, err
	}
	if resp.StatusCode() != http.StatusCreated {
		return nil, jsend.NewError(body.Error(), err, resp.StatusCode())
	}

	var response domain.APICategoryCreateResponse
	err = maptostruct.Convert(body.Data, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

func (r *categoriesRepository) GetAll(token string) (*[]domain.APICategoryGetResponse, error) {
	resp, err := r.client.R().
		SetAuthScheme("Bearer").
		SetAuthToken(token).
		Get(fmt.Sprintf("%s%s", r.baseURL, defines.APICategoriesGetAllURL))
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
	var response []domain.APICategoryGetResponse
	err = json.Unmarshal(jsonBody, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}
func (r *categoriesRepository) GetByID(ID string, token string) (*domain.APICategoryGetResponse, error) {
	resp, err := r.client.R().
		SetAuthScheme("Bearer").
		SetAuthToken(token).
		SetPathParam(defines.ParamCategoryID, ID).
		Get(fmt.Sprintf("%s%s", r.baseURL, defines.APICategoriesGetByID))
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
	var response domain.APICategoryGetResponse
	err = json.Unmarshal(jsonBody, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}
