package domain

import "time"

type APIUsersCreateRequest struct {
	ChatID         int64 `json:"chat_id"`
	TimezoneOffset int   `json:"timezone_offset"`
}
type APIUsersCreateResponse struct {
	UserID         string    `json:"user_id"`
	CreatedAt      time.Time `json:"created_at"`
	ChatID         int       `json:"chat_id"`
	TimezoneOffset int       `json:"timezone_offset"`
}
type APIUsersAuthRequest struct {
	ChatID int64 `json:"chat_id"`
}
type APIUsersAuthResponse struct {
	AccessToken string `json:"access_token"`
}
