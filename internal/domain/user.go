package domain

import "time"

type UserDTO struct {
	ID             int64      `json:"id"`
	CreatedAt      *time.Time `json:"created_at"`
	ChatID         int64      `json:"chat_id"`
	UpdatedAt      *time.Time `json:"updated_at"`
	TimezoneOffset int        `json:"timezone_offset"`
}
