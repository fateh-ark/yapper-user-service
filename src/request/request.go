package request

import "time"

type CreateUserReq struct {
	Username    string `json:"username"`
	Email       string `json:"email"`
	DisplayName string `json:"display_name"`
}

type UpdateUserReq struct {
	Username    *string `json:"username"`
	Email       *string `json:"email"`
	DisplayName *string `json:"display_name"`
}

type FollowReq struct {
	FollowingID int64 `json:"following_id"`
}

type UserProfileReq struct {
	Bio        *string    `json:"bio"`         // Use pointer for nullable fields
	Job        *string    `json:"job"`         // Use pointer for nullable fields
	Location   *string    `json:"location"`    // Use pointer for nullable fields
	WebsiteURL *string    `json:"website_url"` // Use pointer for nullable fields
	BirthDate  *time.Time `json:"birth_date"`
}

type UserPreferenceReq struct {
	NotificationsEnabled bool `json:"notifications_enabled"`
	AccountPrivate       bool `json:"account_private"`
}
