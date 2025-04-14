package repositories

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type FollowerRepository interface {
	FollowUser(ctx context.Context, followerID, followingID int64) error
	UnfollowUser(ctx context.Context, followerID, followingID int64) error
	IsFollowing(ctx context.Context, followerID, followingID int64) (bool, error)
	GetFollowers(ctx context.Context, userID int64) ([]int64, error)
	GetFollowing(ctx context.Context, userID int64) ([]int64, error)
	GetFollowersCount(ctx context.Context, userID int64) (int64, error)
	GetFollowingCount(ctx context.Context, userID int64) (int64, error)
}

type followerRepositoryImpl struct {
	db *pgxpool.Pool
}

func NewFollowerRepository(db *pgxpool.Pool) FollowerRepository {
	return &followerRepositoryImpl{db: db}
}

// Insert a new follow relationship.
func (r *followerRepositoryImpl) FollowUser(ctx context.Context, followerID, followingID int64) error {
	_, err := r.db.Exec(ctx,
		`INSERT INTO followers (follower_id, following_id)
		 VALUES ($1, $2)`,
		followerID, followingID,
	)
	return err
}

// Removes a follow relationship.
func (r *followerRepositoryImpl) UnfollowUser(ctx context.Context, followerID, followingID int64) error {
	_, err := r.db.Exec(ctx,
		`DELETE FROM followers
		 WHERE follower_id = $1 AND following_id = $2`,
		followerID, followingID,
	)
	return err
}

// Checks if a user is following another user.
func (r *followerRepositoryImpl) IsFollowing(ctx context.Context, followerID, followingID int64) (bool, error) {
	var exists bool
	err := r.db.QueryRow(ctx,
		`SELECT EXISTS (
			SELECT 1
			FROM followers
			WHERE follower_id = $1 AND following_id = $2
		)`,
		followerID, followingID,
	).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

// Retrieves a list of the user ids of all follower of a given user.
func (r *followerRepositoryImpl) GetFollowers(ctx context.Context, userID int64) ([]int64, error) {
	rows, err := r.db.Query(ctx,
		`SELECT follower_id
		 FROM followers
		 WHERE following_id = $1`,
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var followerIDs []int64
	for rows.Next() {
		var followerID int64
		if err := rows.Scan(&followerID); err != nil {
			return nil, fmt.Errorf("failed to scan follower ID: %w", err)
		}
		followerIDs = append(followerIDs, followerID)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over follower rows: %w", err)
	}

	return followerIDs, nil
}

// Retrieves a list of user ids of users that a given user is following.
func (r *followerRepositoryImpl) GetFollowing(ctx context.Context, userID int64) ([]int64, error) {
	rows, err := r.db.Query(ctx,
		`SELECT following_id
		 FROM followers
		 WHERE follower_id = $1`,
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var followingIDs []int64
	for rows.Next() {
		var followingID int64
		if err := rows.Scan(&followingID); err != nil {
			return nil, fmt.Errorf("failed to scan following ID: %w", err)
		}
		followingIDs = append(followingIDs, followingID)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over following rows: %w", err)
	}

	return followingIDs, nil
}

// Retrieves the count of followers for a given user.
func (r *followerRepositoryImpl) GetFollowersCount(ctx context.Context, userID int64) (int64, error) {
	var count int64
	err := r.db.QueryRow(ctx,
		`SELECT COUNT(*)
		 FROM followers
		 WHERE following_id = $1`,
		userID,
	).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

// Retrieves the count of users that a given user is following.
func (r *followerRepositoryImpl) GetFollowingCount(ctx context.Context, userID int64) (int64, error) {
	var count int64
	err := r.db.QueryRow(ctx,
		`SELECT COUNT(*)
		 FROM followers
		 WHERE follower_id = $1`,
		userID,
	).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}
