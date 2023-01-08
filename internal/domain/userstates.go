package domain

type UserStateDTO struct {
	Status string         `json:"status"`
	Data   TransactionDTO `json:"data"`
}
