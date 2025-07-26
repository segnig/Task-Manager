package Intrastructures

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func GetFromEnv(key string) (value string) {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatal("Error Loading .env file")
	}
	value = os.Getenv(key)
	return
}
