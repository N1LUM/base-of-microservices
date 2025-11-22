package repository

import (
	"context"
	"site-constructor/internal/models"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type User interface {
	Create(user models.User) (*models.User, error)
	GetByID(id uuid.UUID) (*models.User, error)
	GetByUsername(username string) (*models.User, error)
	GetAll() ([]models.User, error)
	Update(user *models.User) (*models.User, error)
	Delete(id uuid.UUID) error
	DeleteRefreshTokenByUserID(userID uuid.UUID, ctx context.Context) error
}

type Auth interface {
	GetUserByUsername(username string) (*models.User, error)
	SaveUserRefreshToken(userID string, refreshToken string, ctx context.Context) error
	GetRefreshToken(userID string, ctx context.Context) (string, error)
}
type Repository struct {
	User
	Auth
}

func NewRepository(postgres *gorm.DB, redis *redis.Client) *Repository {
	return &Repository{
		Auth: NewAuthRepo(postgres, redis),
		User: NewUserRepository(postgres, redis),
	}
}
