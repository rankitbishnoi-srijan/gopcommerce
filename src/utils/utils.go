package utils

import (
	"github.com/joho/godotenv"
)

// loadEnv loads environment variables from a .env file
func LoadEnv() error {
	err := godotenv.Load()
	if err != nil {
		return err
	}

	return nil
}
