package busx

import (
	"os"
	"strconv"
	"time"
)

type Provider string

const (
	ProviderRabbit Provider = "rabbitmq"
)

type Config struct {
	Provider          Provider
	URL               string
	Prefetch          int
	Heartbeat         time.Duration
	Locale            string
	PublisherConfirms bool
	Observer          Observer
}

func DefaultConfig() Config {
	return Config{
		Provider:          ProviderRabbit,
		URL:               "amqp://guest:guest@rabbitmq:5672/",
		Prefetch:          16,
		Heartbeat:         10 * time.Second,
		Locale:            "en_US",
		PublisherConfirms: true,
		Observer:          NopObserver{},
	}
}

func FromEnv() Config {
	cfg := DefaultConfig()
	if v := os.Getenv("BUS_PROVIDER"); v != "" {
		cfg.Provider = Provider(v)
	}
	cfg.URL = getEnvURL(cfg.URL)
	if v := os.Getenv("BUS_PREFETCH"); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			cfg.Prefetch = n
		}
	}
	if v := os.Getenv("BUS_HEARTBEAT_SEC"); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			cfg.Heartbeat = time.Duration(n) * time.Second
		}
	}
	if v := os.Getenv("BUS_PUBLISHER_CONFIRMS"); v == "0" || v == "false" {
		cfg.PublisherConfirms = false
	}
	return cfg
}

func getEnvURL(defaultURL string) string {
	if v := os.Getenv("BUS_URL"); v != "" {
		return v
	}
	if v := os.Getenv("AMQP_URL"); v != "" {
		return v
	}
	return defaultURL
}
