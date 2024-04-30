package service

import (
	"loquegasto-telegram/internal/defines"
	"loquegasto-telegram/internal/domain"
	"loquegasto-telegram/internal/repository"
)

type UsersService interface {
	Create(userID int64) error
}
type usersService struct {
	usersRepo  repository.UsersRepository
	sessRepo   repository.SessionsRepository
	walletRepo repository.WalletsRepository
	catRepo    repository.CategoriesRepository
}

func NewUsersService(ur repository.UsersRepository, sr repository.SessionsRepository, wr repository.WalletsRepository, cr repository.CategoriesRepository) UsersService {
	return &usersService{
		usersRepo:  ur,
		sessRepo:   sr,
		walletRepo: wr,
		catRepo:    cr,
	}
}
func (s *usersService) Create(userID int64) error {
	// Create user
	_, err := s.usersRepo.Create(&domain.APIUsersCreateRequest{
		ChatID:         userID,
		TimezoneOffset: defines.DefaultUserTimeZone,
	})
	if err != nil {
		return err
	}

	// Auth user
	authResp, err := s.usersRepo.Auth(&domain.APIUsersAuthRequest{
		ChatID: userID,
	})

	if err != nil {
		return err
	}

	// Create session
	s.sessRepo.Add(userID, authResp.AccessToken)

	// Create default wallet
	_, err = s.walletRepo.Create(&domain.APIWalletCreateRequest{
		Name:          defines.DefaultWalletName,
		InitialAmount: 0,
		Emoji:         defines.DefaultWalletEmoji,
	}, authResp.AccessToken)

	if err != nil {
		return err
	}

	// Create default category
	_, err = s.catRepo.Create(&domain.APICategoryCreateRequest{
		Name:  defines.DefaultCategoryNameOthers,
		Emoji: defines.DefaultCategoryEmojiOthers,
	}, authResp.AccessToken)

	return err
}
