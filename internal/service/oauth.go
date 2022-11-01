package service

import (
	"fmt"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"log"
	"loquegasto-telegram/internal/defines"
	"loquegasto-telegram/internal/domain"
	"net/url"
	"os"
)

type OAuthService interface {
	GetLoginURL(userID int) string
}

type oAuthService struct {
	creds domain.OAuthCredentials
}

func NewOAuthService() OAuthService {
	s := oAuthService{}

	clientID := os.Getenv(defines.EnvGoogleClientID)
	projectID := os.Getenv(defines.EnvGoogleProjectID)
	secret := os.Getenv(defines.EnvGoogleClientSecret)

	s.creds = domain.NewCredentials(clientID, projectID, secret)

	return &s
}

func (srv *oAuthService) GetLoginURL(userID int) string {
	redirectURI := url.URL{
		Scheme:   "https",
		Host:     os.Getenv(defines.EnvBackendBaseURL),
		Path:     defines.APIAuthURL,
		RawQuery: fmt.Sprintf("%s=%d", defines.QueryParamUserID, userID),
	}

	credsBytes := srv.creds.AddRedirectURI(redirectURI.String()).Bytes()

	config, err := google.ConfigFromJSON(credsBytes, defines.ScopeSheetsRW)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}

	return config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
}
