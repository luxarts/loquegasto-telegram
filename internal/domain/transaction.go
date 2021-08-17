package domain

type TransactionDTO struct {
	ID			string `json:"id,omitempty"`
	MsgID       int    `json:"msg_id"`
	Amount      int64  `json:"amount"`
	Description string `json:"description"`
	Source      string `json:"source,omitempty"`
}
