package service

import "errors"

var (
	ErrUserNotFound          = errors.New("user not found")
	ErrUserAlreadyExist      = errors.New("user already exists")
	ErrUsernameAlreadyInUse  = errors.New("username already in use")
	ErrEmailAlreadyInUse     = errors.New("email already in use")
	ErrCannotFollowSelf      = errors.New("follower and following id cannot be the same")
	ErrIsAlreadyFollowing    = errors.New("user is already following")
	ErrIsAlreadyNotFollowing = errors.New("user is already not following")
	ErrUserPrefNotFound      = errors.New("user preferences not found")
	ErrUserProfileNotFound   = errors.New("user profile not found")
)
