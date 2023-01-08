package domain

type TransactionStatusDTO struct {
	Status string         `json:"status"`
	Data   TransactionDTO `json:"data"`
}
