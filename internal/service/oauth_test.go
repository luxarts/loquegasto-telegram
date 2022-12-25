package service

import (
	"github.com/stretchr/testify/suite"
	"loquegasto-telegram/internal/defines"
	"os"
	"testing"
)

type OAuthServiceTestSuite struct {
	suite.Suite
}

func (s *OAuthServiceTestSuite) SetupSuite() {
	err := os.Setenv(defines.EnvGoogleProjectID, "googleProjectID")
	if err != nil {
		panic(err)
	}
	err = os.Setenv(defines.EnvGoogleClientID, "googleClientID")
	if err != nil {
		panic(err)
	}
	err = os.Setenv(defines.EnvGoogleClientSecret, "googleClientSecret")
	if err != nil {
		panic(err)
	}
	err = os.Setenv(defines.EnvBackendBaseURL, "localhost:8080")
	if err != nil {
		panic(err)
	}
}

func (s *OAuthServiceTestSuite) TestGetLoginURL() {
	// Given
	expected := "https://accounts.google.com/o/oauth2/auth?access_type=offline&client_id=googleClientID&redirect_uri=https%3A%2F%2Flocalhost%3A8080%2Fauth%3Fuser_id%3D1234&response_type=code&scope=https%3A%2F%2Fwww.googleapis.com%2Fauth%2Fspreadsheets&state=state-token"
	var userID int64 = 1234

	srv := NewOAuthService()

	// When
	loginURL := srv.GetLoginURL(userID)

	// Then
	s.Equal(expected, loginURL)
}

func TestOAuthServiceTestSuite(t *testing.T) {
	suite.Run(t, new(OAuthServiceTestSuite))
}
