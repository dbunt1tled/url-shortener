package storage

import (
	"go_first/internal/lib/common/model/urlshort"
	"go_first/internal/lib/common/model/user"
	"go_first/storage/mysql"
)

type Storage interface {
	CreateURL(s string, alias string) (*urlshort.URLShort, error)
	GetURL(filter mysql.URLFilter) (*urlshort.URLShort, error)
	CreateUser(firstName string, secondName string, email string, phoneNumber string, password string, status int) (*user.User, error)
	GetUser(filter mysql.UserFilter) (*user.User, error)
	GetUserIdentity(login string) (*user.User, error)
}
