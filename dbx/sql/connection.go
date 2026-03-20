package sql

import (
	"database/sql"
	"errors"
	"github.com/Mozart-SymphonIA/infra-mz/dbx"
	"strings"
	"time"

	_ "github.com/lib/pq"
)

var (
	_ dbx.Conn = (*sqlConn)(nil)
)

type sqlConn struct {
	c *sql.DB
}

func (s *sqlConn) Close() error {
	return s.c.Close()
}

func openConnection(connString string) (*sql.DB, error) {
	if strings.TrimSpace(connString) == "" {
		return nil, errors.New("empty connection string")
	}

	// Default to postgres as we have migrated away from sqlserver
	db, err := sql.Open("postgres", connString)

	if err != nil {
		return nil, err
	}

	db.SetConnMaxLifetime(5 * time.Minute)
	db.SetConnMaxIdleTime(1 * time.Minute)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)

	return db, nil
}
