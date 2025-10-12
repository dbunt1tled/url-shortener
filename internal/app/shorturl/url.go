package shorturl

import (
	"errors"
	"time"

	"github.com/dbunt1tled/url-shortener/internal/domain/enum"
)

var (
	ErrURLNotFound = errors.New("url not found")
	ErrURLExists   = errors.New("url already exists")
	ErrURLInvalid  = errors.New("url invalid")
)

type URL struct {
	ID            int64       `db:"id" json:"id"`
	Code          string      `db:"code" json:"code"`
	URL           string      `db:"url" json:"url"`
	UserID        *int64      `db:"user_id" json:"userId"`
	Status        enum.Status `db:"status" json:"status"`
	CreatedAt     time.Time   `db:"created_at" json:"createdAt"`
	LastVisitedAt time.Time   `db:"last_visited_at" json:"lastVisitedAt"`
	ExpiredAt     time.Time   `db:"expired_at" json:"expiredAt"`
	Count         int64       `db:"count" json:"Count"`
}

type URLCreate struct {
	URL       string     `json:"url" vd:"@:(len($)<=2048 && regexp('^https?://[^\\s]+$')); msg:'not valid url'"`
	ExpiredAt *time.Time `json:"expired_at" validate:"omitempty"`
	UserID    *int64
}

type URLRedirect struct {
	Alias string `path:"alias" vd:"@:(len($)<=20 && regexp('^[A-Za-z0-9]+$')); msg:'not valid url'"`
}

type URLBase struct {
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}
