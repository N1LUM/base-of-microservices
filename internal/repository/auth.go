package repository

import (
	"context"
	"site-constructor/internal/auth"
	"site-constructor/internal/models"

	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type AuthRepository struct {
	postgres *gorm.DB
	redis    *redis.Client
}

func NewAuthRepo(postgres *gorm.DB, redis *redis.Client) *AuthRepository {
	logrus.Info("[AuthRepository] Initialized AuthRepository")
	return &AuthRepository{postgres: postgres, redis: redis}
}

func (r *AuthRepository) GetUserByUsername(username string) (*models.User, error) {
	var user models.User
	if err := r.postgres.First(&user, "username = ?", username).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *AuthRepository) SaveUserRefreshToken(userID string, refreshToken string, ctx context.Context) error {
	return r.redis.Set(ctx, "refresh:"+userID, refreshToken, auth.RefreshTokenExpiration).Err()
}
func (r *AuthRepository) GetRefreshToken(userID string, ctx context.Context) (string, error) {
	return r.redis.Get(ctx, "refresh:"+userID).Result()
}
