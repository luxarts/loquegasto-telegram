package domain

import "time"

type UserDTO struct {
	ID        int        `json:"id"`
	CreatedAt *time.Time `json:"created_at"`
	ChatID    int        `json:"chat_id"`
}
