package domain

import "time"

type APITransactionCreateRequest struct {
	MsgID       int64      `json:"msg_id"`
	Amount      float64    `json:"amount"`
	Description string     `json:"description"`
	WalletID    string     `json:"wallet_id"`
	CategoryID  string     `json:"category_id"`
	CreatedAt   *time.Time `json:"created_at"`
}
