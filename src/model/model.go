package model

import "time"

// User represents the users table.
type User struct {
	ID              int64     `json:"id"`
	Username        string    `json:"username"`
	Email           string    `json:"email"`
	DisplayName     string    `json:"display_name"`
	ProfileImageURL *string   `json:"profile_image_url"` // Use pointer for nullable fields
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// Follower represents the followers table (for follow/unfollow functionality).
type Follower struct {
	FollowerID  int64     `json:"follower_id"`
	FollowingID int64     `json:"following_id"`
	FollowedAt  time.Time `json:"followed_at"`
}

// UserProfile represents the user_profiles table.
type UserProfile struct {
	UserID         int64      `json:"user_id"`
	Bio            *string    `json:"bio"`              // Use pointer for nullable fields
	Job            *string    `json:"job"`              // Use pointer for nullable fields
	Location       *string    `json:"location"`         // Use pointer for nullable fields
	WebsiteURL     *string    `json:"website_url"`      // Use pointer for nullable fields
	BirthDate      *time.Time `json:"birth_date"`       // Use pointer for nullable fields
	BannerImageURL *string    `json:"banner_image_url"` // Use pointer for nullable fields
	UpdatedAt      time.Time  `json:"updated_at"`
}

// UserPreference represents the user_preferences table.
type UserPreference struct {
	UserID               int64     `json:"user_id"`
	NotificationsEnabled bool      `json:"notifications_enabled"`
	AccountPrivate       bool      `json:"account_private"`
	UpdatedAt            time.Time `json:"updated_at"`
}

// Relationships
type FollowStats struct {
	FollowersCount int64 `json:"followers_count"`
	FollowingCount int64 `json:"following_count"`
}

// UserWithProfile combines User and UserProfile information.
type UserWithProfile struct {
	User
	Profile *UserProfile `json:"profile,omitempty"`
}

// UserWithPreferences combines User and UserPreference information.
type UserWithPreferences struct {
	User
	Preferences *UserPreference `json:"preferences,omitempty"`
}

// UserWithFollowers combines User and their followers/following.
// You might have separate structs for "Following" and "Followers" lists
// depending on your use case.
type UserWithFollowers struct {
	User
	Followers []int64 `json:"followers,omitempty"` // List of follower IDs
	Following []int64 `json:"following,omitempty"` // List of following IDs
}
