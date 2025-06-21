package postgres

import (
	"context"
	"fmt"
	"github.com/go-faster/errors"
	"github.com/jackc/pgx/v5"
	jsoniter "github.com/json-iterator/go"
	"myapp/internal/models"
	"myapp/internal/usecase"
	"myapp/pkg/postgres"
)

type GoodsRepo struct {
	pgconn *pgx.Conn
}

func NewGoodsRepo(pg *postgres.Postgres) usecase.GoodsStrore {
	return &GoodsRepo{
		pgconn: pg.Pgconn,
	}
}

func (db *GoodsRepo) CreateGoods(ctx context.Context, data models.DataFromRequestGoodsAdd) (*models.Goods, error) {
	dataJson, err := jsoniter.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("json marshal dataFromRequestGoodsAdd err: %w", err)
	}

	dbData, err := getDataFromDB[models.Goods](ctx, db.pgconn,
		`SELECT * FROM goods_upd($1);`, dataJson,
	)
	if err != nil {
		return nil, err
	}

	return dbData, nil
}

func (db *GoodsRepo) UpdateGoods(ctx context.Context, data models.DataFromRequestGoodsUpdate) (*models.Goods, error) {
	dataJson, err := jsoniter.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("json marshal dataFromRequestGoodsAdd err: %w", err)
	}

	dbData, err := getDataFromDB[models.Goods](ctx, db.pgconn,
		`SELECT * FROM goods_upd($1);`, dataJson,
	)
	if err != nil {
		return nil, errors.Wrap(err, "UpdateGoods getDataFromDB")
	}

	return dbData, nil
}

func (db *GoodsRepo) DeleteGoods(ctx context.Context, data models.DataFromRequestGoodsDelete) (*models.Goods, error) {
	dataJson, err := jsoniter.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("json marshal dataFromRequestGoodsAdd err: %w", err)
	}

	dbData, err := getDataFromDB[models.Goods](ctx, db.pgconn,
		`SELECT * FROM goods_upd($1);`, dataJson,
	)
	if err != nil {
		return nil, errors.Wrap(err, "CreateGoods getDataFromDB")
	}

	return dbData, nil
}

func (db *GoodsRepo) ListGoods(ctx context.Context, data models.DataFromRequestGoodsList) (*models.GoodsListDBResponse, error) {
	dbData, err := getDataFromDB[*models.GoodsListDBResponse](ctx, db.pgconn,
		`SELECT * FROM goods_list($1,$2,$3,$4);`, data.GoodsID, data.ProjectID, data.Limit, data.Offset,
	)
	if err != nil {
		return nil, errors.Wrap(err, "ListGoods getDataFromDB")
	}

	return *dbData, nil
}

func (db *GoodsRepo) ReprioritizeGood(ctx context.Context, data models.DataFromRequestReprioritizeGood) (*models.Goods, error) {
	dataJson, err := jsoniter.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("json marshal dataFromRequestGoodsAdd err: %w", err)
	}

	dbData, err := getDataFromDB[models.Goods](ctx, db.pgconn,
		`SELECT * FROM goods_upd($1);`, dataJson,
	)
	if err != nil {
		return nil, errors.Wrap(err, "UpdateGoods getDataFromDB")
	}

	return dbData, nil
}
