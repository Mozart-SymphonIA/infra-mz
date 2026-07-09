package factory

import (
	"fmt"
	"os"

	"github.com/Mozart-SymphonIA/infra-mz/notifyx"
	"github.com/Mozart-SymphonIA/infra-mz/notifyx/pq"
)

// NewFromEnv construye un Listener usando las variables de entorno estándar de Postgres.
// Mismas variables que dbx/factory: POSTGRES_USER, POSTGRES_PASSWORD, POSTGRES_HOST, POSTGRES_PORT.
func NewFromEnv(dbName string) notifyx.Listener {
	user := firstEnv("POSTGRES_USER", "postgres")
	pass := firstEnv("POSTGRES_PASSWORD", "postgres")
	host := firstEnv("POSTGRES_HOST", "localhost")
	port := firstEnv("POSTGRES_PORT", "5432")
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		user, pass, host, port, dbName,
	)
	return pq.New(connStr)
}

func firstEnv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
