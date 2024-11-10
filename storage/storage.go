package storage

import (
	"github.com/pkg/errors"
	"go_first/storage/mysql"
)

var (
	ErrUrlNotFound = errors.New("url not found")
	ErrUrlExists   = errors.New("url exists")
)

type Storage interface {
	CreateURL(s string, alias string) (int64, error)
	GetURL(filter mysql.URLFilter) (*mysql.URL, error)
	UpdateURL(id int64, alias string) (*mysql.URL, error)
	DeleteURL(id int64) error
}
