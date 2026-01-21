package config

import (
	"context"
	"crypto/tls"
	"fmt"
	"os"
	"time"
	"tracker/app/models"

	"github.com/valkey-io/valkey-go"
)

func LoadValKey(v *models.Config) (valkey.Client, error) {
	host := v.ValKey.Host
	port := v.ValKey.Port
	address := host + ":" + port

	// Log connection details
	fmt.Printf("[VALKEY] Attempting connection to: %s\n", address)
	fmt.Printf("[VALKEY] APP_ENV: %s\n", os.Getenv("APP_ENV"))
	fmt.Printf("[VALKEY] Config: Host=%s, Port=%s\n", host, port)

	if os.Getenv("APP_ENV") != "dev" {
		fmt.Printf("[VALKEY] Using PRODUCTION mode with TLS\n")

		tls := &tls.Config{
			MinVersion: tls.VersionTLS12,
			ServerName: host,
		}

		fmt.Printf("[VALKEY] TLS Config: ServerName=%s\n", host)

		client, err := valkey.NewClient(valkey.ClientOption{
			InitAddress: []string{address},
			Username:    "default",
			Password:    "on ~* &* +@all",
			TLSConfig:   tls,
		})

		if err != nil {
			fmt.Printf("[VALKEY ERROR] Failed to connect: %v\n", err)
			fmt.Printf("[VALKEY DEBUG] Connection details: Host=%s, TLS=%v\n", address, tls != nil)
			return nil, fmt.Errorf("unable to connect to valkey client: %v", err)
		}

		// Test connection
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := client.Do(ctx, client.B().Ping().Build()).Error(); err != nil {
			fmt.Printf("[VALKEY ERROR] Ping failed: %v\n", err)
			client.Close()
			return nil, fmt.Errorf("valkey ping failed: %v", err)
		}

		fmt.Printf("[VALKEY SUCCESS] Connected to Valkey cache at %s\n", address)
		return client, nil

	} else {
		fmt.Printf("[VALKEY] Using DEVELOPMENT mode (no TLS)\n")

		client, err := valkey.NewClient(valkey.ClientOption{
			InitAddress: []string{address},
			Password:    "s3cr3t",
		})

		if err != nil {
			fmt.Printf("[VALKEY ERROR] Failed to connect: %v\n", err)
			return nil, fmt.Errorf("unable to connect to valkey client: %v", err)
		}

		// Test connection
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := client.Do(ctx, client.B().Ping().Build()).Error(); err != nil {
			fmt.Printf("[VALKEY ERROR] Ping failed: %v\n", err)
			client.Close()
			return nil, fmt.Errorf("valkey ping failed: %v", err)
		}

		fmt.Printf("[VALKEY SUCCESS] Connected to Valkey cache at %s (dev mode)\n", address)
		return client, nil
	}
}
