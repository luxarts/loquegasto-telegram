package service

import (
	"loquegasto-telegram/internal/repository"
	"loquegasto-telegram/internal/utils/jwt"
)

type OAuthService interface {
	GetLoginURL(userID int64) (*string, error)
}

type oAuthService struct {
	repo repository.OAuthRepository
}

func NewOAuthService(repo repository.OAuthRepository) OAuthService {
	return &oAuthService{
		repo: repo,
	}
}

func (svc *oAuthService) GetLoginURL(userID int64) (*string, error) {
	token := jwt.GenerateToken(nil, &jwt.Payload{
		Subject: userID,
	})

	return svc.repo.GetLoginURL(token)
}
