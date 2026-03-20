package sql

import (
	"context"
	"database/sql"
	"github.com/Mozart-SymphonIA/infra-mz/dbx"
)

var (
	_ dbx.Executer = (*sqlExecuter)(nil)
)

type sqlExecuter struct {
	c *sql.DB
}

func (s *sqlExecuter) Execute(ctx context.Context, query string, args ...any) error {
	_, err := s.c.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}
	return nil
}
func (s *sqlExecuter) MultipleExecute(ctx context.Context, queries []dbx.MultipleQuery) error {
	tx, err := s.c.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	for _, query := range queries {
		if _, err := tx.ExecContext(ctx, query.Query, query.Params...); err != nil {
			_ = tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}
