package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	DB      DatabaseConfig
	Redis   RedisConfig
	Server  ServerConfig
	Session SessionConfig
	JWT     JWTConfig
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

type RedisConfig struct {
	Addr     string
	Password string
	DB       int
}

type ServerConfig struct {
	Port string
}

type SessionConfig struct {
	SecretKey string
}

type JWTConfig struct {
	SecretKey string
}

var AppConfig Config

func InitConfig() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file, %s", err)
	}

	err := viper.Unmarshal(&AppConfig)
	if err != nil {
		log.Fatalf("Unable to decode into struct, %v", err)
	}
}
