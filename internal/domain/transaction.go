package domain

import "time"

type TransactionDTO struct {
	ID          string    `json:"id,omitempty"`
	UserID      int       `json:"user_id,omitempty"`
	MsgID       int       `json:"msg_id"`
	Amount      float64   `json:"amount"`
	Description string    `json:"description"`
	Source      string    `json:"source,omitempty"`
	Timestamp   time.Time `json:"created_at"`
}
