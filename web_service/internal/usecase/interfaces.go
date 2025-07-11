package usecase

import (
	"context"
	"github.com/Georgiy136/go_test/web_service/internal/models"
)

type GoodsStrore interface {
	CreateGoods(ctx context.Context, data models.DataFromRequestGoodsAdd) (*models.Goods, error)
	ListGoods(ctx context.Context, data models.DataFromRequestGoodsList) (*models.GoodsList, error)
	DeleteGoods(ctx context.Context, data models.DataFromRequestGoodsDelete) (*models.Goods, error)
	UpdateGoods(ctx context.Context, data models.DataFromRequestGoodsUpdate) (*models.Goods, error)
	ReprioritizeGood(ctx context.Context, data models.DataFromRequestReprioritizeGood) (*models.Goods, error)
}

type GoodsCache interface {
	GetGoods(ctx context.Context, goodsID, projectID int) (*models.Goods, error)
	SaveGoods(ctx context.Context, goodsID, projectID int, goods models.Goods) error
	ClearGoods(ctx context.Context, goodsID, projectID int) error
}
