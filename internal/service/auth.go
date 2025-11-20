package service

import (
	"context"
	"errors"
	"site-constructor/internal/auth"
	"site-constructor/internal/models"
	"site-constructor/internal/repository"

	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	repo       repository.Auth
	jwtManager *auth.JWTManager
}

func NewAuthService(repo repository.Auth, jwtManager *auth.JWTManager) *AuthService {
	logrus.Info("[AuthService] Initialized AuthService")
	return &AuthService{repo: repo, jwtManager: jwtManager}
}

func (s *AuthService) Login(username, password string, ctx context.Context) (*models.User, string, string, error) {
	user, err := s.repo.GetUserByUsername(username)
	if err != nil {
		return nil, "", "", err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		// TODO создать кастомный тип ошибок
		return nil, "", "", errors.New("invalid password")
	}

	accessToken, refreshToken, err := s.generateAccessAndRefreshTokens(user.ID.String(), ctx)
	if err != nil {
		return nil, "", "", err
	}

	return user, accessToken, refreshToken, nil
}

func (s *AuthService) Refresh(userID string, token string, ctx context.Context) (string, string, error) {
	storedToken, err := s.repo.GetRefreshToken(userID, ctx)
	if err != nil {
		return "", "", errors.New("refresh token not found or expired")
	}

	if storedToken != token {
		return "", "", errors.New("invalid refresh token")
	}

	accessToken, refreshToken, err := s.generateAccessAndRefreshTokens(userID, ctx)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func (s *AuthService) generateAccessAndRefreshTokens(userID string, ctx context.Context) (string, string, error) {
	accessToken, err := s.jwtManager.AccessToken(userID)
	if err != nil {
		return "", "", err
	}

	refreshToken, err := s.jwtManager.RefreshToken(userID)
	if err != nil {
		return "", "", err
	}

	if err := s.repo.SaveUserRefreshToken(userID, refreshToken, ctx); err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}
