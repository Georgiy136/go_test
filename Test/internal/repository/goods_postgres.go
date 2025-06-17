package repository

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

type Good struct {
	Bun *bun.DB
}

func NewGoodsRepo(pg *postgres.Postgres) *Good {
	return &Good{
		Bun: pg.Conn,
	}
}

func (db *Good) CreateGoods(ctx context.Context, data models.DataFromRequestGoodsAdd) (*models.Goods, error) {
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

	result, err := GetDataFromDB[models.Goods](rows)
	if err != nil {
		return nil, fmt.Errorf("Goods - CreateGoods - GetDataFromDB: %w", err)
	}

	logrus.Infof("Goods - CreateGoods - db.Bun.NewInsert: %v", result)

	return result, nil
}

func (db *Good) ListGoods(ctx context.Context, data models.DataFromRequestGoodsList) (*models.GoodsListDBResponse, error) {
	res := models.GoodsListDBResponse{}

	//err := db.Bun.QueryRowContext(ctx,
	//	`SELECT id,
	//				  project_id,
	//				  name,
	//				  description,
	//				  priority,
	//				  removed,
	//				  created_at
	//			FROM goods
	//					  `,
	//	data.ProjectID, data.Name, data.Description, data.Priority,
	//).Scan(&goods.Id, &goods.ProjectID, &goods.Name, &goods.Description, &goods.Priority, &goods.Removed, &goods.CreatedAt)
	/*err := db.Bun.NewSelect().Model(&models.Goods{}).Where("uuid IN (?)", bun.In(data.ID)).Scan(ctx, &Goods)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		logrus.Error(err)
		return nil, fmt.Errorf("Goods - ListGoods - db.Bun.NewSelect: %w", err)
	}*/
	return &res, nil
}

func (db *Good) DeleteGoods(ctx context.Context, data models.DataFromRequestGoodsDelete) error {
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

func (db *Good) UpdateGoods(ctx context.Context, data models.DataFromRequestGoodsUpdate) (*models.Goods, error) {
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
