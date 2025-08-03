package repository

import (
	"go-library-service/cmd/api/entity"
	errmap "go-library-service/internal/error_map"

	"github.com/pkg/errors"
	"gorm.io/gorm"
)

// CreateUser creates a new user
func (r *PostgresRepository) CreateUser(user entity.User) (*uint, error) {
	err := r.postgres.Table("users").Create(&user).Error
	if err != nil {
		return nil, errors.Wrap(err, "[PostgresRepository.CreateUser]: unable to create user")
	}
	return &user.ID, nil
}

// GetUserByID retrieves a user by their ID
func (r *PostgresRepository) GetUserByID(userID uint) (*entity.UserResponse, error) {
	var user entity.UserResponse
	err := r.postgres.Table("users").First(&user, userID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errmap.ErrmapNotFound
		}
		
		return nil, errors.Wrap(err, "[PostgresRepository.GetUserByID]: unable to get user")
	}
	return &user, nil
}

// GetUserByUsername retrieves a user by username
func (r *PostgresRepository) GetUserByUsername(username string) (*entity.User, error) {
	var user entity.User
	err := r.postgres.Table("users").Where("username = ?", username).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errmap.ErrmapNotFound
		}
		return nil, errors.Wrapf(err, "[PostgresRepository.GetUserByUsername]: unable to find user with username %s", username)
	}

	return &user, nil
}

// UpdateUser updates an existing user
func (r *PostgresRepository) UpdateUser(user entity.User) error {
    err := r.postgres.Table("users").Model(&user).Update("name", user.Name).Error
    if err != nil {
        return errors.Wrap(err, "[PostgresRepository.UpdateUser]: unable to update user")
    }
	return nil
}

// DeleteUser deletes a user by ID
func (r *PostgresRepository) DeleteUser(userID uint) error {
	err := r.postgres.Table("users").Delete(&entity.User{}, userID).Error
	if err != nil {
		return errors.Wrap(err, "[PostgresRepository.DeleteUser]: unable to delete user")
	}
	return nil
}
