package service

import (
	"loquegasto-telegram/internal/domain"
	"loquegasto-telegram/internal/repository"
	"loquegasto-telegram/internal/utils/jwt"
)

type CategoriesService interface {
	GetAll(userID int64) (*[]domain.CategoryDTO, error)
	GetByID(ID int, userID int64) (*domain.CategoryDTO, error)
	Create(userID int64, name string, emoji string) (*domain.CategoryDTO, error)
}
type categoriesService struct {
	repo repository.CategoriesRepository
}

func NewCategoriesService(repo repository.CategoriesRepository) CategoriesService {
	return &categoriesService{repo: repo}
}

func (s *categoriesService) GetAll(userID int64) (*[]domain.CategoryDTO, error) {
	token := jwt.GenerateToken(nil, &jwt.Payload{
		Subject: userID,
	})

	return s.repo.GetAll(token)
}
func (s *categoriesService) GetByID(ID int, userID int64) (*domain.CategoryDTO, error) {
	token := jwt.GenerateToken(nil, &jwt.Payload{
		Subject: userID,
	})

	return s.repo.GetByID(ID, token)
}
func (s *categoriesService) Create(userID int64, name string, emoji string) (*domain.CategoryDTO, error) {
	token := jwt.GenerateToken(nil, &jwt.Payload{
		Subject: userID,
	})

	dto := &domain.CategoryDTO{
		Name:  name,
		Emoji: emoji,
	}
	return s.repo.Create(dto, token)
}
