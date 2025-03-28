package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DB     PostGresConfig
	Port   string
	JWT    string
	SignUp SignUpConfig
}

type PostGresConfig struct {
	Username string
	Password string
	URL      string
	Port     string
	Host     string
	Name     string
}

type SignUpConfig struct {
	ORG_ADD_KEY   string
	ORG_ADD_KEY_1 string
	ORG_ADD_KEY_2 string
	ORG_ADD_KEY_3 string
}

func LoadConfig() (*Config, error) {

	err := godotenv.Load(".env.development")
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
		SignUp: SignUpConfig{
			ORG_ADD_KEY:   os.Getenv("ORG_ADD_KEY"),
			ORG_ADD_KEY_1: os.Getenv("ORG_ADD_KEY_1"),
			ORG_ADD_KEY_2: os.Getenv("ORG_ADD_KEY_2"),
			ORG_ADD_KEY_3: os.Getenv("ORG_ADD_KEY_3"),
		},
	}
	return cfg, nil
}
