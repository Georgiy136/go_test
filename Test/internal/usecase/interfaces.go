package usecase

import (
	"context"
	"myapp/internal/models"
)

type GoodsStrore interface {
	CreateGoods(ctx context.Context, data models.DataFromRequestGoodsAdd) (*models.Goods, error)
	ListGoods(ctx context.Context, data models.DataFromRequestGoodsList) (*models.GoodsListDBResponse, error)
	DeleteGoods(ctx context.Context, data models.DataFromRequestGoodsDelete) error
	UpdateGoods(ctx context.Context, data models.DataFromRequestGoodsUpdate) (*models.Goods, error)
}

type GoodsCache interface {
	GetGoods(ctx context.Context, goodsID, projectID int) (*models.Goods, error)
	SaveGoods(ctx context.Context, goodsID, projectID int, goods models.Goods) error
	ClearGoods(ctx context.Context, goodsID, projectID int) error
}
