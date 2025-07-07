package postgres

import (
	"context"
	"fmt"
	"github.com/Georgiy136/go_test/web_service/internal/errors/common"
	"github.com/Georgiy136/go_test/web_service/internal/models"
	"github.com/Georgiy136/go_test/web_service/internal/usecase"
	"github.com/Georgiy136/go_test/web_service/pkg/postgres"
	"github.com/jackc/pgx/v5"
	jsoniter "github.com/json-iterator/go"
	"strings"
)

type AuthPostgres struct {
	pgconn *pgx.Conn
}

func NewAuthPostgres(pg *postgres.Postgres) usecase.AuthStore {
	return &AuthPostgres{
		pgconn: pg.Pgconn,
	}
}

func (db *AuthPostgres) CreateGoods(ctx context.Context, data models.DataFromRequestGoodsAdd) (*models.Goods, error) {
	const sp = "goods_upd"

	dataJson, err := jsoniter.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("json marshal dataFromRequestGoodsAdd err: %w", err)
	}

	pg := new(postgres.PgSpec)
	pg.SetStoredProcedure(sp)
	pg.SetParams(dataJson)
	pg.SetUseFunction()

	dbData, err := GetDataFromDBAndUnmarshal[models.GoodsUpdDBResponse](ctx, db.pgconn, pg)
	if err != nil {
		return nil, err
	}

	return &dbData.Data, nil
}

func (db *AuthPostgres) UpdateGoods(ctx context.Context, data models.DataFromRequestGoodsUpdate) (*models.Goods, error) {
	const sp = "goods_upd"

	dataJson, err := jsoniter.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("json marshal dataFromRequestGoodsAdd err: %w", err)
	}

	pg := new(postgres.PgSpec)
	pg.SetStoredProcedure(sp)
	pg.SetParams(dataJson)
	pg.SetUseFunction()

	dbData, err := GetDataFromDBAndUnmarshal[models.GoodsUpdDBResponse](ctx, db.pgconn, pg)
	if err != nil {
		if strings.Contains(err.Error(), common.GoodsNotFoundDbError) {
			return nil, &common.CustomError{Description: err.Error(), Err: &common.NotFoundError}
		}
		return nil, err
	}

	return &dbData.Data, nil
}

func (db *AuthPostgres) DeleteGoods(ctx context.Context, data models.DataFromRequestGoodsDelete) (*models.Goods, error) {
	const sp = "goods_upd"

	dataJson, err := jsoniter.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("json marshal dataFromRequestGoodsAdd err: %w", err)
	}

	pg := new(postgres.PgSpec)
	pg.SetStoredProcedure(sp)
	pg.SetParams(dataJson)
	pg.SetUseFunction()

	dbData, err := GetDataFromDBAndUnmarshal[models.GoodsUpdDBResponse](ctx, db.pgconn, pg)
	if err != nil {
		if strings.Contains(err.Error(), common.GoodsNotFoundDbError) {
			return nil, &common.CustomError{Description: err.Error(), Err: &common.NotFoundError}
		}
		return nil, err
	}

	return &dbData.Data, nil
}

func (db *AuthPostgres) ListGoods(ctx context.Context, data models.DataFromRequestGoodsList) (*models.GoodsList, error) {
	const sp = "goods_list"

	pg := new(postgres.PgSpec)
	pg.SetStoredProcedure(sp)
	pg.SetParams(data.GoodsID, data.ProjectID, data.Limit, data.Offset)
	pg.SetUseFunction()

	dbData, err := GetDataFromDBAndUnmarshal[models.GoodsListDBResponse](ctx, db.pgconn, pg)
	if err != nil {
		return nil, err
	}

	return dbData.Data, nil
}

func (db *AuthPostgres) ReprioritizeGood(ctx context.Context, data models.DataFromRequestReprioritizeGood) (*models.Goods, error) {
	const sp = "goods_upd"

	dataJson, err := jsoniter.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("json marshal dataFromRequestGoodsAdd err: %w", err)
	}

	pg := new(postgres.PgSpec)
	pg.SetStoredProcedure(sp)
	pg.SetParams(dataJson)
	pg.SetUseFunction()

	dbData, err := GetDataFromDBAndUnmarshal[models.GoodsUpdDBResponse](ctx, db.pgconn, pg)
	if err != nil {
		if strings.Contains(err.Error(), common.GoodsNotFoundDbError) {
			return nil, &common.CustomError{Description: err.Error(), Err: &common.NotFoundError}
		}
		return nil, err
	}

	return &dbData.Data, nil
}
