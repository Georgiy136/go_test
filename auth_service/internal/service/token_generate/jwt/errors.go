package jwt

import "errors"

var (
	TokenIsExpiredError = errors.New("token_generate is expired")
)
