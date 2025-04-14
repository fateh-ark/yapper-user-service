package service

import (
	"context"
	"errors"

	"fateh-ark/yapper-user-service/model"
	"fateh-ark/yapper-user-service/request"

	"github.com/jackc/pgx/v5"
)

// UpsertUserPreference handles the business logic for creating or updating user preferences.
func (s *userServiceImpl) UpsertUserPreference(ctx context.Context, userID int64, request *request.UserPreferenceReq) (*model.UserPreference, error) {
	_, err := s.GetUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	existingPrefs, err := s.preferenceRepo.GetUserPreferenceByUserID(ctx, userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			newPreference := &model.UserPreference{
				NotificationsEnabled: request.NotificationsEnabled,
				AccountPrivate:       request.AccountPrivate,
			}
			if err := s.preferenceRepo.UpsertUserPreference(ctx, newPreference); err != nil {
				return nil, err
			}
			return s.preferenceRepo.GetUserPreferenceByUserID(ctx, userID)
		}
		return nil, err
	}

	existingPrefs.NotificationsEnabled = request.NotificationsEnabled
	existingPrefs.AccountPrivate = request.AccountPrivate

	if err = s.preferenceRepo.UpsertUserPreference(ctx, existingPrefs); err != nil {
		return nil, err
	}
	return s.preferenceRepo.GetUserPreferenceByUserID(ctx, userID)
}

func (s *userServiceImpl) GetUserPreferenceByUserID(ctx context.Context, userID int64) (*model.UserPreference, error) {
	if _, err := s.GetUserByID(ctx, userID); err != nil {
		return nil, err
	}

	prefs, err := s.preferenceRepo.GetUserPreferenceByUserID(ctx, userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrUserPrefNotFound
		}
		return nil, err
	}
	return prefs, nil
}
