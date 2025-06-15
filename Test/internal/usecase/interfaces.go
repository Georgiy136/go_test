package usecase

import (
	"context"
	"myapp/internal/models"
	"time"

	"github.com/google/uuid"
)

type GoodsStrore interface {
	CreateGoods(ctx context.Context, p models.Goods) error
	GetAllGoods(ctx context.Context) ([]models.Goods, error)
	DeleteGoods(ctx context.Context, id uuid.UUID) error
	UpdateGoods(ctx context.Context, id uuid.UUID, p models.Goods) (*models.Goods, error)
	GetOneGoods(ctx context.Context, id uuid.UUID) (*models.Goods, error)
}

type RoleStore interface {
	GetRoleRights(ctx context.Context, role string) ([]string, error)
	AddRoleRights(ctx context.Context, role string, rights []string, period time.Duration) error
}
