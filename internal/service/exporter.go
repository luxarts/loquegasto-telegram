package service

import (
	"loquegasto-telegram/internal/domain"
	"loquegasto-telegram/internal/repository"
	"time"
)

type ExporterService interface {
	Create(userID int64) error
	GetFile(userID int64) string
	Delete(userID int64) error
	AddEntry(id string, amount float64, description string, walletName string, categoryName string, createdAt *time.Time, userID int64) error
}
type exporterService struct {
	repo repository.ExporterRepository
}

func NewExporterService(repo repository.ExporterRepository) ExporterService {
	return &exporterService{repo: repo}
}
func (s *exporterService) Create(userID int64) error {
	return s.repo.Create(userID)
}
func (s *exporterService) GetFile(userID int64) string {
	return s.repo.GetFilePath(userID)
}
func (s *exporterService) Delete(userID int64) error {
	return s.repo.Delete(userID)
}
func (s *exporterService) AddEntry(id string, amount float64, description string, walletName string, categoryName string, createdAt *time.Time, userID int64) error {
	row := domain.ExporterRow{
		ID:           id,
		Amount:       amount,
		Description:  description,
		WalletName:   walletName,
		CategoryName: categoryName,
		Date:         createdAt,
		Time:         createdAt,
	}

	return s.repo.AddRow(row, userID)
}
