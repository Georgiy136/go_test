package usecase

import (
	"context"
	"myapp/internal/models"
	"time"

	"github.com/google/uuid"
)

type GoodsStrore interface {
	CreateGoods(ctx context.Context, p models.Operator) error
	GetAllGoods(ctx context.Context) ([]models.Operator, error)
	DeleteGoods(ctx context.Context, id uuid.UUID) error
	UpdateGoods(ctx context.Context, id uuid.UUID, p models.Operator) (*models.Operator, error)
	GetOneGoods(ctx context.Context, id uuid.UUID) (*models.Operator, error)
}

type RoleStore interface {
	GetRoleRights(ctx context.Context, role string) ([]string, error)
	AddRoleRights(ctx context.Context, role string, rights []string, period time.Duration) error
}
