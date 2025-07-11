package jwt

import "errors"

var (
	TokenIsExpiredError = errors.New("token is expired")
)
