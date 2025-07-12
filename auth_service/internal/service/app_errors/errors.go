package app_errors

import "errors"

var (
	TokenIsExpiredError      = errors.New("token_generate is expired")
	UserNotFoundError        = errors.New("User not found")
	SessionUserNotFoundError = errors.New("Session user not found")
)
