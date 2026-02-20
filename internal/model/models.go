package model

import "time"

type Link struct {
	ShortURL    string    `db:"short_url"`
	OriginalURL string    `db:"original_url"`
	CustomAlias string    `db:"custom_alias"`
	CreatedAt   time.Time `db:"created_at"`
}

type Analytics struct {
	ID        int64     `db:"id"`
	ShortURL  string    `db:"short_url"`
	Timestamp time.Time `db:"timestamp"`
	UserAgent string    `db:"user_agent"`
	Referer   string    `db:"referer"`
}
