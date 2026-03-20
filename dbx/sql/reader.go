package sql

import (
	"context"
	dbsql "database/sql"
	"errors"
	"fmt"
	"github.com/Mozart-SymphonIA/infra-mz/dbx"
)

var (
	_ dbx.Reader = (*sqlReader)(nil)
)

type sqlReader struct {
	c *dbsql.DB
}

type sqlRow struct {
	row *dbsql.Row
}

func (r sqlRow) Scan(dest ...any) error {
	return r.row.Scan(dest...)
}

func (s *sqlReader) Query(ctx context.Context, query string, args ...any) (string, error) {
	var config string
	err := s.c.QueryRowContext(ctx, query, args...).Scan(&config)
	if err != nil {
		if errors.Is(err, dbsql.ErrNoRows) {
			return "", dbsql.ErrNoRows
		}
		return "", fmt.Errorf("scan capability config: %w", err)
	}
	return config, nil
}

func (s *sqlReader) QueryRow(ctx context.Context, query string, args ...any) dbx.Row {
	return sqlRow{row: s.c.QueryRowContext(ctx, query, args...)}
}
