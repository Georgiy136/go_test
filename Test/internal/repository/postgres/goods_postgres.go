package postgres

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	jsoniter "github.com/json-iterator/go"
	"myapp/internal/models"
	"myapp/pkg/postgres"
)

type GoodsRepo struct {
	Pgconn *pgx.Conn
}

func NewGoodsRepo(pg *postgres.Postgres) *GoodsRepo {
	return &GoodsRepo{
		Pgconn: pg.Pgconn,
	}
}

func (db *GoodsRepo) CreateGoods(ctx context.Context, data models.DataFromRequestGoodsAdd) (*models.Goods, error) {
	dataBytes, err := jsoniter.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("json marshal dataFromRequestGoodsAdd err: %w", err)
	}

	dbData, err := getDataFromDB[models.Goods](ctx, db.Pgconn,
		`SELECT * FROM goods_upd($1);`, dataBytes,
	)
	if err != nil {
		return nil, fmt.Errorf("CreateGoods getDataFromDB: %w", err)
	}

	return dbData, nil
}

func (db *GoodsRepo) ListGoods(ctx context.Context, data models.DataFromRequestGoodsList) (*models.GoodsListDBResponse, error) {
	/*rows, err := db.Bun.QueryContext(ctx,
		`
			WITH goods_cte AS (
				SELECT g.id,
				       g.project_id,
				       g.name,
				       g.description,
				       g.priority,
				       g.removed,
				       g.created_at
				FROM goods g
			    WHERE g.project_id = COALESCE(?, g.project_id))
				LIMIT COAL
			)

			WITH cte AS (
			    INSERT INTO goods AS g (project_id,
			                      		name,
			                      		description,
			                      		priority)
			    VALUES (?, ?, ?, ?)
			    RETURNING g.id,
			              g.project_id,
			              g.name,
			              g.description,
			              g.priority,
			              g.removed,
			              g.created_at
			)
			SELECT jsonb_build_object('data', row_to_json(cte))
			FROM cte;
			`, data.ID, data.ProjectID, data.Limit, data.Offset,
	)
	if err != nil {
		return nil, fmt.Errorf("Goods - CreateGoods - db.Bun.NewInsert: %w", err)
	}

	goodsList, err := GetDataFromDB[models.GoodsListDBResponse](rows)
	if err != nil {
		return nil, fmt.Errorf("Goods - CreateGoods - GetDataFromDB: %w", err)
	}*/

	return nil, nil
}

func (db *GoodsRepo) DeleteGoods(ctx context.Context, data models.DataFromRequestGoodsDelete) error {
	/*Goods := &models.Goods{}
	err := db.Bun.NewDelete().
		Model(Goods).
		Where(`uuid = ?`, data.ID).
		Returning("uuid").
		Scan(ctx, &data.ID)

	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("Goods - DeleteGoods - db.Bun.NewDelete: %s", fmt.Sprintf("оператора с id = %s не существует", data.ID))
		}
		logrus.Error(err)
		return fmt.Errorf("Goods - DeleteGoods - db.Bun.NewDelete: %w", err)
	}*/
	return nil
}

func (db *GoodsRepo) UpdateGoods(ctx context.Context, data models.DataFromRequestGoodsUpdate) (*models.Goods, error) {
	/*err := db.Bun.NewUpdate().
		Model(&data).
		Column("first_name", "last_name", "patronymic", "city", "phone", "email").
		Where(`uuid = ?`, data.ID).
		Returning("uuid, password").
		Scan(ctx, &data.ID)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("Goods - UpdateGoods - db.Bun.NewUpdate: %s", fmt.Sprintf("оператора с id = %s не существует", data.ID))
		}
		logrus.Error(err)
		return nil, fmt.Errorf("Goods - UpdateGoods - db.Bun.NewUpdate: %w", err)
	}*/

	return nil, nil
}
