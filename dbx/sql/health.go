package sql

import (
	"context"
	"database/sql"
	"github.com/Mozart-SymphonIA/infra-mz/dbx"
)

var (
	_ dbx.ConnectionInspector = (*sqlConnectionInspector)(nil)
)

type sqlConnectionInspector struct {
	c *sql.DB
}

func (s *sqlConnectionInspector) Ping(ctx context.Context) error {
	return s.c.PingContext(ctx)
}
