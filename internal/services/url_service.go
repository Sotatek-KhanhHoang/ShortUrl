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

func NewPostgres(db *sqlx.DB) *URLService {
	return &URLService{DB: db}
}

func (repo *URLService) SaveUrl(originalUrl string, shortendUrl string, userid uuid.UUID) (*models.Url, error) {
	url := &models.Url{
		ID:            uuid.New(),
		Original_url:  originalUrl,
		Shortened_url: shortendUrl,
		UserId:        userid,
	}

	_, err := repo.DB.NamedExec(`INSERT INTO urls (id, original_url, shortened_url, user_id) 
							VALUES (:id, :originalUrl, :shortendUrl, :userid")`, url)
	if err != nil {
		return nil, err
	}

	repo.Cache.Set(context.Background(), shortendUrl, originalUrl, 0)

	return url, nil
}

func (repo *URLService) GetUrlById(userid uuid.UUID) (string, error) {
	return "", nil
}
