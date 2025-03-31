package config

import (
	"fmt"
	"log"
	"os"
	"tracker/app/models"

	"github.com/joho/godotenv"
)

func LoadConfig() (*models.Config, error) {
	env := os.Getenv("APP_ENV")
	fmt.Println("Current ENV---", env)
	if env == "" {
		env = "development"
	}
	envFile := ".env." + env
	if env == "production" {
		err := awsConfig()
		if err != nil {
			log.Println("Error creating .env file")
		}
	}

	err := godotenv.Load(envFile)
	if err != nil {
		log.Println("Error loading .env file")
	}
	cfg := &models.Config{
		JWT:  os.Getenv("JWT"),
		Port: os.Getenv("B_PORT"),
		DB: models.PostGresConfig{
			Username: os.Getenv("POSTGRES_USER"),
			Password: os.Getenv("POSTGRES_PASSWORD"),
			URL:      os.Getenv("POSTGRES_URL"),
			Port:     os.Getenv("POSTGRES_PORT"),
			Host:     os.Getenv("POSTGRES_HOST"),
			Name:     os.Getenv("DB_NAME"),
		},
		SignUp: models.SignUpConfig{
			ORG_ADD_KEY:   os.Getenv("ORG_ADD_KEY"),
			ORG_ADD_KEY_1: os.Getenv("ORG_ADD_KEY_1"),
			ORG_ADD_KEY_2: os.Getenv("ORG_ADD_KEY_2"),
			ORG_ADD_KEY_3: os.Getenv("ORG_ADD_KEY_3"),
		},
	}
	return cfg, nil
}
