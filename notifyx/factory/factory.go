package factory

import (
	"fmt"
	"os"

	"github.com/Mozart-SymphonIA/infra-mz/notifyx"
	"github.com/Mozart-SymphonIA/infra-mz/notifyx/pq"
)

// NewListenerFromEnv construye un Listener usando las variables de entorno estándar de Postgres.
// Mismas variables que dbx/factory: POSTGRES_USER, POSTGRES_PASSWORD, POSTGRES_HOST, POSTGRES_PORT.
func NewListenerFromEnv(dbName string) notifyx.Listener {
	return pq.New(connStr(dbName))
}

// NewPublisherFromEnv construye un Publisher usando las variables de entorno estándar de Postgres.
func NewPublisherFromEnv(dbName string) (notifyx.Publisher, error) {
	return pq.NewPublisher(connStr(dbName))
}

func connStr(dbName string) string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		firstEnv("POSTGRES_USER", "postgres"),
		firstEnv("POSTGRES_PASSWORD", "postgres"),
		firstEnv("POSTGRES_HOST", "localhost"),
		firstEnv("POSTGRES_PORT", "5432"),
		dbName,
	)
}

func firstEnv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
