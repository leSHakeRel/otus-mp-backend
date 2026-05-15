package config

import (
	"os"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	os.Setenv("DB_HOST", "localhost")
	os.Setenv("DB_PORT", "5432")
	os.Setenv("DB_USER", "postgres")
	os.Setenv("DB_PASSWORD", "postgres")
	os.Setenv("DB_NAME", "test_db")
	os.Setenv("DB_SSLMODE", "disable")
	os.Setenv("JWT_SECRET", "test-secret")
	os.Setenv("TMDB_API_KEY", "test-key")
	os.Setenv("SERVER_PORT", "8080")
	defer os.Clearenv()

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	if cfg.Database.Host != "localhost" {
		t.Errorf("Expected Database.Host to be 'localhost', got '%s'", cfg.Database.Host)
	}

	if cfg.Server.Port != "8080" {
		t.Errorf("Expected Server.Port to be '8080', got '%s'", cfg.Server.Port)
	}
}

func TestDatabaseConfigDSN(t *testing.T) {
	dbConfig := DatabaseConfig{
		Host:     "localhost",
		Port:     "5432",
		User:     "testuser",
		Password: "testpass",
		DBName:   "testdb",
		SSLMode:  "disable",
	}

	expected := "host=localhost port=5432 user=testuser password=testpass dbname=testdb sslmode=disable"
	actual := dbConfig.DSN()

	if actual != expected {
		t.Errorf("Expected DSN to be '%s', got '%s'", expected, actual)
	}
}
