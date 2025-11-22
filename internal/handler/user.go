package handler

import (
	"errors"
	"net/http"
	"site-constructor/internal/dto/user_context"

	"github.com/google/uuid"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func (h *Handler) createUser(c *gin.Context) {
	var input user_context.CreateUserInput
	if err := c.ShouldBindJSON(&input); err != nil {
		logrus.Warnf("[UserHandler] Failed to bind JSON for createUser: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	createdUser, err := h.services.User.CreateUser(input)
	if err != nil {
		logrus.Errorf("[UserHandler] Failed to create user: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, createdUser)
}

func (h *Handler) getUsers(c *gin.Context) {
	users, err := h.services.User.GetAllUsers()
	if err != nil {
		logrus.Errorf("[UserHandler] Failed to get users: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, users)
}

func (h *Handler) getUserByID(c *gin.Context) {
	id := c.Param("id")

	user, err := h.services.User.GetByID(parseUUID(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logrus.Warnf("[UserHandler] User not found by ID=%s", id)
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		logrus.Errorf("[UserHandler] Failed to get user by ID=%s: %v", id, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}

func (h *Handler) updateUser(c *gin.Context) {
	id := c.Param("id")

	var input user_context.UpdateUserInput
	if err := c.ShouldBindJSON(&input); err != nil {
		logrus.Warnf("[UserHandler] Failed to bind JSON for updateUser: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{"error": "User not found"}})
		return
	}

	updatedUser, err := h.services.User.Update(parseUUID(id), input)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logrus.Warnf("[UserHandler] User not found for update by ID=%s", id)
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}
		logrus.Errorf("[UserHandler] Failed to update user by ID=%s: %v", id, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, updatedUser)
}

func (h *Handler) getUserByUsername(c *gin.Context) {
	username := c.Param("username")

	user, err := h.services.User.GetUserByUsername(username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) || user == nil {
			logrus.Warnf("[Handler] User not found by Username=%s", username)
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}
		logrus.Errorf("[Handler] Failed to get user by Username=%s: %v", username, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}

func (h *Handler) deleteUser(c *gin.Context) {
	id := c.Param("id")

	if err := h.services.User.DeleteUser(parseUUID(id)); err != nil {
		logrus.Errorf("[UserHandler] Failed to delete user by ID=%s: %v", id, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusNoContent, gin.H{})
}

func parseUUID(id string) uuid.UUID {
	uid, err := uuid.Parse(id)
	if err != nil {
		logrus.Warnf("[UserHandler] Invalid UUID: %s", id)
		return uuid.Nil
	}
	return uid
}
