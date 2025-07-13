package app_errors

import "errors"

var (
	TokenIsExpiredError      = errors.New("token is expired")
	SessionUserNotFoundError = errors.New("user session not found")
	UserNotFoundError        = errors.New("user not found")
	UserAgentNotMatchInDB    = errors.New("User-Agent does not match in db")
)
