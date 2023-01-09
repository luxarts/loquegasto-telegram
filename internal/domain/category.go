package domain

type CategoryDTO struct {
	ID            int64  `json:"id"`
	UserID        int64  `json:"user_id"`
	Name          string `json:"name"`
	SanitizedName string `json:"sanitized_name"`
	Emoji         string `json:"emoji"`
}
