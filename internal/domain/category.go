package domain

type CategoryDTO struct {
	ID            int    `json:"id"`
	UserID        int    `json:"user_id"`
	Name          string `json:"name"`
	SanitizedName string `json:"sanitized_name"`
	Emoji         string `json:"emoji"`
}
