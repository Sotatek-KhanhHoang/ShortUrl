package models

import "github.com/google/uuid"

type Url struct {
	ID            uuid.UUID `json:"id" db:"id"`
	Original_url  string    `json:"original_url" db:"original_url"`
	Shortened_url string    `json:"shortened_url" db:"shortened_url"`
	UserId        uuid.UUID `json:"user_id" db:"user_id"`
	ClickCount    int       `db:"click_count" json:"click_count"`
}
