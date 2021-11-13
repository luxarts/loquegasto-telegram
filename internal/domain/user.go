package domain

import "time"

type UserDTO struct {
	ID        int        `json:"id"`
	CreatedAt *time.Time `json:"created_at"`
	ChatID    int64      `json:"chat_id"`
}
