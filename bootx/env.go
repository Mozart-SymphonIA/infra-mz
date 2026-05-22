package bootx

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

func LoadEnv() {
	if err := godotenv.Load(); err == nil {
		return
	}

	for i := 0; i < 2; i++ {
		if err := godotenv.Load(strings.Repeat("../", i+1) + ".env"); err == nil {
			return
		}
	}
}

func FirstEnv(keys ...string) string {
	for _, key := range keys {
		if key == "" {
			continue
		}
		if val := strings.TrimSpace(os.Getenv(key)); val != "" {
			return val
		}
	}
	return ""
}

func RequireEnv(keys ...string) (string, error) {
	if val := FirstEnv(keys...); val != "" {
		return val, nil
	}
	switch len(keys) {
	case 0:
		return "", fmt.Errorf("environment variable name is required")
	case 1:
		return "", fmt.Errorf("%s is required", keys[0])
	default:
		return "", fmt.Errorf("one of %s is required", strings.Join(keys, ", "))
	}
}

// ParseDuration parses a duration string (Go format, e.g. "15m") with an optional
// plain-seconds fallback (e.g. "900" is treated as "900s").
// Returns fallback if raw is empty.
func ParseDuration(raw string, fallback time.Duration) (time.Duration, error) {
	if strings.TrimSpace(raw) == "" {
		return fallback, nil
	}
	if d, err := time.ParseDuration(raw); err == nil {
		return d, nil
	}
	secs, err := time.ParseDuration(raw + "s")
	if err != nil {
		return 0, fmt.Errorf("invalid duration %q", raw)
	}
	return secs, nil
}

// LoadFileFromEnv reads a file's contents from one of two env vars:
// - inlineEnvVar: the value itself (e.g. an inline PEM key)
// - pathEnvVar: a filesystem path to read the file from
// Returns an error if neither is set.
func LoadFileFromEnv(inlineEnvVar, pathEnvVar string) ([]byte, error) {
	if raw := strings.TrimSpace(os.Getenv(inlineEnvVar)); raw != "" {
		return []byte(raw), nil
	}
	path := strings.TrimSpace(os.Getenv(pathEnvVar))
	if path == "" {
		return nil, fmt.Errorf("%s or %s is required", inlineEnvVar, pathEnvVar)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read %s: %w", pathEnvVar, err)
	}
	return data, nil
}
