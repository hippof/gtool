package genv

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func init() {
	EnvInit(".env")
}

func EnvInit(filenames ...string) {
	err := godotenv.Load(filenames...)
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func Env(keys ...string) string {
	if v := os.Getenv(keys[0]); v != "" {
		return v
	}
	if len(keys) > 1 {
		return keys[1]
	}
	return ""
}
