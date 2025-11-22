package handler

import (
	"site-constructor/internal/dto/auth_context"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"errors"
	"net/http"
)

func (h *Handler) Login(c *gin.Context) {
	var input auth_context.LoginInput
	ctx := c.Request.Context()

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, accessToken, refreshToken, err := h.services.Auth.Login(input.Username, input.Password, ctx)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":           user.ID,
		"username":     user.Username,
		"accessToken":  accessToken,
		"refreshToken": refreshToken,
	})
}

func (h *Handler) Refresh(c *gin.Context) {
	var input auth_context.RefreshInput
	ctx := c.Request.Context()

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	claims, err := h.jwtManager.ParseRefreshToken(input.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
		return
	}

	accessToken, refreshToken, err := h.services.Auth.Refresh(claims.Subject, input.RefreshToken, ctx)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":           claims.Subject,
		"accessToken":  accessToken,
		"refreshToken": refreshToken,
	})
}
