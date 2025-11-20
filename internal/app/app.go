package app

import (
	"errors"
	"log"
	"net/http"
	"os"
	"site-constructor/configs"
	"site-constructor/internal/auth"
	"site-constructor/internal/database"
	"site-constructor/internal/handler"
	"site-constructor/internal/repository"
	"site-constructor/internal/service"

	"github.com/sirupsen/logrus"
)

func Run() {
	logrus.SetFormatter(new(logrus.JSONFormatter))

	redisConfig := configs.InitRedisConfig()

	redis, err := database.ConnectRedis(redisConfig)
	if err != nil {
		log.Fatalf("Failed to connect to redis: %v", err)
	}

	postgresConfig := configs.InitPostgresConfig()

	postgres, err := database.ConnectPostgres(postgresConfig)
	if err != nil {
		log.Fatalf("Failed to connect to postgres: %v", err)
	}

	database.MigratePostgres(postgres)

	jwtManager := auth.NewJWTManager(
		[]byte(os.Getenv("JWT_ACCESS_SECRET")),
		[]byte(os.Getenv("JWT_REFRESH_SECRET")),
	)

	repositories := repository.NewRepository(postgres, redis)
	services := service.NewService(repositories, jwtManager)
	handlers := handler.NewHandler(services, jwtManager)

	srv := new(configs.Server)
	address := configs.BuildAppAddress()
	logrus.Infof("🚀 Starting server on %s", address)

	if err := srv.Run(address, handlers.InitRoutes()); err != nil && !errors.Is(err, http.ErrServerClosed) {
		logrus.Fatalf("❌ Error running HTTP server: %s", err)
	}
}
