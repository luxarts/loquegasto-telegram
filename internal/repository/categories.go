package repository

import (
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/luxarts/jsend-go"
	"loquegasto-telegram/internal/defines"
	"loquegasto-telegram/internal/domain"
	"os"
	"strconv"
)

type CategoriesRepository interface {
	GetAll(token string) (*[]domain.CategoryDTO, error)
	GetByID(ID int, token string) (*domain.CategoryDTO, error)
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

func (r *categoriesRepository) GetAll(token string) (*[]domain.CategoryDTO, error) {
	req := r.client.R()
	req = req.SetAuthScheme("Bearer")
	req = req.SetAuthToken(token)
	resp, err := req.Get(fmt.Sprintf("%s%s", r.baseURL, defines.APICategoriesGetAllURL))
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
	var response []domain.CategoryDTO
	err = json.Unmarshal(jsonBody, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}
func (r *categoriesRepository) GetByID(ID int, token string) (*domain.CategoryDTO, error) {
	req := r.client.R()
	req = req.SetAuthScheme("Bearer")
	req = req.SetAuthToken(token)
	req = req.SetPathParam(defines.ParamCategoryID, strconv.Itoa(ID))
	resp, err := req.Get(fmt.Sprintf("%s%s", r.baseURL, defines.APICategoriesGetByID))
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
	var response domain.CategoryDTO
	err = json.Unmarshal(jsonBody, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}
