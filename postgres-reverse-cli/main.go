package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/jackc/pgx/v5"
)

var ctx = context.Background()

func main() {
	// Connect to PostgreSQL
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable"
	}

	conn, err := pgx.Connect(ctx, dbURL)
	if err != nil {
		fmt.Printf("Error connecting to PostgreSQL: %v\n", err)
		fmt.Println("Make sure PostgreSQL is running:")
		fmt.Println("  docker run -d -p 5432:5432 -e POSTGRES_PASSWORD=postgres postgres")
		os.Exit(1)
	}
	defer conn.Close(ctx)

	// Create table if not exists
	if err := setupTable(conn); err != nil {
		fmt.Printf("Error creating table: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("=== PostgreSQL Reverse String CLI ===")
	fmt.Println("• Type text to store (key=input, value=reversed)")
	fmt.Println("• Type existing key to delete it")
	fmt.Println("• Ctrl+C to exit and clear all entries")
	fmt.Println()

	// Handle Ctrl+C - cleanup all entries
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigChan
		fmt.Println("\n\nCleaning up entries...")
		cleanupEntries(conn)
		fmt.Println("Goodbye!")
		os.Exit(0)
	}()

	scanner := bufio.NewScanner(os.Stdin)

	// Show initial table
	printTable(conn)

	for {
		fmt.Print("\nInput: ")
		if !scanner.Scan() {
			break
		}

		input := strings.TrimSpace(scanner.Text())
		if input == "" {
			continue
		}

		// Check if key exists
		exists, err := keyExists(conn, input)
		if err != nil {
			fmt.Printf("Error checking key: %v\n", err)
			continue
		}

		if exists {
			// Key exists - delete it
			if err := deleteEntry(conn, input); err != nil {
				fmt.Printf("Error deleting entry: %v\n", err)
				continue
			}
			fmt.Printf("Deleted: %s\n", input)
		} else {
			// New key - store reversed string
			reversed := reverseString(input)
			if err := insertEntry(conn, input, reversed); err != nil {
				fmt.Printf("Error inserting entry: %v\n", err)
				continue
			}
			fmt.Printf("Stored: %s → %s\n", input, reversed)
		}

		printTable(conn)
	}
}

func setupTable(conn *pgx.Conn) error {
	_, err := conn.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS reverse_strings (
			key TEXT PRIMARY KEY,
			value TEXT NOT NULL,
			created_at TIMESTAMP DEFAULT NOW()
		)
	`)
	return err
}

func reverseString(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

func keyExists(conn *pgx.Conn, key string) (bool, error) {
	var exists bool
	err := conn.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM reverse_strings WHERE key = $1)", key).Scan(&exists)
	return exists, err
}

func insertEntry(conn *pgx.Conn, key, value string) error {
	_, err := conn.Exec(ctx, "INSERT INTO reverse_strings (key, value) VALUES ($1, $2)", key, value)
	return err
}

func deleteEntry(conn *pgx.Conn, key string) error {
	_, err := conn.Exec(ctx, "DELETE FROM reverse_strings WHERE key = $1", key)
	return err
}

func printTable(conn *pgx.Conn) {
	rows, err := conn.Query(ctx, "SELECT key, value FROM reverse_strings ORDER BY created_at")
	if err != nil {
		fmt.Printf("Error querying entries: %v\n", err)
		return
	}
	defer rows.Close()

	type entry struct {
		key   string
		value string
	}
	entries := []entry{}

	maxKeyLen := 3 // minimum "Key" header
	maxValLen := 5 // minimum "Value" header

	for rows.Next() {
		var e entry
		if err := rows.Scan(&e.key, &e.value); err != nil {
			continue
		}
		entries = append(entries, e)

		if len(e.key) > maxKeyLen {
			maxKeyLen = len(e.key)
		}
		if len(e.value) > maxValLen {
			maxValLen = len(e.value)
		}
	}

	if len(entries) == 0 {
		fmt.Println("\n(no entries)")
		return
	}

	// Print table
	fmt.Println()
	divider := "+" + strings.Repeat("-", maxKeyLen+2) + "+" + strings.Repeat("-", maxValLen+2) + "+"

	fmt.Println(divider)
	fmt.Printf("| %-*s | %-*s |\n", maxKeyLen, "Key", maxValLen, "Value")
	fmt.Println(divider)

	for _, e := range entries {
		fmt.Printf("| %-*s | %-*s |\n", maxKeyLen, e.key, maxValLen, e.value)
	}

	fmt.Println(divider)
	fmt.Printf("Total: %d entries\n", len(entries))
}

func cleanupEntries(conn *pgx.Conn) {
	result, err := conn.Exec(ctx, "DELETE FROM reverse_strings")
	if err != nil {
		fmt.Printf("Error deleting entries: %v\n", err)
		return
	}

	fmt.Printf("Deleted %d entries\n", result.RowsAffected())
}
