package configs

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

type PostgresConfig struct {
	Host     string
	Port     string
	Username string
	Password string
	DBName   string
	SSLMode  string
}

func InitPostgresConfig() *PostgresConfig {
	logrus.Info("Trying initialize postgres configuration...")

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	cfg := PostgresConfig{
		Host:     os.Getenv("POSTGRES_HOST"),
		Port:     os.Getenv("POSTGRES_PORT"),
		Username: os.Getenv("POSTGRES_USER"),
		Password: os.Getenv("POSTGRES_PASSWORD"),
		DBName:   os.Getenv("POSTGRES_DB"),
		SSLMode:  os.Getenv("POSTGRES_SSLMODE"),
	}

	logrus.Info("Successfully initialized postgres configuration")

	return &cfg
}
