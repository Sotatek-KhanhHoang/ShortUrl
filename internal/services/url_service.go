package services

import (
	"ShortUrl/internal/config"
	"ShortUrl/internal/models"
	"context"
	"crypto/rand"
	"encoding/base64"

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

	s.Cache.Set(context.Background(), config.AppConfig.URL.SecretKey+shortened_url, original_url, 0)

	return shortened_url, nil
}

func (s *URLService) SaveRandomUrl(original_url string, shortened_url string) (string, error) {

	s.Cache.Set(context.Background(), config.AppConfig.URL.SecretKey+shortened_url, original_url, 0)

	return shortened_url, nil
}

func (s *URLService) GetOriginalURL(shortenedURL string) (string, error) {
	// Kiểm tra Redis
	originalURL, err := s.Cache.Get(context.Background(), config.AppConfig.URL.SecretKey+shortenedURL).Result()
	if err == redis.Nil {
		// Kiểm tra Postgres
		var url models.Url
		err := s.DB.Get(&url, "SELECT original_url FROM urls WHERE shortened_url=$1", shortenedURL)
		if err != nil {
			return "", err
		}
		originalURL = url.Original_url
		// Lưu lại vào Redis
		s.Cache.Set(context.Background(), config.AppConfig.URL.SecretKey+shortenedURL, originalURL, 0)
	} else if err != nil {
		return "", err
	}

	return originalURL, nil
}

func (s *URLService) GetById(userId uuid.UUID) ([]models.Url, error) {
	var urls []models.Url

	err := s.DB.Select(&urls, "SELECT * FROM urls WHERE user_id = $1", userId)
	if err != nil {
		return nil, err
	}

	for i := range urls {
		urls[i].Shortened_url = config.AppConfig.URL.SecretKey + urls[i].Shortened_url
	}

	return urls, nil
}

func (s *URLService) IncreaseClick(ShortURL string) error {
	_, err := s.DB.Exec("UPDATE urls SET click_count = click_count + 1 WHERE shortened_url = $1", ShortURL)
	return err
}

func (s *URLService) RandomUrl(n int) (string, error) {
	//Tạo 1 slide có độ dài n
	b := make([]byte, n)

	//Đọc ngẫu nhiên n byte vào slide b
	_, err := rand.Read(b)
	if err != nil {
		return "", nil
	}

	return base64.URLEncoding.EncodeToString(b), nil
}

func (s *URLService) GetAllUrlFromCache() (map[string]string, error) {
	c := context.Background()
	urls := make(map[string]string)

	iter := s.Cache.Scan(c, 0, "*", 0).Iterator()
	for iter.Next(c) {
		key := iter.Val()

		value, err := s.Cache.Get(c, key).Result()

		if err != nil {
			return nil, err
		}

		urls[key] = value

	}
	if err := iter.Err(); err != nil {
		return nil, err
	}

	return urls, nil

}
