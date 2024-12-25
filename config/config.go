package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	AppEnv      string
	AppPort     int
	DBHost      string
	DBPort      int
	DBUser      string
	DBPass      string
	DBName      string
	OIDCIssuer  string
	OIDCClientID string
	JWTSecret   string
}

func LoadConfig() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		log.Printf("Error loading .env file: %v", err)
		// return nil, err
	}

	appPort, err := strconv.Atoi(os.Getenv("APP_PORT"))
	if err != nil {
		log.Printf("Error converting APP_PORT to int: %v", err)
		return nil, err
	}

	dbPort, err := strconv.Atoi(os.Getenv("DB_PORT"))
	if err != nil {
		log.Printf("Error converting DB_PORT to int: %v", err)
		return nil, err
	}

	config := &Config{
		AppEnv:      os.Getenv("APP_ENV"),
		AppPort:     appPort,
		DBHost:      os.Getenv("DB_HOST"),
		DBPort:      dbPort,
		DBUser:      os.Getenv("DB_USER"),
		DBPass:      os.Getenv("DB_PASSWORD"),
		DBName:      os.Getenv("DB_NAME"),
		OIDCIssuer:  os.Getenv("OIDC_ISSUER"),
		OIDCClientID: os.Getenv("OIDC_CLIENT_ID"),
		JWTSecret:   os.Getenv("JWT_SECRET"),
	}

	return config, nil
}
