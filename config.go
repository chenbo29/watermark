package main

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	MaxFileNum          int64
	VolcServiceId       string
	VolcAccessKeyId     string
	VolcSecretAccessKey string
}

func LoadConfig() *Config {
	_ = godotenv.Load()
	config := &Config{
		MaxFileNum:          50,
		VolcServiceId:       os.Getenv("VOLC_SERVICE_ID"),
		VolcAccessKeyId:     os.Getenv("VOLC_ACCESS_KEY_ID"),
		VolcSecretAccessKey: os.Getenv("VOLC_SECRET_ACCESS_KEY"),
	}
	if os.Getenv("max_file_num") != "" {
		config.MaxFileNum, _ = strconv.ParseInt(os.Getenv("max_file_num"), 10, 64)
	}
	return config
}
