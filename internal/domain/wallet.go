package domain

import "time"

type APIWalletCreateRequest struct {
	Name          string  `json:"name,omitempty"`
	InitialAmount float64 `json:"initial_amount"`
	Emoji         string  `json:"emoji"`
}

type APIWalletCreateResponse struct {
	ID            string    `json:"id"`
	Name          string    `json:"name"`
	SanitizedName string    `json:"sanitized_name"`
	Emoji         string    `json:"emoji"`
	Balance       float64   `json:"balance"`
	CreatedAt     time.Time `json:"created_at"`
}

type APIWalletGetResponse struct {
	ID            string    `json:"id"`
	Name          string    `json:"name"`
	SanitizedName string    `json:"sanitized_name"`
	Emoji         string    `json:"emoji"`
	Balance       float64   `json:"balance"`
	CreatedAt     time.Time `json:"created_at"`
}
