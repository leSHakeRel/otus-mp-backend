package config

import (
	"os"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	if err := os.Setenv("DB_HOST", "localhost"); err != nil {
		t.Fatal(err)
	}
	if err := os.Setenv("DB_PORT", "5432"); err != nil {
		t.Fatal(err)
	}
	if err := os.Setenv("DB_USER", "postgres"); err != nil {
		t.Fatal(err)
	}
	if err := os.Setenv("DB_PASSWORD", "postgres"); err != nil {
		t.Fatal(err)
	}
	if err := os.Setenv("DB_NAME", "test_db"); err != nil {
		t.Fatal(err)
	}
	if err := os.Setenv("DB_SSLMODE", "disable"); err != nil {
		t.Fatal(err)
	}
	if err := os.Setenv("JWT_SECRET", "test-secret"); err != nil {
		t.Fatal(err)
	}
	if err := os.Setenv("TMDB_API_KEY", "test-key"); err != nil {
		t.Fatal(err)
	}
	if err := os.Setenv("SERVER_PORT", "8080"); err != nil {
		t.Fatal(err)
	}
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
