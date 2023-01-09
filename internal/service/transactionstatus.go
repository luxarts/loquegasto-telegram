package service

import (
	"loquegasto-telegram/internal/domain"
	"loquegasto-telegram/internal/repository"
	"time"
)

type UserStateService interface {
	SetState(userID int64, state string) error
	Create(userID int64, amount float64, description string, createdAt time.Time, msgID int64, status string) error
	GetByUserID(userID int64) (*domain.UserStateDTO, error)
	UpdateByUserID(userID int64, dto *domain.UserStateDTO) error
	DeleteByUserID(userID int64) error
}

type userStateService struct {
	repo repository.UserStateRepository
}

func NewUserStateService(repo repository.UserStateRepository) UserStateService {
	return &userStateService{repo: repo}
}

func (s *userStateService) SetState(userID int64, state string) error {
	usrStateDTO, err := s.repo.GetByUserID(userID)
	if err != nil {
		return err
	}

	// If not exists, create a new key with state
	if usrStateDTO == nil {
		usrStateDTO = &domain.UserStateDTO{}
	}

	usrStateDTO.State = state

	return s.repo.UpdateByUserID(userID, usrStateDTO)
}

func (s *userStateService) Create(userID int64, amount float64, description string, createdAt time.Time, msgID int64, status string) error {
	dto := &domain.UserStateDTO{
		State: status,
		Data: domain.TransactionDTO{
			MsgID:       msgID,
			UserID:      userID,
			Amount:      amount,
			Description: description,
			CreatedAt:   &createdAt,
		},
	}
	return s.repo.Create(dto)
}
func (s *userStateService) GetByUserID(userID int64) (*domain.UserStateDTO, error) {
	return s.repo.GetByUserID(userID)
}
func (s *userStateService) UpdateByUserID(userID int64, dto *domain.UserStateDTO) error {
	return s.repo.UpdateByUserID(userID, dto)
}
func (s *userStateService) DeleteByUserID(userID int64) error {
	return s.repo.DeleteByUserID(userID)
}
