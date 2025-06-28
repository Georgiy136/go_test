package clickhouse

import (
	"fmt"
	"github.com/Georgiy136/go_test/Reader_to_click/internal/proto_models"
	"github.com/Georgiy136/go_test/Reader_to_click/pkg/clickhouse"
)

type Clickhouse struct {
	click *clickhouse.Clickhouse
}

func NewClickhouse(click *clickhouse.Clickhouse) *Clickhouse {
	return &Clickhouse{
		click: click,
	}
}

func (c *Clickhouse) SaveLogsToClick(logs []proto_models.Log) error {
	const insertDataFormat = "Insert into %s values ($1)"

	if _, err := c.click.Conn.Exec(fmt.Sprintf(insertDataFormat, c.click.Cfg.Dbname), logs); err != nil {
		return fmt.Errorf("can not SaveLogsToClick, err: %w", err)
	}
	return nil
}
