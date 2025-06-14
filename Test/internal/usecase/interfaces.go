package usecase

import (
	"context"
	"myapp/internal/models"
	"time"

	"github.com/google/uuid"
)

type OperatorStrore interface {
	CreateOperator(ctx context.Context, p models.Operator) error
	GetAllOperators(ctx context.Context) ([]models.Operator, error)
	DeleteOperator(ctx context.Context, id uuid.UUID) error
	UpdateOperator(ctx context.Context, id uuid.UUID, p models.Operator) (*models.Operator, error)
	GetOneOperator(ctx context.Context, id uuid.UUID) (*models.Operator, error)
}

type RoleStore interface {
	GetRoleRights(ctx context.Context, role string) ([]string, error)
	AddRoleRights(ctx context.Context, role string, rights []string, period time.Duration) error
}
