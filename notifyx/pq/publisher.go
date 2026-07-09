package pq

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"

	"github.com/Mozart-SymphonIA/infra-mz/notifyx"
)

type publisher struct {
	db *sql.DB
}

func NewPublisher(connStr string) (notifyx.Publisher, error) {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("notifyx publisher: open: %w", err)
	}
	db.SetMaxOpenConns(2)
	return &publisher{db: db}, nil
}

func (p *publisher) Publish(ctx context.Context, channel, payload string) error {
	_, err := p.db.ExecContext(ctx, "SELECT pg_notify($1, $2)", channel, payload)
	return err
}
