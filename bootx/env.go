package bootx

import (
	"fmt"
	"os"
	"strings"

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
