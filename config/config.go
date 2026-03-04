package config

import (
	"fmt"
	"log"
	"os"
	"tracker/app/models"

	"github.com/joho/godotenv"
)

func LoadConfig() (*models.Config, error) {
	appEnv := os.Getenv("APP_ENV")

	// 2. Map shorthand to full names using a simple map
	envMapping := map[string]string{
		"dev":  "development",
		"prod": "production",
	}

	if fullEnv, ok := envMapping[appEnv]; ok {
		appEnv = fullEnv
	} else if appEnv == "" {
		appEnv = "production" // Fallback
	}

	// 3. Construct filename
	envFile := fmt.Sprintf(".env.%s", appEnv)

	// 4. Try to load.
	// Note: It's common in Prod to NOT use a file and use K8s env vars directly.
	err := godotenv.Load(envFile)
	if err != nil {
		log.Printf("Note: %s not found. Using system environment variables.", envFile)
	} else {
		log.Printf("Successfully loaded configuration from %s", envFile)
	}

	// Set Stripe API key

	cfg := &models.Config{
		JWT:           os.Getenv("JWT"),
		Port:          os.Getenv("B_PORT"),
		StripeWebhook: os.Getenv("STRIPE_WEBHOOK"),
		DB: models.PostGresConfig{
			Username: os.Getenv("POSTGRES_USER"),
			Password: os.Getenv("POSTGRES_PASSWORD"),
			URL:      os.Getenv("POSTGRES_URL"),
			Port:     os.Getenv("POSTGRES_PORT"),
			Host:     os.Getenv("POSTGRES_HOST"),
			Name:     os.Getenv("DB_NAME"),
			Env:      os.Getenv("APP_ENV"),
		},
		SignUp: models.SignUpConfig{
			ORG_ADD_KEY:   os.Getenv("ORG_ADD_KEY"),
			ORG_ADD_KEY_1: os.Getenv("ORG_ADD_KEY_1"),
			ORG_ADD_KEY_2: os.Getenv("ORG_ADD_KEY_2"),
			ORG_ADD_KEY_3: os.Getenv("ORG_ADD_KEY_3"),
		},
		MQ: models.MQ{
			Host:     os.Getenv("RABBIT_HOST"),
			Username: os.Getenv("RABBIT_USERNAME"),
			Password: os.Getenv("RABBIT_PASSWORD"),
			Port:     os.Getenv("RABBIT_PORT"),
			AmazonMQ: os.Getenv("AMAZON_MQ"),
		},
		S3: models.S3Client{
			AccessKey: os.Getenv("AWS_ACCESS_KEY_ID"),
			SecretKey: os.Getenv("AWS_SECRET_ACCESS_KEY"),
			Region:    os.Getenv("AWS_REGION"),
		},
		ValKey: models.ValKey{
			Port: os.Getenv("VALKEY_PORT"),
			Host: os.Getenv("VALKEY_HOST"),
		},
		SQS: models.SQSConfig{
			DataReportsQueue:      os.Getenv("PROD_QUEUE_NAME_DATA_REPORTS"),
			AssessmentGraderQueue: os.Getenv("PROD_QUEUE_NAME_DATA_ASSESSMENTS"),
			GenerateContentQueue:  os.Getenv("PROD_QUEUE_NAME_GENERATE_CONTENT"),
		},
	}
	return cfg, nil
}
