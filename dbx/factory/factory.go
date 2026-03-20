package factory

import (
	"fmt"
	"os"

	"github.com/Mozart-SymphonIA/infra-mz/dbx"
	"github.com/Mozart-SymphonIA/infra-mz/dbx/sql"
)

func NewBundle(connectionName string) (*dbx.Bundle, error) {
	// Reusing NewBundleFromEnv logic
	cfg := dbx.FromEnv(connectionName)
	return NewBundleWithConfig(cfg)
}

func NewBundleFromEnv(dbName string) (*dbx.Bundle, error) {
	// Standard Postgres Env Vars
	user := firstEnv("POSTGRES_USER", "postgres")
	pass := firstEnv("POSTGRES_PASSWORD", "postgres")
	host := firstEnv("POSTGRES_HOST", "localhost")
	port := firstEnv("POSTGRES_PORT", "5432")

	// Format: postgres://user:pass@host:port/dbname?sslmode=disable
	// Assuming disable ssl for dev/internal docker
	sslmode := "disable"
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		user, pass, host, port, dbName, sslmode)

	cfg := dbx.Config{
		Provider: dbx.ProviderPostgres,
		URL:      connStr,
	}
	return NewBundleWithConfig(cfg)
}

func firstEnv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func NewBundleWithConfig(cfg dbx.Config) (*dbx.Bundle, error) {
	switch cfg.Provider {
	case dbx.ProviderPostgres:
		return sql.BuildSQLDB(cfg)
	default:
		return nil, fmt.Errorf("dbx: unknown provider %q", cfg.Provider)
	}
}
