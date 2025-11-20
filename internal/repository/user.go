package repository

import (
	"site-constructor/internal/models"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type UserRepository struct {
	postgres *gorm.DB
}

func NewUserRepository(postgres *gorm.DB) *UserRepository {
	logrus.Info("[UserRepository] Initialized UserRepository")
	return &UserRepository{postgres: postgres}
}

func (r *UserRepository) Create(user models.User) (*models.User, error) {
	logrus.Infof("[UserRepository] Creating user_context: Username=%s", user.Username)
	if err := r.postgres.Create(&user).Error; err != nil {
		logrus.Errorf("[UserRepository] Failed to create user_context '%s': %v", user.Username, err)
		return nil, err
	}
	logrus.Infof("[UserRepository] User created successfully: ID=%s", user.ID)
	return &user, nil
}

func (r *UserRepository) GetAll() ([]models.User, error) {
	logrus.Info("[UserRepository] Fetching all users")
	var users []models.User
	if err := r.postgres.Order("created_at asc").Find(&users).Error; err != nil {
		logrus.Errorf("[UserRepository] Failed to fetch users: %v", err)
		return nil, err
	}
	logrus.Infof("[UserRepository] Retrieved %d users", len(users))
	return users, nil
}

func (r *UserRepository) GetByID(id uuid.UUID) (*models.User, error) {
	logrus.Infof("[UserRepository] Fetching user_context by ID=%s", id)
	var user models.User
	if err := r.postgres.First(&user, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			logrus.Warnf("[UserRepository] User not found: ID=%s", id)
			return nil, nil
		}
		logrus.Errorf("[UserRepository] Error fetching user_context by ID=%s: %v", id, err)
		return nil, err
	}
	logrus.Infof("[UserRepository] User found: ID=%s, Username=%s", user.ID, user.Username)
	return &user, nil
}

func (r *UserRepository) GetByUsername(username string) (*models.User, error) {
	logrus.Infof("[UserRepository] Fetching user_context by Username=%s", username)
	var user models.User
	if err := r.postgres.Where("username = ?", username).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			logrus.Warnf("[UserRepository] User not found: Username=%s", username)
			return nil, nil
		}
		logrus.Errorf("[UserRepository] Error fetching user_context by Username=%s: %v", username, err)
		return nil, err
	}
	logrus.Infof("[UserRepository] User found: ID=%s, Username=%s", user.ID, user.Username)
	return &user, nil
}

func (r *UserRepository) Update(user *models.User) (*models.User, error) {
	logrus.Infof("[UserRepository] Updating user_context: ID=%s", user.ID)
	if err := r.postgres.Save(user).Error; err != nil {
		logrus.Errorf("[UserRepository] Failed to update user_context ID=%s: %v", user.ID, err)
		return nil, err
	}
	logrus.Infof("[UserRepository] User updated successfully: ID=%s, Username=%s", user.ID, user.Username)
	return user, nil
}

func (r *UserRepository) Delete(id uuid.UUID) error {
	logrus.Infof("[UserRepository] Deleting user_context ID=%s", id)
	if err := r.postgres.Delete(&models.User{}, "id = ?", id).Error; err != nil {
		logrus.Errorf("[UserRepository] Failed to delete user_context ID=%s: %v", id, err)
		return err
	}
	logrus.Infof("[UserRepository] User deleted successfully: ID=%s", id)
	return nil
}
