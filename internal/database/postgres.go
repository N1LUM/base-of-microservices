package database

import (
	"fmt"
	"site-constructor/configs"
	"site-constructor/internal/models"

	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func ConnectPostgres(cfg *configs.PostgresConfig) (*gorm.DB, error) {
	logrus.Info("Trying connect to postgres...")

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		cfg.Host, cfg.Username, cfg.Password, cfg.DBName, cfg.Port, cfg.SSLMode,
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return nil, fmt.Errorf("failed connect to postgres: %w", err)
	}

	logrus.Info("Successfully connected to postgres")

	return db, nil
}

func MigratePostgres(db *gorm.DB) {
	logrus.Info("Running auto-migrations...")

	if err := db.AutoMigrate(
		&models.User{},
	); err != nil {
		logrus.Fatalf("migration failed: %v", err)
	}

	logrus.Info("Auto-migrations completed successfully")
}
