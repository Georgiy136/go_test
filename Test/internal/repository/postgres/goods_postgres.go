package postgres

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"github.com/uptrace/bun"
	"myapp/internal/models"
	"myapp/pkg/postgres"
)

type GoodsRepo struct {
	Bun *bun.DB
}

func NewGoodsRepo(pg *postgres.Postgres) *GoodsRepo {
	return &GoodsRepo{
		Bun: pg.Conn,
	}
}

func (db *GoodsRepo) CreateGoods(ctx context.Context, data models.DataFromRequestGoodsAdd) (*models.Goods, error) {
	rows, err := db.Bun.QueryContext(ctx,
		`
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
			`, data.ProjectID, data.Name, data.Description, data.Priority,
	)
	if err != nil {
		return nil, fmt.Errorf("Goods - CreateGoods - db.Bun.NewInsert: %w", err)
	}

	goods, err := GetDataFromDB[models.Goods](rows)
	if err != nil {
		return nil, fmt.Errorf("Goods - CreateGoods - GetDataFromDB: %w", err)
	}

	return goods, nil
}

func (db *GoodsRepo) ListGoods(ctx context.Context, data models.DataFromRequestGoodsList) (*models.GoodsListDBResponse, error) {
	rows, err := db.Bun.QueryContext(ctx,
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
	}

	return goodsList, nil
}

func (db *GoodsRepo) DeleteGoods(ctx context.Context, data models.DataFromRequestGoodsDelete) error {
	Goods := &models.Goods{}
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
	}
	return nil
}

func (db *GoodsRepo) UpdateGoods(ctx context.Context, data models.DataFromRequestGoodsUpdate) (*models.Goods, error) {
	err := db.Bun.NewUpdate().
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
	}

	return nil, nil
}
