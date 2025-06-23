package config

import (
	"log"
	"os"
)

var (
	_config = config{}
)

func init() {
	_config = config{
		Port:               getEnv("PORT"),
		ApiPort:            getEnv("API_PORT"),
		AdminPort:          getEnv("ADMIN_PORT"),
		ApiAddress:         getEnv("API_ADDRESS"),
		GeniusClientId:     getEnv("GENIUS_CLIENT_ID"),
		GeniusClientSecret: getEnv("GENIUS_CLIENT_SECRET"),
		JwtSecret:          getEnv("JWT_SECRET"),
		DB: struct {
			Name     string
			Host     string
			Username string
			Password string
		}{
			Name:     getEnv("DB_NAME"),
			Host:     getEnv("DB_HOST"),
			Username: getEnv("DB_USERNAME"),
			Password: getEnv("DB_PASSWORD"),
		},
		Smtp: struct {
			Host     string
			Port     string
			Username string
			Password string
		}{
			Host:     getEnv("SMTP_HOST"),
			Port:     getEnv("SMTP_PORT"),
			Username: getEnv("SMTP_USER"),
			Password: getEnv("SMTP_PASSWORD"),
		},
	}
}

type config struct {
	Port               string
	ApiPort            string
	AdminPort          string
	ApiAddress         string
	GeniusClientId     string
	GeniusClientSecret string
	JwtSecret          string
	DB                 struct {
		Name     string
		Host     string
		Username string
		Password string
	}
	Smtp struct {
		Host     string
		Port     string
		Username string
		Password string
	}
}

// Env returns the thing's config values :)
func Env() config {
	return _config
}

func getEnv(key string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		log.Fatalln("The \"" + key + "\" variable is missing.")
	}
	return value
}
