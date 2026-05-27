package dbx

import "context"

type Conn interface {
	Close() error
}

type ConnectionInspector interface {
	Ping(ctx context.Context) error
}

type Executer interface {
	Execute(ctx context.Context, query string, args ...any) error
	MultipleExecute(ctx context.Context, queries []MultipleQuery) error
}

type Row interface {
	Scan(dest ...any) error
}

type Reader interface {
	Query(ctx context.Context, query string, args ...any) (string, error)
	QueryAsTOON(ctx context.Context, section string, query string, args ...any) (string, error)
	QueryRow(ctx context.Context, query string, args ...any) Row
	QueryRows(ctx context.Context, query string, args ...any) (Rows, error)
	QueryStruct(ctx context.Context, query string, dest any, args ...any) error
	QueryStructs(ctx context.Context, query string, dest any, args ...any) error
	QueryMap(ctx context.Context, query string, args ...any) (map[string]any, error)
	QueryMaps(ctx context.Context, query string, args ...any) ([]map[string]any, error)
}

type Rows interface {
	Next() bool
	Scan(dest ...any) error
	Close() error
	Err() error
}

type Bundle struct {
	Conn      Conn
	Inspector ConnectionInspector
	Executer  Executer
	Reader    Reader
}

type MultipleQuery struct {
	Query  string
	Params []any
}
