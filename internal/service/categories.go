package service

import (
	"loquegasto-telegram/internal/domain"
	"loquegasto-telegram/internal/repository"
	"loquegasto-telegram/internal/utils/jwt"
)

type CategoriesService interface {
	GetAll(token string) (*[]domain.APICategoryGetResponse, error)
	GetByID(ID string, token string) (*domain.APICategoryGetResponse, error)
	Create(userID int64, name string, emoji string) (*domain.APICategoryCreateResponse, error)
}
type categoriesService struct {
	repo repository.CategoriesRepository
}

func NewCategoriesService(repo repository.CategoriesRepository) CategoriesService {
	return &categoriesService{repo: repo}
}

func (s *categoriesService) GetAll(token string) (*[]domain.APICategoryGetResponse, error) {
	return s.repo.GetAll(token)
}
func (s *categoriesService) GetByID(ID string, token string) (*domain.APICategoryGetResponse, error) {
	return s.repo.GetByID(ID, token)
}
func (s *categoriesService) Create(userID int64, name string, emoji string) (*domain.APICategoryCreateResponse, error) {
	token := jwt.GenerateToken(nil, &jwt.Payload{
		Subject: userID,
	})

	dto := &domain.APICategoryCreateRequest{
		Name:  name,
		Emoji: emoji,
	}
	return s.repo.Create(dto, token)
}
