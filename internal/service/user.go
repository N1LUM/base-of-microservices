package service

import (
	"errors"
	"site-constructor/internal/dto/user_context"
	"site-constructor/internal/models"
	"site-constructor/internal/repository"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserService struct {
	repo repository.User
}

func NewUserService(repo repository.User) *UserService {
	logrus.Info("[UserService] Initialized UserService")
	return &UserService{repo: repo}
}

func (s *UserService) CreateUser(input user_context.CreateUserInput) (*models.User, error) {
	logrus.Infof("[UserService] Creating User: username = %s", input.Username)

	existingUser, _ := s.repo.GetByUsername(input.Username)
	if existingUser != nil {
		logrus.Warnf("[UserService] User '%s' already exists", input.Username)
		return nil, errors.New("Имя пользователя уже занято")
	}

	hashedPassword, err := hashPassword(input.Password)
	if err != nil {
		logrus.Errorf("[UserService] Failed to hash password for '%s': '%v'", input.Username, err)
		return nil, err
	}

	user := models.User{
		Name:     input.Username,
		Username: input.Username,
		Password: hashedPassword,
	}

	created, err := s.repo.Create(user)
	if err != nil {
		logrus.Errorf("[UserService] Failed to create user_context '%s': '%v'", input.Username, err)
		return nil, err
	}

	logrus.Infof("[UserService] User created successfully: id=%s username=%s", created.ID, created.Username)
	return created, nil
}

func (s *UserService) LoginUser(username, password string) (*models.User, error) {
	logrus.Infof("[UserService] Login attempt for username = %s", username)

	user, err := s.repo.GetByUsername(username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logrus.Warnf("[UserService] User not found: username=%s", username)
			return nil, gorm.ErrRecordNotFound
		}
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		logrus.Warnf("[UserService] Invalid password for username=%s", username)
		return nil, errors.New("invalid password")
	}

	return user, nil
}

func (s *UserService) GetAllUsers() ([]models.User, error) {
	logrus.Infof("[UserService] Fetching all Users")
	users, err := s.repo.GetAll()
	if err != nil {
		logrus.Errorf("[UserService] Failed to get Users: '%v'", err)
		return nil, err
	}
	logrus.Infof("[UserService] Retrieved %d Users", len(users))
	return users, nil
}

func (s *UserService) GetUserByID(id uuid.UUID) (*models.User, error) {
	logrus.Infof("[UserService] Fetching User by ID: %s", id)

	user, err := s.repo.GetByID(id)
	if err != nil {
		logrus.Errorf("[UserService] User not found by ID=%s", id)
		return nil, err
	}

	logrus.Infof("[UserService] User found: username=%s", user.Username)
	return user, nil
}

func (s *UserService) GetUserByUsername(username string) (*models.User, error) {
	logrus.Infof("[UserService] Fetching User by Username: %s", username)
	user, err := s.repo.GetByUsername(username)
	if err != nil {
		logrus.Errorf("[UserService] User not found by Username=%s", username)
		return nil, err
	}

	logrus.Infof("[UserService] User found: ID=%s", user.ID)
	return user, nil
}

func (s *UserService) UpdateUser(id uuid.UUID, input user_context.UpdateUserInput) (*models.User, error) {
	logrus.Infof("[UserService] Updating user_context ID=%s", id)
	user, err := s.repo.GetByID(id)
	if err != nil {
		logrus.Warnf("[UserService] User not found for update: ID=%s", id)
		return nil, errors.New("user_context not found")
	}

	if input.Name != nil {
		logrus.Debugf("[UserService] Updating name for ID=%s", id)
		user.Name = *input.Name
	}
	if input.Username != nil {
		logrus.Debugf("[UserService] Updating username for ID=%s", id)
		user.Username = *input.Username
	}
	if input.Password != nil {
		hashed, err := hashPassword(*input.Password)
		if err != nil {
			logrus.Errorf("[UserService] Failed to hash new password: %v", err)
			return nil, err
		}
		user.Password = hashed
	}

	updated, err := s.repo.Update(user)
	if err != nil {
		logrus.Errorf("[UserService] Failed to update user_context ID=%s: %v", id, err)
		return nil, err
	}

	logrus.Infof("[UserService] User updated successfully: ID=%s", id)
	return updated, nil
}

// TODO Добавить удаление токена
func (s *UserService) DeleteUser(id uuid.UUID) error {
	logrus.Infof("[UserService] Deleting user_context ID=%s", id)
	if err := s.repo.Delete(id); err != nil {
		logrus.Errorf("[UserService] Failed to delete user_context ID=%s: %v", id, err)
		return err
	}

	logrus.Infof("[UserService] User deleted successfully: ID=%s", id)
	return nil
}

func hashPassword(password string) (string, error) {
	if password == "" {
		return "", errors.New("password is empty")
	}
	bs, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(bs), nil
}
