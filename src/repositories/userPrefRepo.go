package repositories

import (
	"context"

	"fateh-ark/yapper-user-service/model"

	"github.com/jackc/pgx/v5/pgxpool"
)

type UserPreferenceRepository interface {
	UpsertUserPreference(ctx context.Context, preference *model.UserPreference) error
	GetUserPreferenceByUserID(ctx context.Context, userID int64) (*model.UserPreference, error)
}

type userPreferenceRepositoryImpl struct {
	db *pgxpool.Pool
}

func NewUserPreferenceRepository(db *pgxpool.Pool) UserPreferenceRepository {
	return &userPreferenceRepositoryImpl{db: db}
}

// Upserts a new user preferences or updates existing ones.
func (r *userPreferenceRepositoryImpl) UpsertUserPreference(ctx context.Context, preference *model.UserPreference) error {
	_, err := r.db.Exec(ctx,
		`INSERT INTO user_preferences (user_id, notifications_enabled, account_private)
		 VALUES ($1, $2, $3)
		 ON CONFLICT (user_id)
		 DO UPDATE SET notifications_enabled = $2, account_private = $3`,
		preference.UserID, preference.NotificationsEnabled, preference.AccountPrivate,
	)
	return err
}

// Retrieves user preferences by the user ID.
func (r *userPreferenceRepositoryImpl) GetUserPreferenceByUserID(ctx context.Context, userID int64) (*model.UserPreference, error) {
	preference := &model.UserPreference{}
	err := r.db.QueryRow(ctx,
		`SELECT user_id, notifications_enabled, account_private, updated_at
		 FROM user_preferences
		 WHERE user_id = $1`,
		userID,
	).Scan(&preference.UserID, &preference.NotificationsEnabled, &preference.AccountPrivate, &preference.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return preference, nil
}
