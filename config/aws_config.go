package config

// Use this code snippet in your app.
// If you need more information about configurations or implementing the sample code, visit the AWS docs:
// https://aws.github.io/aws-sdk-go-v2/docs/getting-started/

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"reflect"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
)

type Secrets struct {
	B_PORT            string `env:"DB_HOST"`
	POSTGRES_HOST     string `env:"POSTGRES_HOST"`
	POSTGRES_USER     string `env:"POSTGRES_USER"`
	POSTGRES_PORT     string `env:"POSTGRES_PORT"`
	POSTGRES_PASSWORD string `env:"POSTGRES_PASSWORD"`
	POSTGRES_URL      string `env:"POSTGRES_URL"`
	DB_NAME           string `env:"DB_NAME"`
	JWT               string `env:"JWT"`
	ORG_ADD_KEY       string `env:"ORG_ADD_KEY"`
	ORG_ADD_KEY_1     string `env:"ORG_ADD_KEY_1"`
	ORG_ADD_KEY_2     string `env:"ORG_ADD_KEY_2"`
	ORG_ADD_KEY_3     string `env:"ORG_ADD_KEY_3"`
}

func awsConfig() error {
	secretName := "tracker-backend-keys"
	region := "us-west-1"

	config, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(region))
	if err != nil {
		return err
	}

	// Create Secrets Manager client
	svc := secretsmanager.NewFromConfig(config)

	input := &secretsmanager.GetSecretValueInput{
		SecretId:     aws.String(secretName),
		VersionStage: aws.String("AWSCURRENT"), // VersionStage defaults to AWSCURRENT if unspecified
	}

	result, err := svc.GetSecretValue(context.TODO(), input)
	if err != nil {
		// For a list of exceptions thrown, see
		// https://docs.aws.amazon.com/secretsmanager/latest/apireference/API_GetSecretValue.html
		return err
	}

	// Decrypts secret using the associated KMS key.
	var secretString string = *result.SecretString
	secrets := &Secrets{}
	if err := json.Unmarshal([]byte(secretString), &secrets); err != nil {
		return fmt.Errorf("Error decoding JSON: %v\n", err)
	}

	err = writeEnvFile(".env.production", secrets)
	if err != nil {
		return fmt.Errorf("unable to write to file")
	}
	fmt.Println("Successfully created config.env")
	// Your code goes here.
	return nil
}

func writeEnvFile(filename string, config interface{}) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	val := reflect.ValueOf(config)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	fmt.Println("SIZE: ", val.NumField())
	for i := 0; i < val.NumField(); i++ {
		field := val.Type().Field(i)
		envTag := field.Tag.Get("env")
		if envTag == "" {
			continue
		}
		value := val.Field(i)
		strValue := value.String()
		println(strValue)
		_, err := fmt.Fprintf(file, "%s=%s\n", envTag, strValue)
		if err != nil {
			fmt.Print("unable to write %v", err)
			return err
		}
	}
	return nil
}
