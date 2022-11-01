package domain

import (
	"encoding/json"
	"log"
)

const (
	oauthAuthUri      = "https://accounts.google.com/o/oauth2/auth"
	oauthTokenUri     = "https://oauth2.googleapis.com/token"
	oauthAuthProvider = "https://www.googleapis.com/oauth2/v1/certs"
)

type OAuthCredentials struct {
	Installed struct {
		ClientId                string   `json:"client_id"`
		ProjectId               string   `json:"project_id"`
		AuthUri                 string   `json:"auth_uri"`
		TokenUri                string   `json:"token_uri"`
		AuthProviderX509CertUrl string   `json:"auth_provider_x509_cert_url"`
		ClientSecret            string   `json:"client_secret"`
		RedirectUris            []string `json:"redirect_uris"`
	} `json:"installed"`
}

func NewCredentials(clientID string, projectID string, clientSecret string) OAuthCredentials {
	c := OAuthCredentials{}

	c.Installed.AuthUri = oauthAuthUri
	c.Installed.TokenUri = oauthTokenUri
	c.Installed.AuthProviderX509CertUrl = oauthAuthProvider
	c.Installed.ClientId = clientID
	c.Installed.ProjectId = projectID
	c.Installed.ClientSecret = clientSecret

	return c
}

func (c OAuthCredentials) AddRedirectURI(redirectURI string) *OAuthCredentials {
	c.Installed.RedirectUris = []string{redirectURI}
	return &c
}

func (c *OAuthCredentials) Bytes() []byte {
	b, err := json.Marshal(c)
	if err != nil {
		log.Fatalln(err)
	}
	return b
}
