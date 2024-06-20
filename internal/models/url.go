package models

import "github.com/google/uuid"

type Url struct {
	ID            uuid.UUID `db:"id"`
	Original_url  string    `db:"original_url"`
	Shortened_url string    `db:"shortened_url"`
	UserId        uuid.UUID `db:"user_id"`
}
