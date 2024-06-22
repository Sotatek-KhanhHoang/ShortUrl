package services

import (
	"ShortUrl/internal/models"
	"context"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"

	"github.com/jmoiron/sqlx"
)

type URLService struct {
	DB    *sqlx.DB
	Cache *redis.Client
}

func (s *URLService) SaveUrl(original_url string, shortened_url string, user_id uuid.UUID) (string, error) {

	url := &models.Url{
		ID:            uuid.New(),
		Original_url:  original_url,
		Shortened_url: shortened_url,
		UserId:        user_id,
	}

	_, err := s.DB.NamedExec(`INSERT INTO urls (id, original_url, shortened_url, user_id) 
	VALUES (:id, :original_url, :shortened_url, :user_id)`, url)
	if err != nil {
		return "", err
	}

	s.Cache.Set(context.Background(), shortened_url, original_url, 0)

	return shortened_url, nil
}

func (s *URLService) GetOriginalURL(shortenedURL string) (string, error) {
	// Kiểm tra Redis
	originalURL, err := s.Cache.Get(context.Background(), shortenedURL).Result()
	if err == redis.Nil {
		// Kiểm tra Postgres
		var url models.Url
		err := s.DB.Get(&url, "SELECT original_url FROM urls WHERE shortened_url=$1", shortenedURL)
		if err != nil {
			return "", err
		}
		originalURL = url.Original_url
		// Lưu lại vào Redis
		s.Cache.Set(context.Background(), shortenedURL, originalURL, 0)
	} else if err != nil {
		return "", err
	}

	return originalURL, nil
}
