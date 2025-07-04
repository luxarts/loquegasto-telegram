package domain

import "time"

type APICategoryCreateRequest struct {
	Name  string `json:"name"`
	Emoji string `json:"emoji"`
}
type APICategoryCreateResponse struct {
	ID            string    `json:"id"`
	Name          string    `json:"name"`
	SanitizedName string    `json:"sanitized_name"`
	Emoji         string    `json:"emoji"`
	CreatedAt     time.Time `json:"created_at"`
}
type APICategoryGetResponse struct {
	ID            string    `json:"id"`
	Name          string    `json:"name"`
	SanitizedName string    `json:"sanitized_name"`
	Emoji         string    `json:"emoji"`
	CreatedAt     time.Time `json:"created_at"`
}
