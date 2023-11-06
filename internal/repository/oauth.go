package repository

import (
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"loquegasto-telegram/internal/defines"
	"os"
)

type OAuthRepository interface {
	GetLoginURL(token string) (*string, error)
}

type oAuthRepository struct {
	client  *resty.Client
	baseURL string
}

func NewOAuthRepository(client *resty.Client) OAuthRepository {
	return &oAuthRepository{
		client:  client,
		baseURL: os.Getenv(defines.EnvBackendBaseURL),
	}
}

func (repo *oAuthRepository) GetLoginURL(token string) (*string, error) {
	req := repo.client.R()
	req = req.SetAuthScheme("Bearer")
	req = req.SetAuthToken(token)

	resp, err := req.Get(fmt.Sprintf("%s%s", repo.baseURL, defines.APIOAuthGetLoginURL))
	if err != nil {
		return nil, err
	}

	type getURLResponse struct {
		URL string `json:"url"`
	}

	var body getURLResponse
	if err := json.Unmarshal(resp.Body(), &body); err != nil {
		return nil, err
	}

	return &body.URL, nil

}
