package service

import (
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"log"
	"loquegasto-telegram/internal/defines"
	"loquegasto-telegram/internal/domain"
	"net/url"
	"os"
	"strconv"
)

type OAuthService interface {
	GetLoginURL(userID int64) string
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

func (srv *oAuthService) GetLoginURL(userID int64) string {
	redirectURI := url.URL{
		Scheme: "http",
		Host:   os.Getenv(defines.EnvBackendBaseURL),
		Path:   defines.APIAuthURL + strconv.FormatInt(userID, 10),
	}

	credsBytes := srv.creds.AddRedirectURI(redirectURI.String()).Bytes()

	config, err := google.ConfigFromJSON(credsBytes, defines.ScopeSheetsRW)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}

	return config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
}
