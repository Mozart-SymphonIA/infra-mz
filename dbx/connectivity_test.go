package dbx_test

import (
	"context"
	"github.com/Mozart-SymphonIA/infra-mz/dbx/factory"
	"os"
	"testing"
	"time"
)

func TestPostgresConnectivity(t *testing.T) {
	// Read connection settings from env or use defaults relevant for local/test
	if os.Getenv("TEST_INTEGRATION") != "true" {
		t.Skip("Skipping integration test; set TEST_INTEGRATION=true to run")
	}

	dbName := os.Getenv("POSTGRES_DB")
	if dbName == "" {
		dbName = "postgres" // Default often used in docker
	}

	// Initialize bundle
	bundle, err := factory.NewBundleFromEnv(dbName)
	if err != nil {
		t.Fatalf("Failed to create bundle: %v", err)
	}
	defer bundle.Conn.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Check Health (Ping)
	if err := bundle.Inspector.Ping(ctx); err != nil {
		t.Fatalf("Failed to ping database: %v", err)
	}

	// Run a simple query
	config, err := bundle.Reader.Query(ctx, "SELECT 1")
	if err != nil {
		t.Fatalf("Failed to execute query: %v", err)
	}

	// Assuming the SQL reader returns the result as string for this generic Query method 
	// (based on previous exploration of sqlReader.Query which scans into a string)
	if config != "1" {
		t.Errorf("Expected '1', got '%s'", config)
	}
}
