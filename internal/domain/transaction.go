package domain

import "time"

type TransactionDTO struct {
	ID          string     `json:"id,omitempty"`
	MsgID       int64      `json:"msg_id,omitempty"`
	UserID      int64      `json:"user_id,omitempty"`
	Amount      float64    `json:"amount"`
	Description string     `json:"description"`
	WalletID    int64      `json:"wallet_id"`
	CategoryID  int64      `json:"category_id"`
	CreatedAt   *time.Time `json:"created_at"`
}
