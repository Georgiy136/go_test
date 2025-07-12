package app_errors

import "errors"

var (
	TokenIsExpiredError      = errors.New("token is expired")
	SessionUserNotFoundError = errors.New("Session user not found")
	UserNotFoundError        = errors.New("User not found")
)
