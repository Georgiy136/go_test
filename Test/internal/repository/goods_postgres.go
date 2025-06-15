package repository

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/sirupsen/logrus"
	"myapp/internal/models"
	"myapp/pkg/postgres"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"github.com/uptrace/bun"
)

func NewGoods(pg *postgres.Postgres) *Good {
	return &Good{
		Bun: pg.Conn,
	}
}

type Good struct {
	Bun *bun.DB
}

func (db *Good) CreateGoods(ctx context.Context, o models.Goods) error {
	_, err := db.Bun.NewInsert().Model(&o).Exec(ctx)
	if err != nil {
		logrus.Error(err)
		return fmt.Errorf("Goods - CreateGoods - db.Bun.NewInsert: %w", err)
	}
	return nil
}

func (db *Good) GetAllGoods(ctx context.Context) ([]models.Goods, error) {

	Goods := []models.Goods{}
	Goods := &models.Goods{}

	err := db.Bun.NewSelect().Model(Goods).Scan(ctx, &Goods)
	if err != nil {
		logrus.Error(err)
		return nil, fmt.Errorf("Goods - GetAllGoods - db.Bun.NewSelect: %w", err)
	}

	return Goods, nil
}

func (db *Good) DeleteGoods(ctx context.Context, id uuid.UUID) error {
	Goods := &models.Goods{}
	err := db.Bun.NewDelete().
		Model(Goods).
		Where(`uuid = ?`, id.String()).
		Returning("uuid").
		Scan(ctx, &id)

	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("Goods - DeleteGoods - db.Bun.NewDelete: %s", fmt.Sprintf("оператора с id = %s не существует", id.String()))
		}
		logrus.Error(err)
		return fmt.Errorf("Goods - DeleteGoods - db.Bun.NewDelete: %w", err)
	}
	return nil
}

func (db *Good) UpdateGoods(ctx context.Context, id uuid.UUID, p models.Goods) (*models.Goods, error) {

	err := db.Bun.NewUpdate().
		Model(&p).
		Column("first_name", "last_name", "patronymic", "city", "phone", "email").
		Where(`uuid = ?`, id.String()).
		Returning("uuid, password").
		Scan(ctx, &p.Id, &p.Password)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("Goods - UpdateGoods - db.Bun.NewUpdate: %s", fmt.Sprintf("оператора с id = %s не существует", id.String()))
		}
		logrus.Error(err)
		return nil, fmt.Errorf("Goods - UpdateGoods - db.Bun.NewUpdate: %w", err)
	}

	return &p, nil
}

func (db *Good) GetOneGoods(ctx context.Context, id uuid.UUID) (*models.Goods, error) {
	Goods := models.Goods{}

	err := db.Bun.NewSelect().Model(&Goods).Where("uuid = ?", id.String()).Scan(ctx)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("Goods - GetOneGoods - db.Bun.NewSelect: %s", fmt.Sprintf("оператора с id = %s не существует", id.String()))
		}
		logrus.Error(err)
		return nil, fmt.Errorf("Goods - GetOneGoods - db.Bun.NewSelect: %w", err)
	}

	return &Goods, nil
}

func (db *Good) GetGoodsById(ctx context.Context, id []string) ([]*models.Goods, error) {
	Goods := models.Goods{}
	Goods := []*models.Goods{}

	err := db.Bun.NewSelect().Model(&Goods).Where("uuid IN (?)", bun.In(id)).Scan(ctx, &Goods)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		logrus.Error(err)
		return nil, fmt.Errorf("Goods - GetGoodsById - db.Bun.NewSelect: %w", err)
	}
	return Goods, nil
}
