package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const (
	AccessTokenExpiration  = 24 * time.Hour
	RefreshTokenExpiration = 30 * 24 * time.Hour
)

type JWTManager struct {
	AccessSecret  []byte
	RefreshSecret []byte
	signingMethod jwt.SigningMethod
}

func NewJWTManager(accessSecret, refreshSecret []byte) *JWTManager {
	return &JWTManager{
		AccessSecret:  accessSecret,
		RefreshSecret: refreshSecret,
		signingMethod: jwt.SigningMethodHS256,
	}
}

func (j *JWTManager) AccessToken(userID string) (string, error) {
	return j.generateToken(userID, AccessTokenExpiration, j.AccessSecret)
}

func (j *JWTManager) RefreshToken(userID string) (string, error) {
	return j.generateToken(userID, RefreshTokenExpiration, j.RefreshSecret)
}

func (j *JWTManager) ParseAccessToken(token string) (*jwt.RegisteredClaims, error) {
	return j.parseToken(token, j.AccessSecret)
}

func (j *JWTManager) ParseRefreshToken(token string) (*jwt.RegisteredClaims, error) {
	return j.parseToken(token, j.RefreshSecret)
}

func (j *JWTManager) parseToken(token string, jwtSecret []byte) (*jwt.RegisteredClaims, error) {
	parsedToken, err := jwt.ParseWithClaims(token, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	if parsedToken.Method != j.signingMethod {
		return nil, errors.New("invalid signing method")
	}

	if claims, ok := parsedToken.Claims.(*jwt.RegisteredClaims); ok && parsedToken.Valid {
		return claims, nil
	}

	return nil, err
}

func (j *JWTManager) generateToken(userID string, expiration time.Duration, secret []byte) (string, error) {
	token := jwt.NewWithClaims(j.signingMethod, jwt.RegisteredClaims{
		Subject:   userID,
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiration)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	})

	return token.SignedString(secret)
}
