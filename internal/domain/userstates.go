package domain

type UserStateDTO struct {
	State string `json:"state"`
	Data  any    `json:"data"`
}
