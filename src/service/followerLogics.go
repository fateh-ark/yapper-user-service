package service

import (
	"context"
	"fateh-ark/yapper-user-service/model"
	"fateh-ark/yapper-user-service/request"
	"fmt"
)

func (s *userServiceImpl) FollowUser(ctx context.Context, userID int64, request *request.FollowReq) error {
	isFollowing, err := s.Isfollowing(ctx, userID, request)
	if err != nil {
		return err
	}
	if isFollowing == true {
		return ErrIsAlreadyFollowing
	}

	return s.followerRepo.FollowUser(ctx, userID, request.FollowingID)
}

func (s *userServiceImpl) UnfollowUser(ctx context.Context, userID int64, request *request.FollowReq) error {
	isFollowing, err := s.Isfollowing(ctx, userID, request)
	if err != nil {
		return err
	}
	if isFollowing == false {
		return ErrIsAlreadyNotFollowing
	}

	return s.followerRepo.UnfollowUser(ctx, userID, request.FollowingID)
}

// GetFollowers retrieves the IDs of users who follow a given user.
func (s *userServiceImpl) GetFollowers(ctx context.Context, userID int64) ([]int64, error) {
	_, err := s.GetUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	return s.followerRepo.GetFollowers(ctx, userID)
}

// GetFollowing retrieves the IDs of users that a given user is following.
func (s *userServiceImpl) GetFollowing(ctx context.Context, userID int64) ([]int64, error) {
	_, err := s.GetUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	return s.followerRepo.GetFollowing(ctx, userID)
}

func (s *userServiceImpl) Isfollowing(ctx context.Context, userID int64, request *request.FollowReq) (bool, error) {
	if userID == request.FollowingID {
		return false, ErrCannotFollowSelf
	}

	if _, err := s.GetUserByID(ctx, userID); err != nil {
		return false, fmt.Errorf("failed to check follower user: %w", err)
	}

	if _, err := s.GetUserByID(ctx, request.FollowingID); err != nil {
		return false, fmt.Errorf("failed to check following user: %w", err)
	}
	return s.followerRepo.IsFollowing(ctx, userID, request.FollowingID)
}

func (s *userServiceImpl) GetFollowStats(ctx context.Context, userID int64) (*model.FollowStats, error) {
	_, err := s.GetUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	followers, err := s.followerRepo.GetFollowersCount(ctx, userID)
	if err != nil {
		return nil, err
	}
	following, err := s.followerRepo.GetFollowingCount(ctx, userID)
	if err != nil {
		return nil, err
	}
	return &model.FollowStats{
		FollowersCount: followers,
		FollowingCount: following,
	}, nil
}
