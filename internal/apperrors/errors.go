package apperrors

import "errors"

//TODO после роста проекта вынести контексты в разные директории

var (
	// Authentication context

	ErrNotFoundRefreshToken = errors.New("refresh token not found or expired")
	ErrInvalidCredentials   = errors.New("invalid credentials")
	ErrInvalidRefreshToken  = errors.New("invalid refresh token")

	// User context

	ErrUserExists   = errors.New("user already exists")
	ErrUserNotFound = errors.New("user not found")
)
