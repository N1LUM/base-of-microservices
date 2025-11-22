package repository

import (
	"site-constructor/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRepository struct {
	postgres *gorm.DB
}

func NewUserRepository(postgres *gorm.DB) *UserRepository {
	return &UserRepository{postgres: postgres}
}

func (r *UserRepository) Create(user models.User) (*models.User, error) {
	if err := r.postgres.Create(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) GetAll() ([]models.User, error) {
	var users []models.User
	if err := r.postgres.Order("created_at asc").Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func (r *UserRepository) GetByID(id uuid.UUID) (*models.User, error) {
	var user models.User
	if err := r.postgres.First(&user, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) GetByUsername(username string) (*models.User, error) {
	var user models.User
	if err := r.postgres.Where("username = ?", username).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) Update(user *models.User) (*models.User, error) {
	if err := r.postgres.Save(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (r *UserRepository) Delete(id uuid.UUID) error {
	if err := r.postgres.Delete(&models.User{}, "id = ?", id).Error; err != nil {
		return err
	}
	return nil
}
