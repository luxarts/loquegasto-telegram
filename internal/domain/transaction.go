package domain

import "time"

type TransactionDTO struct {
	ID          string     `json:"id,omitempty"`
	MsgID       int        `json:"msg_id,omitempty"`
	UserID      int64      `json:"user_id,omitempty"`
	Amount      float64    `json:"amount"`
	Description string     `json:"description"`
	WalletID    int        `json:"wallet_id"`
	CreatedAt   *time.Time `json:"created_at"`
}
