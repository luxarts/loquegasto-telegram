package domain

import "time"

type WalletDTO struct {
	ID        int        `json:"id"`
	UserID    int        `json:"userID"`
	Name      string     `json:"name,omitempty"`
	Balance   float64    `json:"balance"`
	CreatedAt *time.Time `json:"created_at"`
}
