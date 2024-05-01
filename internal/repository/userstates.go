package repository

import (
	"context"
	"encoding/json"
	"github.com/go-redis/redis/v9"
	"loquegasto-telegram/internal/domain"
	"strconv"
)

type UserStateRepository interface {
	Create(userID int64, dto *domain.UserStateDTO) error
	GetByUserID(userID int64) (*domain.UserStateDTO, error)
	UpdateByUserID(userID int64, dto *domain.UserStateDTO) error
	DeleteByUserID(userID int64) error
}
type userStateRepository struct {
	rc *redis.Client
}

func NewUserStateRepository(rc *redis.Client) UserStateRepository {
	return &userStateRepository{
		rc: rc,
	}
}

func (r *userStateRepository) Create(userID int64, dto *domain.UserStateDTO) error {
	dtoBytes, err := json.Marshal(dto)
	if err != nil {
		return err
	}

	ctx := context.Background()
	return r.rc.Set(ctx, strconv.FormatInt(userID, 10), string(dtoBytes), 0).Err()
}
func (r *userStateRepository) GetByUserID(userID int64) (*domain.UserStateDTO, error) {
	uID := strconv.FormatInt(userID, 10)
	ctx := context.Background()
	res, err := r.rc.Get(ctx, uID).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return nil, err
	}

	var dto domain.UserStateDTO
	err = json.Unmarshal([]byte(res), &dto)
	if err != nil {
		return nil, err
	}

	return &dto, nil
}
func (r *userStateRepository) UpdateByUserID(userID int64, dto *domain.UserStateDTO) error {
	dtoBytes, err := json.Marshal(dto)
	if err != nil {
		return err
	}

	ctx := context.Background()
	return r.rc.Set(ctx, strconv.FormatInt(userID, 10), string(dtoBytes), 0).Err()
}
func (r *userStateRepository) DeleteByUserID(userID int64) error {
	ctx := context.Background()
	return r.rc.Del(ctx, strconv.FormatInt(userID, 10)).Err()
}
