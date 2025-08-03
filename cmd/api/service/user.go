package service

import (
	"go-library-service/cmd/api/constant"
	"go-library-service/cmd/api/entity"
	errmap "go-library-service/internal/error_map"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

// CreateUser creates a new user
func (s *Service) CreateUser(req entity.UserCreateRequest) (*uint, error) {
	existUser, err := s.deps.PostgresRepo.GetUserByUsername(req.Username)
	if err != nil {
		if !errors.Is(err, errmap.ErrmapNotFound) {
			log.Error(errors.Wrap(err, "[Service.CreateUser]: unable to check exist user"))
			return nil, errors.Wrap(err, "[Service.CreateUser]: unable to check exist user")
		}
	}

	if existUser != nil {
		log.Error(errors.Wrap(err, "[Service.CreateUser]: username already exists"))
		return nil, errors.Wrap(errmap.ErrmapConflict, "[Service.CreateUser]: username already exists")
	}

	password, err := s.deps.BcryptService.Hash(req.Password)
	if err != nil {
		log.Error(errors.Wrap(err, "[Service.CreateUser]: unable to hash password"))
		return nil, errors.Wrap(err, "[Service.CreateUser]: unable to hash password")
	}

	user := entity.User{
		Name:  	req.Name,
		Username: req.Username,
		Password: password,
		Role: 	constant.UserTypeUser,
	}

	userID, err := s.deps.PostgresRepo.CreateUser(user)
	if err != nil {
		log.Error(errors.Wrap(err, "[Service.CreateUser]: unable to create user"))
		return nil, errors.Wrap(err, "[Service.CreateUser]: unable to create user")
	}

	return userID, nil
}

// LoginUser login user
func (s *Service) LoginUser(username, password string) (*entity.User, error) {
	user, err := s.deps.PostgresRepo.GetUserByUsername(username)
	if err != nil {
		if errors.Is(err, errmap.ErrmapNotFound) {
			log.Error(errors.Wrapf(err, "[Service.LoginUser]: user with username %s not found", username))
			return nil, errmap.ErrmapNotFound
		}
		log.Error(errors.Wrap(err, "[Service.LoginUser]: unable to get user by username"))
		return nil, errors.Wrap(err, "[Service.LoginUser]: unable to get user by username")
	}

	if err := s.deps.BcryptService.Compare(password, user.Password); err != nil {
		log.Error(errors.Wrap(err, "[Service.LoginUser]: invalid password"))
		return nil, errmap.ErrmapInvalidPassword
	}
	return user, nil
}

// GetUserByID retrieves a user by ID
func (s *Service) GetUserByID(userID uint) (*entity.UserResponse, error) {
	user, err := s.deps.PostgresRepo.GetUserByID(userID)
	if err != nil {
		if errors.Is(err, errmap.ErrmapNotFound) {
			log.Error(errors.Wrap(err, "[Service.GetUserByID]: user not found"))
			return nil, errmap.ErrmapNotFound
		}
		
		log.Error(errors.Wrap(err, "[Service.GetUserByID]: user not found"))
		return nil, errors.Wrap(err, "[Service.GetUserByID]: user not found")
	}
	return user, nil
}

// UpdateUser updates user
func (s *Service) UpdateUser(req entity.UserUpdateRequest) error {
	_, err := s.deps.PostgresRepo.GetUserByID(req.ID)
	if err != nil {
		if errors.Is(err, errmap.ErrmapNotFound) {
			log.Error(errors.Wrap(err, "[Service.UpdateUser]: user not found"))
			return errmap.ErrmapNotFound
		}

		log.Error(errors.Wrap(err, "[Service.UpdateUser]: unable to get user by id"))
		return errors.Wrap(err, "[Service.UpdateUser]: unable to get user by id")
	}

	user := entity.User{
		ID: req.ID,
		Name: req.Name,
	}


	if err := s.deps.PostgresRepo.UpdateUser(user); err != nil {
		log.Error(errors.Wrap(err, "[Service.UpdateUser]: unable to update user"))
		return errors.Wrap(err, "[Service.UpdateUser]: unable to update user")
	}
	return nil
}

// DeleteUser deletes a user by ID
func (s *Service) DeleteUser(userID uint) error {
	_, err := s.deps.PostgresRepo.GetUserByID(userID)
	if err != nil {
		if errors.Is(err, errmap.ErrmapNotFound) {
			log.Error(errors.Wrap(err, "[Service.DeleteUser]: user not found"))
			return errmap.ErrmapNotFound
		}
		log.Error(errors.Wrap(err, "[Service.DeleteUser]: unable to get user by id"))
		return errors.Wrap(err, "[Service.DeleteUser]: unable to get user by id")
	}

	if err := s.deps.PostgresRepo.DeleteUser(userID); err != nil {
		log.Error(errors.Wrap(err, "[Service.DeleteUser]: unable to delete user"))
		return errors.Wrap(err, "[Service.DeleteUser]: unable to delete user")
	}
	return nil
}
