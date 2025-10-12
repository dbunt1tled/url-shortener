package shorturl

import (
	"github.com/dbunt1tled/url-shortener/internal/domain/repository"
	"github.com/dbunt1tled/url-shortener/storage/mysql"
)

type URLRepository struct {
	*repository.BaseRepository[URL]
}

func NewURLRepository(db *mysql.Mysql) *URLRepository {
	return &URLRepository{
		BaseRepository: repository.NewBaseRepository[URL](db, "url"),
	}
}
