package database

import (
	"ShortUrl/internal/config"
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
)

func InitDB() *sqlx.DB {
	dbConfig := config.AppConfig.DB
	dns := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", dbConfig.Host, dbConfig.Port, dbConfig.User, dbConfig.Password, dbConfig.DBName, dbConfig.SSLMode)

	db, err := sqlx.Connect("postgres", dns)

	if err != nil {
		log.Fatal(err)
	}

	return db
}

func InitRedis() *redis.Client {
	redisConfig := config.AppConfig.Redis
	client := redis.NewClient(&redis.Options{
		Addr:     redisConfig.Addr,
		Password: redisConfig.Password,
		DB:       redisConfig.DB,
	})

	return client
}
