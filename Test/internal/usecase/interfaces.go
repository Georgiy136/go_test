package usecase

import (
	"context"
	"myapp/internal/models"
	"time"
)

type GoodsStrore interface {
	CreateGoods(ctx context.Context, data models.DataFromRequestGoodsAdd) (*models.Goods, error)
	ListGoods(ctx context.Context, data models.DataFromRequestGoodsList) (*models.GoodsListDBResponse, error)
	DeleteGoods(ctx context.Context, data models.DataFromRequestGoodsDelete) error
	UpdateGoods(ctx context.Context, data models.DataFromRequestGoodsUpdate) (*models.Goods, error)
}

type GoodsCache interface {
	GetRoleRights(ctx context.Context, role string) ([]string, error)
	AddRoleRights(ctx context.Context, role string, rights []string, period time.Duration) error
}
