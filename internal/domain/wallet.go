package domain

import "time"

type WalletDTO struct {
	ID            int64      `json:"id"`
	UserID        int64      `json:"user_id"`
	Name          string     `json:"name,omitempty"`
	SanitizedName string     `json:"sanitized_name"`
	Balance       float64    `json:"balance"`
	CreatedAt     *time.Time `json:"created_at"`
}
