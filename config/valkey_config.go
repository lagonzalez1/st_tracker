package config

import (
	"crypto/tls"
	"fmt"
	"os"
	"tracker/app/models"

	"github.com/valkey-io/valkey-go"
)

func LoadValKey(ValKey models.ValKey) (valkey.Client, error) {
	host := ValKey.Host
	port := ValKey.Port
	if os.Getenv("APP_ENV") == "prod" {
		tls := &tls.Config{
			MinVersion: tls.VersionTLS12,
			ServerName: ValKey.Host,
		}
		client, err := valkey.NewClient(valkey.ClientOption{InitAddress: []string{host + ":" + port}, Username: "default", Password: "on ~* &* +@all", TLSConfig: tls})
		if err != nil {
			return nil, fmt.Errorf("unable to connect to valkey client %v", err)
		}
		fmt.Printf("Connected to cache valkey")
		return client, nil
	} else {
		client, err := valkey.NewClient(valkey.ClientOption{InitAddress: []string{host + ":" + port}, Password: "s3cr3t"})
		if err != nil {
			return nil, fmt.Errorf("unable to connect to valkey client %v", err)
		}
		fmt.Printf("Connected to cache valkey")
		return client, nil
	}
}
