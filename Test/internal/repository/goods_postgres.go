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
		`INSERT INTO goods(project_id, 
                  				 name, 
                  				 description, 
                 				 priority)
			   VALUES (?, ?, ?, ?)
			   RETURNING id,
						 project_id, 
						 name,
						 description, 
					     priority, 
						 removed, 
						 created_at
			   `,
		data.ProjectID, data.Name, data.Description, data.Priority,
	)

	if err != nil {
		return nil, fmt.Errorf("Goods - CreateGoods - db.Bun.NewInsert: %w", err)
	}

	res := models.Goods{}
	for rows.Next() {
		if err = rows.Scan(&res.Id, &res.ProjectID, &res.Name, &res.Description, &res.Priority, &res.Removed, &res.CreatedAt); err != nil {
			return nil, fmt.Errorf("Goods - CreateGoods - Scan: %w", err)
		}
	}

	return &res, nil
}

func (db *Good) ListGoods(ctx context.Context, data models.DataFromRequestGoodsList) ([]models.Goods, error) {
	Goods := []*models.Goods{}

	err := db.Bun.NewSelect().Model(&models.Goods{}).Where("uuid IN (?)", bun.In(data.ID)).Scan(ctx, &Goods)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		logrus.Error(err)
		return nil, fmt.Errorf("Goods - ListGoods - db.Bun.NewSelect: %w", err)
	}
	return nil, nil
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
