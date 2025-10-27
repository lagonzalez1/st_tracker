package config

import (
	"crypto/tls"
	"fmt"
	"os"

	"github.com/valkey-io/valkey-go"
)

func LoadValKey() (valkey.Client, error) {
	env, err := LoadConfig()
	if err != nil {
		return nil, err
	}
	host := env.ValKey.Host
	port := env.ValKey.Port
	var client valkey.Client
	if os.Getenv("APP_ENV") != "" {
		tls := &tls.Config{
			MinVersion: tls.VersionTLS12,
			ServerName: env.ValKey.Host,
		}
		client, err = valkey.NewClient(valkey.ClientOption{InitAddress: []string{host + ":" + port}, Username: "default", Password: "on ~* &* +@all", TLSConfig: tls})
		if err != nil {
			return nil, fmt.Errorf("unable to connect to valkey client %v", err)
		}
	} else {
		client, err = valkey.NewClient(valkey.ClientOption{InitAddress: []string{host + ":" + port}, Password: "s3cr3t"})
		if err != nil {
			return nil, fmt.Errorf("unable to connect to valkey client %v", err)
		}

	}
	return client, nil
}
