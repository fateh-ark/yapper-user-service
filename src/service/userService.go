package service

import (
	"context"

	"fateh-ark/yapper-user-service/model"
	"fateh-ark/yapper-user-service/repositories"
	"fateh-ark/yapper-user-service/request"
)

type UserService interface {
	CreateUser(ctx context.Context, request *request.CreateUserReq) (*model.User, error)
	GetUserByID(ctx context.Context, userID int64) (*model.User, error)
	GetUserByUsername(ctx context.Context, username string) (*model.User, error)
	GetUserByEmail(ctx context.Context, email string) (*model.User, error)
	UpdateUser(ctx context.Context, userID int64, request *request.UpdateUserReq) (*model.User, error)
	DeleteUser(ctx context.Context, userID int64) error

	FollowUser(ctx context.Context, userID int64, request *request.FollowReq) error
	UnfollowUser(ctx context.Context, userID int64, request *request.FollowReq) error
	GetFollowers(ctx context.Context, userID int64) ([]int64, error)
	GetFollowing(ctx context.Context, userID int64) ([]int64, error)
	Isfollowing(ctx context.Context, userID int64, request *request.FollowReq) (bool, error)
	GetFollowStats(ctx context.Context, userID int64) (*model.FollowStats, error)

	UpsertUserProfile(ctx context.Context, userID int64, request *request.UserProfileReq) (*model.UserProfile, error)
	GetUserProfileByUserID(ctx context.Context, userID int64) (*model.UserProfile, error)

	UpsertUserPreference(ctx context.Context, userID int64, request *request.UserPreferenceReq) (*model.UserPreference, error)
	GetUserPreferenceByUserID(ctx context.Context, userID int64) (*model.UserPreference, error)
}

type userServiceImpl struct {
	userRepo       repositories.UserRepository
	followerRepo   repositories.FollowerRepository
	profileRepo    repositories.UserProfileRepository
	preferenceRepo repositories.UserPreferenceRepository
}

func NewUserService(
	userRepo repositories.UserRepository,
	followerRepo repositories.FollowerRepository,
	profileRepo repositories.UserProfileRepository,
	preferenceRepo repositories.UserPreferenceRepository,
) UserService {
	return &userServiceImpl{
		userRepo:       userRepo,
		followerRepo:   followerRepo,
		profileRepo:    profileRepo,
		preferenceRepo: preferenceRepo,
	}
}
