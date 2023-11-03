package config

import (
	"fmt"
	"log"
	"os"
)

func SanityCheck() {
	envProps := []string{
		"POSTGRES_USER",
		"POSTGRES_PASSWORD",
		"POSTGRES_HOST",
		"POSTGRES_PORT",
		"POSTGRES_DB",
		"CLIENT_ORIGIN",
		"PORT",
		"ACCESS_SECRET_KEY",
		"REFRESH_SECRET_KEY",
		"SENDER_MAIL",
		"SENDER_PASS",
		"REDIS_HOST",
		"REDIS_PORT",
		"REDIS_DB",
	}
	for _, k := range envProps {
		if os.Getenv(k) == "" {
			log.Fatal(fmt.Sprintf("Environment variable %s not defined. Terminating application...", k))
		}
	}
}
