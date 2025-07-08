package service

import (
	"context"
)

type AuthStrore interface {
	CheckUser(ctx context.Context, userID int) error
}
