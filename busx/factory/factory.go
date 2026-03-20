package factory

import (
	"fmt"
	"os"

	"github.com/Mozart-SymphonIA/infra-mz/busx"
	"github.com/Mozart-SymphonIA/infra-mz/busx/rabbit"
)

func NewMinimalBundle() (*busx.Bundle, error) {
	return NewBundleFromEnv()
}

func NewBundleFromEnv() (*busx.Bundle, error) {
	user := firstEnv("RABBITMQ_USER", "guest")
	pass := firstEnv("RABBITMQ_PASSWORD", "guest")
	host := firstEnv("RABBITMQ_HOST", "rabbitmq")
	port := firstEnv("RABBITMQ_PORT", "5672")

	url := fmt.Sprintf("amqp://%s:%s@%s:%s/", user, pass, host, port)

	cfg := busx.Config{
		Provider: busx.ProviderRabbit,
		URL:      url,
	}
	return NewBundle(cfg)
}

func firstEnv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func NewBundle(cfg busx.Config) (*busx.Bundle, error) {
	switch cfg.Provider {
	case busx.ProviderRabbit:
		return rabbit.BuildRabbit(cfg)
	default:
		return nil, fmt.Errorf("busx: unknown provider %q", cfg.Provider)
	}
}
