package main

import (
	"ShortUrl/internal/config"
	"ShortUrl/internal/database"
	"ShortUrl/internal/routes"
	"ShortUrl/internal/services"
)

func main() {
	config.InitConfig()

	db := database.InitDB()
	redisClient := database.InitRedis()

	authService := &services.AuthService{DB: db}
	urlService := &services.URLService{DB: db, Cache: redisClient}

	r := routes.SetupRouter(authService, urlService)

	r.Run(":" + config.AppConfig.Server.Port)
}
