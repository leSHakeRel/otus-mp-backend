//go:build ignore

package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

func main() {
	fmt.Print("Enter migration name: ")
	var name string
	fmt.Scanln(&name)

	if name == "" {
		log.Fatal("Migration name cannot be empty")
	}

	// Sanitize name
	name = strings.ReplaceAll(name, " ", "_")
	name = strings.ToLower(name)

	timestamp := time.Now().Format("20060102150405")
	filename := fmt.Sprintf("%s_%s.sql", timestamp, name)

	migrationsPath := "migrations"
	fullPath := filepath.Join(migrationsPath, filename)

	// Create file
	file, err := os.Create(fullPath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// Write template
	file.WriteString("-- +migrate Up\n\n\n-- +migrate Down\n\n")

	fmt.Printf("Created migration: %s\n", filename)

	// Open the file in editor
	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = "notepad"
	}
	cmd := exec.Command(editor, fullPath)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	_ = cmd.Run()
}
