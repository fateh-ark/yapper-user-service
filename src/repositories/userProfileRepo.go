package repositories

import (
	"context"

	"fateh-ark/yapper-user-service/model"

	"github.com/jackc/pgx/v5/pgxpool"
)

type UserProfileRepository interface {
	UpsertUserProfile(ctx context.Context, profile *model.UserProfile) error
	GetUserProfileByUserID(ctx context.Context, userID int64) (*model.UserProfile, error)
}

type userProfileRepositoryImpl struct {
	db *pgxpool.Pool
}

func NewUserProfileRepository(db *pgxpool.Pool) UserProfileRepository {
	return &userProfileRepositoryImpl{db: db}
}

// Upserts a new user profile or updates an existing one.
func (r *userProfileRepositoryImpl) UpsertUserProfile(ctx context.Context, profile *model.UserProfile) error {
	_, err := r.db.Exec(ctx,
		`INSERT INTO user_profiles (user_id, bio, job, location, website_url, birth_date, banner_image_url)
		 VALUES ($1, $2, $3, $4, $5, $6, $7)
		 ON CONFLICT (user_id)
		 DO UPDATE SET
		 	bio = $2,
		 	job = $3,
		 	location = $4,
		 	website_url = $5,
		 	birth_date = $6,
		 	banner_image_url = $7`,
		profile.UserID, profile.Bio, profile.Job, profile.Location, profile.WebsiteURL, profile.BirthDate, profile.BannerImageURL,
	)
	return err
}

// Retrieves a user profile by the user ID.
func (r *userProfileRepositoryImpl) GetUserProfileByUserID(ctx context.Context, userID int64) (*model.UserProfile, error) {
	profile := &model.UserProfile{}
	err := r.db.QueryRow(ctx,
		`SELECT user_id, bio, job, location, website_url, birth_date, banner_image_url, updated_at
		 FROM user_profiles
		 WHERE user_id = $1`,
		userID,
	).Scan(&profile.UserID, &profile.Bio, &profile.Job, &profile.Location, &profile.WebsiteURL, &profile.BirthDate, &profile.BannerImageURL, &profile.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return profile, nil
}
