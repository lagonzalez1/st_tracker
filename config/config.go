package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DB   PostGresConfig
	Port string
	JWT  string
}

type PostGresConfig struct {
	Username string
	Password string
	URL      string
	Port     string
	Host     string
	Name     string
}

func LoadConfig() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
	}
	cfg := &Config{
		JWT:  os.Getenv("JWT"),
		Port: os.Getenv("B_PORT"),
		DB: PostGresConfig{
			Username: os.Getenv("POSTGRES_USER"),
			Password: os.Getenv("POSTGRES_PASSWORD"),
			URL:      os.Getenv("POSTGRES_URL"),
			Port:     os.Getenv("POSTGRES_PORT"),
			Host:     os.Getenv("POSTGRES_HOST"),
			Name:     os.Getenv("DB_NAME"),
		},
	}
	return cfg, nil
}
