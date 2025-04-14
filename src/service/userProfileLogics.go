package service

import (
	"context"
	"errors"

	"fateh-ark/yapper-user-service/model"
	"fateh-ark/yapper-user-service/request"

	"github.com/jackc/pgx/v5"
)

func (s *userServiceImpl) UpsertUserProfile(ctx context.Context, userID int64, request *request.UserProfileReq) (*model.UserProfile, error) {
	if _, err := s.GetUserByID(ctx, userID); err != nil {
		return nil, err
	}

	existingProfile, err := s.profileRepo.GetUserProfileByUserID(ctx, userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			newProfile := &model.UserProfile{
				Bio:        request.Bio,
				Job:        request.Job,
				Location:   request.Location,
				WebsiteURL: request.WebsiteURL,
				BirthDate:  request.BirthDate,
			}
			if err := s.profileRepo.UpsertUserProfile(ctx, newProfile); err != nil {
				return nil, err
			}
			return s.profileRepo.GetUserProfileByUserID(ctx, userID)
		}
		return nil, err
	}

	if request.Bio != nil {
		existingProfile.Bio = request.Bio
	}
	if request.Job != nil {
		existingProfile.Job = request.Job
	}
	if request.Location != nil {
		existingProfile.Location = request.Location
	}
	if request.WebsiteURL != nil {
		existingProfile.WebsiteURL = request.WebsiteURL
	}
	if request.BirthDate != nil {
		existingProfile.BirthDate = request.BirthDate
	}
	// if request.BannerImageURL != nil {
	// 	existingProfile.BannerImageURL = request.BannerImageURL
	// }

	err = s.profileRepo.UpsertUserProfile(ctx, existingProfile)
	if err != nil {
		return nil, err
	}
	return s.profileRepo.GetUserProfileByUserID(ctx, userID)
}

func (s *userServiceImpl) GetUserProfileByUserID(ctx context.Context, userID int64) (*model.UserProfile, error) {
	_, err := s.GetUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	profile, err := s.profileRepo.GetUserProfileByUserID(ctx, userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrUserProfileNotFound
		}
		return nil, err
	}
	return profile, nil
}
