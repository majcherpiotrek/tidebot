package services

import (
	"database/sql"
	"fmt"
	"strings"
	"tidebot/pkg/users/models"
	"tidebot/pkg/users/repositories"

	"github.com/labstack/echo/v4"
)

type UserService interface {
	SaveUser(phoneNumber string, name *string) (models.User, error)
	GetAllUsers() ([]models.User, error)
	GetUserByID(id int) (models.User, error)
	GetUserByPhoneNumber(phoneNumber string) (models.User, error)
}

type userServiceImpl struct {
	userRepository repositories.UserRepository
	db             *sql.DB
	log            echo.Logger
}

func NewUserService(userRepository repositories.UserRepository, db *sql.DB, log echo.Logger) UserService {
	return &userServiceImpl{
		userRepository: userRepository,
		db:             db,
		log:            log,
	}
}

func (s *userServiceImpl) SaveUser(phoneNumber string, name *string) (models.User, error) {
	s.log.Debugf("Attempting to save user - phone: %s, name: %v", phoneNumber, name)

	tx, err := s.db.Begin()
	if err != nil {
		return models.User{}, fmt.Errorf("failed to start transaction: %w", err)
	}
	defer tx.Rollback()

	existingUser, err := s.userRepository.GetByPhoneNumber(phoneNumber)
	if err == nil {
		s.log.Infof("User already exists with phone number '%s' (ID: %d), skipping save", phoneNumber, existingUser.ID)
		tx.Commit()
		return existingUser, nil
	}

	if !strings.Contains(err.Error(), "user not found") {
		return models.User{}, fmt.Errorf("failed to check if user exists: %w", err)
	}

	writeModel := models.UserWriteModel{
		PhoneNumber: phoneNumber,
		Name:        name,
	}

	savedUser, err := s.userRepository.Save(writeModel)
	if err != nil {
		return models.User{}, fmt.Errorf("failed to save new user: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return models.User{}, fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.log.Infof("Successfully saved new user - ID: %d, phone: %s", savedUser.ID, savedUser.PhoneNumber)
	return savedUser, nil
}

func (s *userServiceImpl) GetAllUsers() ([]models.User, error) {
	s.log.Debugf("Getting all users")

	users, err := s.userRepository.ListAll()
	if err != nil {
		return nil, fmt.Errorf("failed to get all users: %w", err)
	}

	s.log.Debugf("Successfully retrieved %d users", len(users))
	return users, nil
}

func (s *userServiceImpl) GetUserByID(id int) (models.User, error) {
	s.log.Debugf("Getting user by ID: %d", id)

	user, err := s.userRepository.GetByID(id)
	if err != nil {
		return models.User{}, fmt.Errorf("failed to get user by ID=%d: %w", id, err)
	}

	s.log.Debugf("Successfully retrieved user with id %d and phone number %s", user.ID, user.PhoneNumber)
	return user, nil
}

func (s *userServiceImpl) GetUserByPhoneNumber(phoneNumber string) (models.User, error) {
	s.log.Debugf("Getting user by phone number: %s", phoneNumber)

	user, err := s.userRepository.GetByPhoneNumber(phoneNumber)
	if err != nil {
		return models.User{}, fmt.Errorf("Failed to get user by phone number %s: %w", phoneNumber, err)
	}

	s.log.Debugf("Successfully retrieved user with id %d and phone number %s", user.ID, user.PhoneNumber)
	return user, nil
}
