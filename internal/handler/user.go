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
		logrus.Errorf("[UserHandler] Failed to create user_context: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	logrus.Infof("[UserHandler] User created successfully: ID=%s, Username=%s", createdUser.ID, createdUser.Username)
	c.JSON(http.StatusCreated, createdUser)
}

func (h *Handler) getUsers(c *gin.Context) {
	users, err := h.services.User.GetAllUsers()
	if err != nil {
		logrus.Errorf("[UserHandler] Failed to get users: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	logrus.Infof("[UserHandler] Retrieved %d users", len(users))
	c.JSON(http.StatusOK, users)
}

func (h *Handler) getUserByID(c *gin.Context) {
	id := c.Param("id")

	user, err := h.services.User.GetUserByID(parseUUID(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logrus.Warnf("[UserHandler] User not found: ID=%s", id)
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}
		logrus.Errorf("[UserHandler] Failed to get user_context by ID=%s: %v", id, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	logrus.Infof("[UserHandler] User fetched: ID=%s, Username=%s", user.ID, user.Username)
	c.JSON(http.StatusOK, user)
}

func (h *Handler) updateUser(c *gin.Context) {
	id := c.Param("id")

	var input user_context.UpdateUserInput
	if err := c.ShouldBindJSON(&input); err != nil {
		logrus.Warnf("[UserHandler] Failed to bind JSON for updateUser: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updatedUser, err := h.services.User.UpdateUser(parseUUID(id), input)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logrus.Warnf("[UserHandler] User not found for update: ID=%s", id)
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}
		logrus.Errorf("[UserHandler] Failed to update user_context ID=%s: %v", id, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	logrus.Infof("[UserHandler] User updated successfully: ID=%s, Username=%s", updatedUser.ID, updatedUser.Username)
	c.JSON(http.StatusOK, updatedUser)
}

func (h *Handler) loginUser(c *gin.Context) {
	var input struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.services.User.LoginUser(input.Username, input.Password)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "user_context not found"})
			return
		}
		if err.Error() == "invalid password" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid password"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":       user.ID,
		"username": user.Username,
	})
}
func (h *Handler) getUserByUsername(c *gin.Context) {
	username := c.Param("username")

	logrus.Infof("[Handler] Fetching user_context by username: %s", username)

	user, err := h.services.User.GetUserByUsername(username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) || user == nil {
			logrus.Warnf("[Handler] User not found: username=%s", username)
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}
		logrus.Errorf("[Handler] Failed to get user_context by username=%s: %v", username, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	logrus.Infof("[Handler] User fetched successfully: username=%s, id=%s", username, user.ID)
	c.JSON(http.StatusOK, user)
}

func (h *Handler) deleteUser(c *gin.Context) {
	id := c.Param("id")

	if err := h.services.User.DeleteUser(parseUUID(id)); err != nil {
		logrus.Errorf("[UserHandler] Failed to delete user_context ID=%s: %v", id, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	logrus.Infof("[UserHandler] User deleted successfully: ID=%s", id)
	c.JSON(http.StatusNoContent, gin.H{})
}

// parseUUID - вспомогательная функция для конвертации строки в uuid.UUID
func parseUUID(id string) uuid.UUID {
	uid, err := uuid.Parse(id)
	if err != nil {
		logrus.Warnf("[UserHandler] Invalid UUID: %s", id)
		return uuid.Nil
	}
	return uid
}
