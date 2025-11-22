package service

import (
	"context"
	"site-constructor/internal/auth"
	"site-constructor/internal/dto/user_context"
	"site-constructor/internal/models"
	"site-constructor/internal/repository"

	"github.com/google/uuid"
)

type User interface {
	CreateUser(input user_context.CreateUserInput) (*models.User, error)
	GetAllUsers() ([]models.User, error)
	GetByID(id uuid.UUID) (*models.User, error)
	GetUserByUsername(username string) (*models.User, error)
	Update(id uuid.UUID, input user_context.UpdateUserInput) (*models.User, error)
	DeleteUser(id uuid.UUID, ctx context.Context) error
}

type Auth interface {
	Login(username, password string, ctx context.Context) (*models.User, string, string, error)
	Refresh(userID string, token string, ctx context.Context) (string, string, error)
}

type Service struct {
	User
	Auth
}

func NewService(repos *repository.Repository, jwtManager *auth.JWTManager) *Service {
	return &Service{
		User: NewUserService(repos.User),
		Auth: NewAuthService(repos.Auth, jwtManager),
	}
}
