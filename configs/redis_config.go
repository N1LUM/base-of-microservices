package configs

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

type RedisConfig struct {
	Host     string
	Port     string
	Password string
	DBNumber int
}

func InitRedisConfig() *RedisConfig {
	logrus.Info("Trying initialize redis configuration...")

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	dbNumber, err := strconv.Atoi(os.Getenv("REDIS_DB"))
	if err != nil {
		log.Fatalf("Error loading DBNumber")
	}

	cfg := RedisConfig{
		Host:     os.Getenv("REDIS_HOST"),
		Port:     os.Getenv("REDIS_PORT"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DBNumber: dbNumber,
	}

	logrus.Info("Successfully initialized redis configuration")

	return &cfg
}
