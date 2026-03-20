package sql

import (
	"context"
	dbsql "database/sql"
)

func (s *sqlReader) QueryMap(ctx context.Context, query string, args ...any) (map[string]any, error) {
	rows, err := s.QueryMaps(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	if len(rows) == 0 {
		return nil, dbsql.ErrNoRows
	}
	return rows[0], nil
}

func (s *sqlReader) QueryMaps(ctx context.Context, query string, args ...any) ([]map[string]any, error) {
	rows, err := s.c.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}
	result := make([]map[string]any, 0)
	for rows.Next() {
		values := make([]any, len(columns))
		scanArgs := make([]any, len(columns))
		for i := range values {
			scanArgs[i] = &values[i]
		}
		if err := rows.Scan(scanArgs...); err != nil {
			return nil, err
		}
		entry := make(map[string]any, len(columns))
		for i, col := range columns {
			entry[col] = values[i]
		}
		result = append(result, entry)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return result, nil
}
