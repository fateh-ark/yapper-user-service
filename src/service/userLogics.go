package service

import (
	"context"
	"errors"

	"fateh-ark/yapper-user-service/model"
	"fateh-ark/yapper-user-service/request"

	"github.com/jackc/pgx/v5"
)

func (s *userServiceImpl) CreateUser(ctx context.Context, request *request.CreateUserReq) (*model.User, error) {
	if existingUser, _ := s.GetUserByUsername(ctx, request.Username); existingUser != nil {
		return nil, ErrUsernameAlreadyInUse
	}

	if existingUser, _ := s.GetUserByEmail(ctx, request.Username); existingUser != nil {
		return nil, ErrEmailAlreadyInUse
	}

	user := &model.User{
		Username:    request.Username,
		Email:       request.Email,
		DisplayName: request.DisplayName,
	}

	if err := s.userRepo.CreateUser(ctx, user); err != nil {
		return nil, err
	}

	return s.userRepo.GetUserByUsername(ctx, request.Username) // Retrieve the created user with ID
}

func (s *userServiceImpl) GetUserByID(ctx context.Context, userID int64) (*model.User, error) {
	user, err := s.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return user, nil
}

func (s *userServiceImpl) GetUserByUsername(ctx context.Context, username string) (*model.User, error) {
	user, err := s.userRepo.GetUserByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return user, nil
}

func (s *userServiceImpl) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	user, err := s.userRepo.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return user, nil
}

func (s *userServiceImpl) UpdateUser(ctx context.Context, userID int64, request *request.UpdateUserReq) (*model.User, error) {
	existingUser, err := s.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	if request.Username != nil {
		if existingUser, _ := s.GetUserByUsername(ctx, *request.Username); existingUser != nil {
			return nil, ErrUsernameAlreadyInUse
		}
		existingUser.Username = *request.Username
	}
	if request.Email != nil {
		if existingUser, _ := s.GetUserByEmail(ctx, *request.Email); existingUser != nil {
			return nil, ErrEmailAlreadyInUse
		}
		existingUser.Email = *request.Email
	}
	if request.DisplayName != nil {
		existingUser.DisplayName = *request.DisplayName
	}
	// if user.ProfileImageURL != nil {
	// 	existingUser.ProfileImageURL = user.ProfileImageURL
	// }

	err = s.userRepo.UpdateUser(ctx, existingUser)
	if err != nil {
		return nil, err
	}
	return s.userRepo.GetUserByID(ctx, userID) // Retrieve the updated user
}

func (s *userServiceImpl) DeleteUser(ctx context.Context, userID int64) error {
	if _, err := s.userRepo.GetUserByID(ctx, userID); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrUserNotFound
		}
		return err
	}

	err := s.userRepo.DeleteUser(ctx, userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrUserNotFound
		}
		return err
	}
	return nil
}
