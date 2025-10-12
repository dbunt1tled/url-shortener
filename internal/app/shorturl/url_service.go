package shorturl

import (
	"context"
	"fmt"
	"net/url"
	"time"

	"github.com/dbunt1tled/url-shortener/internal/domain/enum"
	"github.com/dbunt1tled/url-shortener/internal/domain/repository"
	"github.com/dbunt1tled/url-shortener/internal/lib/hasher"
)

type URLService struct {
	urlRepository *URLRepository
	hasher        *hasher.Hasher
	baseURL       string
}

func NewURLService(urlRepository *URLRepository, hasher *hasher.Hasher, baseURL string) *URLService {
	return &URLService{
		urlRepository: urlRepository,
		hasher:        hasher,
		baseURL:       baseURL,
	}
}

func (s *URLService) ShortenURL(ctx context.Context, dto *URLCreate) (*URLBase, error) {
	var (
		url, m    *URL
		exp       time.Time
		shortCode string
		i         int
		err       error
	)
	if !isValidURL(dto.URL) {
		return nil, ErrURLInvalid
	}

	now := time.Now()
	if dto.ExpiredAt != nil {
		exp = *dto.ExpiredAt
	} else {
		exp = now.Add(time.Hour * 24 * 31) //nolint:mnd // 1 month
	}
	url = &URL{
		URL:       dto.URL,
		Code:      shortCode,
		UserID:    dto.UserID,
		CreatedAt: now,
		ExpiredAt: exp,
		Status:    enum.Active,
		Count:     0,
	}
	i = 0
	for {
		shortCode, err = s.hasher.NewID()
		if err != nil {
			return nil, err
		}
		url.Code = shortCode
		m, _ = s.urlRepository.Create(ctx, url)
		if m != nil {
			break
		}
		i++
		if i > 5 { //nolint:mnd // max attempts
			return nil, ErrURLExists
		}
	}

	return &URLBase{
		ShortURL:    s.GetShortURL(m),
		OriginalURL: m.URL,
	}, nil
}

func (s *URLService) One(ctx context.Context, filters []repository.Filter, sort *repository.Sort) (*URL, error) {
	return s.urlRepository.One(ctx, filters, sort)
}

func (s *URLService) Delete(ctx context.Context, id int64) (*URL, error) {
	return s.urlRepository.Update(ctx, id, map[string]any{
		"status": enum.Inactive,
	})
}

func (s *URLService) Update(ctx context.Context, id int64, data map[string]any) (*URL, error) {
	return s.urlRepository.Update(ctx, id, data)
}

func (s *URLService) GetShortURL(url *URL) string {
	if url == nil {
		return ""
	}
	return fmt.Sprintf("%s/%s", s.baseURL, url.Code)
}

func isValidURL(str string) bool {
	u, err := url.Parse(str)
	return err == nil && (u.Scheme == "http" || u.Scheme == "https") && u.Host != ""
}
