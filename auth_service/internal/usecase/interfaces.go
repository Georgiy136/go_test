package usecase

import (
	"context"
	"github.com/Georgiy136/go_test/web_service/internal/models"
)

type AuthStore interface {
	CreateGoods(ctx context.Context, data models.DataFromRequestGoodsAdd) (*models.Goods, error)
	ListGoods(ctx context.Context, data models.DataFromRequestGoodsList) (*models.GoodsList, error)
	DeleteGoods(ctx context.Context, data models.DataFromRequestGoodsDelete) (*models.Goods, error)
	UpdateGoods(ctx context.Context, data models.DataFromRequestGoodsUpdate) (*models.Goods, error)
	ReprioritizeGood(ctx context.Context, data models.DataFromRequestReprioritizeGood) (*models.Goods, error)
}
