//go:build ignore

package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
)

func main() {
	// Получаем текущую рабочую директорию
	wd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	// Ищем корень проекта
	projectRoot, err := findProjectRoot(wd)
	if err != nil {
		log.Fatal(err)
	}

	// Загружаем .env
	envPath := filepath.Join(projectRoot, ".env")
	if _, err := os.Stat(envPath); err == nil {
		if err := godotenv.Load(envPath); err != nil {
			log.Printf("Warning: failed to load .env file at %s: %v", envPath, err)
		}
	} else {
		log.Printf("Warning: .env file not found at %s, using default values", envPath)
	}

	// Путь к миграциям
	migrationsPath := filepath.Join(projectRoot, "migrations")

	// Проверяем существование папки миграций
	if _, err := os.Stat(migrationsPath); os.IsNotExist(err) {
		log.Fatalf("Migrations directory does not exist: %s", migrationsPath)
	}

	// Формируем URL базы данных
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		// Build DATABASE_URL из .env или переменных окружения
		host := getEnv("DB_HOST", "localhost")
		port := getEnv("DB_PORT", "5432")
		user := getEnv("DB_USER", "postgres")
		password := getEnv("DB_PASSWORD", "postgres")
		name := getEnv("DB_NAME", "movienight")
		sslmode := getEnv("DB_SSLMODE", "disable")

		databaseURL = fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
			user, password, host, port, name, sslmode)
	}

	log.Printf("Project root: %s", projectRoot)
	log.Printf("Migrations path: %s", migrationsPath)
	log.Printf("Database URL: %s", databaseURL)

	// Подключаемся к БД
	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	goose.SetDialect("postgres")

	if len(os.Args) < 2 {
		log.Fatal("Usage: migrate-cli <up|down|version|force>")
	}

	switch os.Args[1] {
	case "up":
		if err := goose.Up(db, migrationsPath); err != nil {
			log.Fatal(err)
		}
		fmt.Println("Migrations applied successfully")
	case "down":
		steps := 1
		if len(os.Args) > 2 {
			// Если указано количество шагов
			fmt.Sscanf(os.Args[2], "%d", &steps)
		}
		successfulSteps := 0
		for i := 0; i < steps; i++ {
			if err := goose.Down(db, migrationsPath); err != nil {
				break
			}
			successfulSteps++
		}
		fmt.Printf("Rolled back %d migration(s)\n", successfulSteps)
	case "version":
		version, err := goose.GetDBVersion(db)
		if err != nil {
			fmt.Println("No migrations have been applied yet")
		} else {
			fmt.Printf("Current version: %d\n", version)
		}
	case "force":
		if len(os.Args) != 3 {
			log.Fatal("Usage: migrate-cli force <version>")
		}
		var version int64
		fmt.Sscanf(os.Args[2], "%d", &version)
		if err := goose.DownTo(db, migrationsPath, version); err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Forced version to %d\n", version)
	default:
		log.Fatal("Unknown command. Use 'up', 'down', 'version', or 'force'")
	}
}

// findProjectRoot ищет корень проекта, начиная с currentDir и поднимаясь вверх,
// пока не найдет директорию, содержащую папку migrations
func findProjectRoot(currentDir string) (string, error) {
	dir := currentDir

	for {
		// Проверяем наличие папки migrations в текущей директории
		migrationsPath := filepath.Join(dir, "migrations")
		if info, err := os.Stat(migrationsPath); err == nil && info.IsDir() {
			// Проверяем, что в папке есть хотя бы один .sql файл
			entries, _ := os.ReadDir(migrationsPath)
			for _, entry := range entries {
				if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".sql") {
					return dir, nil
				}
			}
			// Если нет .sql файлов, но папка существует - тоже считаем корнем
			return dir, nil
		}

		// Поднимаемся на уровень выше
		parent := filepath.Dir(dir)
		if parent == dir {
			// Дошли до корня файловой системы
			break
		}
		dir = parent
	}

	return "", fmt.Errorf("could not find 'migrations' folder in current or parent directories (started from: %s)", currentDir)
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
