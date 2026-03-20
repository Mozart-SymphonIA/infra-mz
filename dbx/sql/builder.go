package sql

import (
	"github.com/Mozart-SymphonIA/infra-mz/dbx"
)

func BuildSQLDB(cfg dbx.Config) (*dbx.Bundle, error) {
	conn, err := openConnection(cfg.URL)
	if err != nil {
		return nil, err
	}

	return &dbx.Bundle{
		Conn:      &sqlConn{c: conn},
		Inspector: &sqlConnectionInspector{c: conn},
		Executer:  &sqlExecuter{c: conn},
		Reader:    &sqlReader{c: conn},
	}, nil
}
