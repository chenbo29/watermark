package main

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	VolcServiceId       string
	VolcAccessKeyId     string
	VolcSecretAccessKey string
}

func LoadConfig() *Config {
	_ = godotenv.Load()

	return &Config{
		VolcServiceId:       os.Getenv("VOLC_SERVICE_ID"),
		VolcAccessKeyId:     os.Getenv("VOLC_ACCESS_KEY_ID"),
		VolcSecretAccessKey: os.Getenv("VOLC_SECRET_ACCESS_KEY"),
	}
}
