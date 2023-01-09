package domain

type UserStateDTO struct {
	State string         `json:"state"`
	Data  TransactionDTO `json:"data"`
}
