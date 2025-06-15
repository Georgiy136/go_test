package models

import (
	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type Goods struct {
	bun.BaseModel `bun:"table:goods"`

	Id         uuid.UUID `bun:"uuid"`
	FirstName  string    `bun:"first_name"`
	LastName   string    `bun:"last_name"`
	Patronymic string    `bun:"patronymic"`
	City       string    `bun:"city"`
	Phone      string    `bun:"phone"`
	Email      string    `bun:"email"`
	Password   string    `bun:"password"`
}
