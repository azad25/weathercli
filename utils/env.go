package utils

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// load env file
func LoadEnv() {
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

// get key from env
func GetEnv(key string) string {
	return os.Getenv(key)
}
