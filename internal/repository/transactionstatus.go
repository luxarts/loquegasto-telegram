package repository

import (
	"context"
	"encoding/json"
	"github.com/go-redis/redis/v9"
	"loquegasto-telegram/internal/domain"
	"strconv"
)

type TransactionStatusRepository interface {
	Create(dto *domain.TransactionStatusDTO) error
	GetByUserID(userID int64) (*domain.TransactionStatusDTO, error)
	UpdateByUserID(dto *domain.TransactionStatusDTO) error
	DeleteByUserID(userID int64) error
}
type transactionStatusRepository struct {
	rc *redis.Client
}

func NewTransactionStatusRepository(rc *redis.Client) TransactionStatusRepository {
	return &transactionStatusRepository{
		rc: rc,
	}
}

func (r *transactionStatusRepository) Create(dto *domain.TransactionStatusDTO) error {
	dtoBytes, err := json.Marshal(dto)
	if err != nil {
		return err
	}

	ctx := context.Background()
	return r.rc.Set(ctx, strconv.FormatInt(dto.Data.UserID, 10), string(dtoBytes), 0).Err()
}
func (r *transactionStatusRepository) GetByUserID(userID int64) (*domain.TransactionStatusDTO, error) {
	uID := strconv.FormatInt(userID, 10)
	ctx := context.Background()
	res, err := r.rc.Get(ctx, uID).Result()
	if err != nil {
		return nil, err
	}

	var dto domain.TransactionStatusDTO
	err = json.Unmarshal([]byte(res), &dto)
	if err != nil {
		return nil, err
	}

	return &dto, nil
}
func (r *transactionStatusRepository) UpdateByUserID(dto *domain.TransactionStatusDTO) error {
	dtoBytes, err := json.Marshal(dto)
	if err != nil {
		return err
	}

	ctx := context.Background()
	return r.rc.Set(ctx, strconv.FormatInt(dto.Data.UserID, 10), string(dtoBytes), 0).Err()
}
func (r *transactionStatusRepository) DeleteByUserID(userID int64) error {
	ctx := context.Background()
	return r.rc.Del(ctx, strconv.FormatInt(userID, 10)).Err()
}
