package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"tidebot/pkg/users/models"

	"github.com/labstack/echo/v4"
)

type UserRepository interface {
	ListAll() ([]models.User, error)
	GetByID(id int) (models.User, error)
	GetByPhoneNumber(phoneNumber string) (models.User, error)
	Save(models.UserWriteModel) (models.User, error)
	Update(id int, writeModel models.UserWriteModel) (models.User, error)
	Delete(id int) error
}

type userRepositoryImpl struct {
	db  *sql.DB
	log echo.Logger
}

func NewUserRepository(db *sql.DB, log echo.Logger) UserRepository {
	return &userRepositoryImpl{db, log}
}


func (r *userRepositoryImpl) ListAll() ([]models.User, error) {
	r.log.Debugf("Attempting to list all users")

	query := `SELECT id, phone_number, name, created_at, updated_at FROM users ORDER BY created_at DESC`

	rows, err := r.db.QueryContext(context.Background(), query)
	if err != nil {
		return []models.User{}, fmt.Errorf("failed to list all users: %w", err)
	}
	defer rows.Close()

	var allUsers []models.User

	for rows.Next() {
		var user models.User
		err := rows.Scan(
			&user.ID,
			&user.PhoneNumber,
			&user.Name,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			return []models.User{}, fmt.Errorf("failed to scan user row: %w", err)
		}
		allUsers = append(allUsers, user)
	}

	if err := rows.Err(); err != nil {
		return []models.User{}, fmt.Errorf("failed to read rows when trying to list all users: %w", err)
	}

	r.log.Debugf("Successfully listed %d users", len(allUsers))
	return allUsers, nil
}

func (r *userRepositoryImpl) GetByID(id int) (models.User, error) {
	r.log.Debugf("Attempting to get user by ID: %d", id)

	query := `SELECT id, phone_number, name, created_at, updated_at FROM users WHERE id = ? LIMIT 1`

	var user models.User
	err := r.db.QueryRowContext(context.Background(), query, id).Scan(
		&user.ID,
		&user.PhoneNumber,
		&user.Name,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return models.User{}, fmt.Errorf("user not found with id '%d'", id)
		}
		return models.User{}, fmt.Errorf("failed to get user by id: %w", err)
	}

	r.log.Debugf("Successfully found user with phone='%s' for ID '%d'", user.PhoneNumber, id)
	return user, nil
}

func (r *userRepositoryImpl) GetByPhoneNumber(phoneNumber string) (models.User, error) {
	r.log.Debugf("Attempting to get user by phone number: %s", phoneNumber)

	query := `SELECT id, phone_number, name, created_at, updated_at FROM users WHERE phone_number = ? LIMIT 1`

	var user models.User
	err := r.db.QueryRowContext(context.Background(), query, phoneNumber).Scan(
		&user.ID,
		&user.PhoneNumber,
		&user.Name,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return models.User{}, fmt.Errorf("user not found with phone number '%s'", phoneNumber)
		}
		return models.User{}, fmt.Errorf("failed to get user by phone number: %w", err)
	}

	r.log.Debugf("Successfully found user with ID='%d' for phone number '%s'", user.ID, phoneNumber)
	return user, nil
}

func (r *userRepositoryImpl) Save(writeModel models.UserWriteModel) (models.User, error) {
	r.log.Debugf("Attempting to save a new user: %+v", writeModel)

	query := `
		INSERT INTO users (phone_number, name, created_at, updated_at) 
		VALUES (?, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP) 
		RETURNING id, phone_number, name, created_at, updated_at`

	var user models.User
	err := r.db.QueryRowContext(
		context.Background(),
		query,
		writeModel.PhoneNumber,
		writeModel.Name,
	).Scan(
		&user.ID,
		&user.PhoneNumber,
		&user.Name,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return models.User{}, fmt.Errorf("failed to save new user: %w", err)
	}

	r.log.Debugf("Saved new user with id='%d'", user.ID)
	return user, nil
}

func (r *userRepositoryImpl) Update(id int, writeModel models.UserWriteModel) (models.User, error) {
	r.log.Debugf("Attempting to update user with id='%d': %+v", id, writeModel)

	query := `
		UPDATE users 
		SET phone_number = ?, name = ?, updated_at = CURRENT_TIMESTAMP 
		WHERE id = ?
		RETURNING id, phone_number, name, created_at, updated_at`

	var user models.User
	err := r.db.QueryRowContext(
		context.Background(),
		query,
		writeModel.PhoneNumber,
		writeModel.Name,
		id,
	).Scan(
		&user.ID,
		&user.PhoneNumber,
		&user.Name,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return models.User{}, fmt.Errorf("user not found with id='%d'", id)
		}
		return models.User{}, fmt.Errorf("failed to update user: %w", err)
	}

	r.log.Debugf("Successfully updated user with id='%d'", id)
	return user, nil
}

func (r *userRepositoryImpl) Delete(id int) error {
	r.log.Debugf("Attempting to delete user with id='%d'", id)

	query := `DELETE FROM users WHERE id = ?`

	result, err := r.db.ExecContext(context.Background(), query, id)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("user not found with id='%d'", id)
	}

	r.log.Debugf("Successfully deleted user with id='%d'", id)
	return nil
}