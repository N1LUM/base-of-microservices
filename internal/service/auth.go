package service

import (
	"context"
	"errors"
	"fmt"
	"site-constructor/internal/apperrors"
	"site-constructor/internal/auth"
	"site-constructor/internal/models"
	"site-constructor/internal/repository"

	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthService struct {
	repo       repository.Auth
	jwtManager *auth.JWTManager
}

func NewAuthService(repo repository.Auth, jwtManager *auth.JWTManager) *AuthService {
	return &AuthService{repo: repo, jwtManager: jwtManager}
}

func (s *AuthService) Login(username, password string, ctx context.Context) (*models.User, string, string, error) {
	user, err := s.repo.GetUserByUsername(username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, "", "", apperrors.ErrUserNotFound
		}
		return nil, "", "", fmt.Errorf("failed to login user: %w", err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, "", "", apperrors.ErrInvalidCredentials
	}

	accessToken, refreshToken, err := s.generateAccessAndRefreshTokens(user.ID.String(), ctx)
	if err != nil {
		return nil, "", "", fmt.Errorf("failed to generate access and refresh tokens: %w", err)
	}

	logrus.Infof("[AuthService] User authenticated successfully: ID=%s, Username=%s, Name=%s", user.ID, user.Username, user.Name)
	return user, accessToken, refreshToken, nil
}

func (s *AuthService) Refresh(userID string, token string, ctx context.Context) (string, string, error) {
	storedToken, err := s.repo.GetRefreshToken(userID, ctx)
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return "", "", apperrors.ErrNotFoundRefreshToken
		}
		return "", "", fmt.Errorf("failed to get refresh token: %w", err)
	}

	if storedToken != token {
		return "", "", apperrors.ErrInvalidRefreshToken
	}

	accessToken, refreshToken, err := s.generateAccessAndRefreshTokens(userID, ctx)
	if err != nil {
		return "", "", fmt.Errorf("failed to generate access and refresh tokens: %w", err)
	}

	logrus.Infof("[AuthService] Refresh and Access token generated successfully for user: ID=%s", userID)
	return accessToken, refreshToken, nil
}

func (s *AuthService) generateAccessAndRefreshTokens(userID string, ctx context.Context) (string, string, error) {
	accessToken, err := s.jwtManager.AccessToken(userID)
	if err != nil {
		return "", "", fmt.Errorf("failed to generate access token: %w", err)
	}

	refreshToken, err := s.jwtManager.RefreshToken(userID)
	if err != nil {
		return "", "", fmt.Errorf("failed to generate refresh token: %w", err)
	}

	if err := s.repo.SaveUserRefreshToken(userID, refreshToken, ctx); err != nil {
		return "", "", fmt.Errorf("failed to save refresh token: %w", err)
	}

	return accessToken, refreshToken, nil
}
