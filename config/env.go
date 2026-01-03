package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

func Env(key string, defaultValue string) string {
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Printf("Warning: Error loading .env file: %v\n", err)
	}

	value := os.Getenv(key)

	if value == "" {
		return defaultValue
	}

	return value
}
