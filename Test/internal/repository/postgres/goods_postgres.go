package postgres

import (
	"context"
	"github.com/jackc/pgx/v5"
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
	query := `
			WITH ins_cte AS (
				INSERT INTO goods AS g (good_id,
				                        project_id,
				                        name,
				                        description,
				                        priority,
				                        created_at,
				                        deleted_at)
				SELECT nextval('good_sq') AS good_id,
    			       $1,
    			       $2,
    			       $3,
    			       $4,
				       NOW(),
				       null
				RETURNING g.*)
		
			SELECT jsonb_build_object('data', row_to_json(ins_cte))
			FROM ins_cte;`

	dbData, err := GetDataFromDB[models.GoodsUpdDBResponse](ctx, db.pgconn, query, data.ProjectID, data.Name, data.Description, data.Priority)
	if err != nil {
		return nil, err
	}

	return dbData.Data, nil
}

func (db *GoodsRepo) UpdateGoods(ctx context.Context, data models.DataFromRequestGoodsUpdate) (*models.Goods, error) {
	query := `
			WITH upd_cte AS (
					UPDATE goods AS g SET name = $3,
										  description = $4
					WHERE good_id = $1 AND project_id = $2
			RETURNING g.*)

			SELECT jsonb_build_object('data', row_to_json(upd_cte))
			FROM upd_cte;`

	dbData, err := GetDataFromDB[models.GoodsUpdDBResponse](ctx, db.pgconn, query, data.GoodID, data.ProjectID, data.Name, data.Description)
	if err != nil {
		return nil, err
	}

	return dbData.Data, nil
}

func (db *GoodsRepo) DeleteGoods(ctx context.Context, data models.DataFromRequestGoodsDelete) (*models.Goods, error) {
	query := `
			WITH upd_cte AS (
					UPDATE goods AS g SET deleted_at = $3
					WHERE good_id = $1 AND project_id = $2
			RETURNING g.*)

			SELECT jsonb_build_object('data', row_to_json(upd_cte))
			FROM upd_cte;`

	dbData, err := GetDataFromDB[models.GoodsUpdDBResponse](ctx, db.pgconn, query, data.GoodID, data.ProjectID, data.DeletedAt)
	if err != nil {
		return nil, err
	}

	return dbData.Data, nil
}

func (db *GoodsRepo) ListGoods(ctx context.Context, data models.DataFromRequestGoodsList) (*models.GoodsList, error) {
	query := `WITH goods_cte AS (SELECT g.good_id,
										g.project_id,
										g.name,
										p.name as project_name,
										g.description,
										g.priority,
										g.created_at,
										g.deleted_at
  								 FROM goods AS g
								 INNER JOIN projects p ON p.project_id = g.project_id
  								 WHERE g.good_id = COALESCE($1, g.good_id)
  								 AND g.project_id = COALESCE($2, g.project_id)
  								 LIMIT $3 OFFSET $4)

		SELECT JSONB_BUILD_OBJECT('data', (
				SELECT JSONB_BUILD_OBJECT('meta', JSONB_BUILD_OBJECT('total',  (SELECT COUNT(*) FROM goods_cte),
										  							 'remove', (SELECT COUNT(*) FROM goods_cte g WHERE g.deleted_at IS NOT NULL),
										  							 'limit',   $3,
   										  							 'offset',  $4),
                                          'goods', JSONB_AGG(c))
				FROM (SELECT g.good_id,
							 g.project_id,
							 g.project_name,
							 g.name,
							 g.description,
							 g.priority,
							 g.created_at,
							 g.deleted_at
				FROM goods_cte g) c));`

	dbData, err := GetDataFromDB[models.GoodsListDBResponse](ctx, db.pgconn, query, data.GoodsID, data.ProjectID, data.Limit, data.Offset)
	if err != nil {
		return nil, err
	}

	return dbData.Data, nil
}

func (db *GoodsRepo) ReprioritizeGood(ctx context.Context, data models.DataFromRequestReprioritizeGood) (*models.Goods, error) {
	query := `
			WITH upd_cte AS (
					UPDATE goods AS g SET priority = $3
					WHERE good_id = $1 AND project_id = $2
			RETURNING g.*)

			SELECT jsonb_build_object('data', row_to_json(upd_cte))
			FROM upd_cte;`

	dbData, err := GetDataFromDB[models.GoodsUpdDBResponse](ctx, db.pgconn, query, data.GoodID, data.ProjectID, data.Priority)
	if err != nil {
		return nil, err
	}

	return dbData.Data, nil
}
