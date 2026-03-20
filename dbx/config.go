package dbx

import "os"

type Provider string

const (
	ProviderSQL      Provider = "sqlserver"
	ProviderPostgres Provider = "postgres"
)

type Config struct {
	Provider Provider
	URL      string
}

func DefaultConfig() Config {
	return Config{
		Provider: ProviderPostgres,
	}
}

func FromEnv(connectionName string) Config {
	cfg := DefaultConfig()
	if v := os.Getenv("DB_PROVIDER"); v != "" {
		cfg.Provider = Provider(v)
	}
	if v := os.Getenv(connectionName); v != "" { //todo validate env var name
		cfg.URL = v
	}
	return cfg
}
