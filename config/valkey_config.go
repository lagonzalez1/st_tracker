package config

import (
	"fmt"

	"github.com/valkey-io/valkey-go"
)

func LoadValKey() (*valkey.Client, error) {
	env, err := LoadConfig()
	if err != nil {
		return nil, err
	}
	host := env.ValKey.Host
	port := env.ValKey.Port

	client, err := valkey.NewClient(valkey.ClientOption{InitAddress: []string{host + ":" + port}})
	if err != nil {
		return nil, fmt.Errorf("unable to connect to valkey client %v", err)
	}
	defer client.Close()

	return &client, nil
}
