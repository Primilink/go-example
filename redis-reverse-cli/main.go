package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/redis/go-redis/v9"
)

const keyPrefix = "reverse:"

var ctx = context.Background()

func main() {
	// Connect to Redis
	redisURL := os.Getenv("REDIS_URL")
	if redisURL == "" {
		redisURL = "redis://localhost:6379"
	}

	opt, err := redis.ParseURL(redisURL)
	if err != nil {
		fmt.Printf("Error parsing Redis URL: %v\n", err)
		os.Exit(1)
	}

	client := redis.NewClient(opt)

	// Test connection
	if err := client.Ping(ctx).Err(); err != nil {
		fmt.Printf("Error connecting to Redis: %v\n", err)
		fmt.Println("Make sure Redis is running: docker run -d -p 6379:6379 redis")
		os.Exit(1)
	}

	fmt.Println("=== Redis Reverse String CLI ===")
	fmt.Println("• Type text to store (key=input, value=reversed)")
	fmt.Println("• Type existing key to delete it")
	fmt.Println("• Ctrl+C to exit and clear all keys")
	fmt.Println()

	// Handle Ctrl+C - cleanup all keys
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigChan
		fmt.Println("\n\nCleaning up keys...")
		cleanupKeys(client)
		fmt.Println("Goodbye!")
		os.Exit(0)
	}()

	scanner := bufio.NewScanner(os.Stdin)

	// Show initial table
	printTable(client)

	for {
		fmt.Print("\nInput: ")
		if !scanner.Scan() {
			break
		}

		input := strings.TrimSpace(scanner.Text())
		if input == "" {
			continue
		}

		key := keyPrefix + input

		// Check if key exists
		exists, err := client.Exists(ctx, key).Result()
		if err != nil {
			fmt.Printf("Error checking key: %v\n", err)
			continue
		}

		if exists > 0 {
			// Key exists - delete it
			if err := client.Del(ctx, key).Err(); err != nil {
				fmt.Printf("Error deleting key: %v\n", err)
				continue
			}
			fmt.Printf("Deleted: %s\n", input)
		} else {
			// New key - store reversed string
			reversed := reverseString(input)
			if err := client.Set(ctx, key, reversed, 0).Err(); err != nil {
				fmt.Printf("Error setting key: %v\n", err)
				continue
			}
			fmt.Printf("Stored: %s → %s\n", input, reversed)
		}

		printTable(client)
	}
}

func reverseString(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

func printTable(client *redis.Client) {
	keys, err := client.Keys(ctx, keyPrefix+"*").Result()
	if err != nil {
		fmt.Printf("Error getting keys: %v\n", err)
		return
	}

	if len(keys) == 0 {
		fmt.Println("\n(no entries)")
		return
	}

	// Get all values
	type entry struct {
		key   string
		value string
	}
	entries := make([]entry, 0, len(keys))

	maxKeyLen := 3 // minimum "Key" header
	maxValLen := 5 // minimum "Value" header

	for _, k := range keys {
		val, err := client.Get(ctx, k).Result()
		if err != nil {
			continue
		}

		displayKey := strings.TrimPrefix(k, keyPrefix)
		entries = append(entries, entry{key: displayKey, value: val})

		if len(displayKey) > maxKeyLen {
			maxKeyLen = len(displayKey)
		}
		if len(val) > maxValLen {
			maxValLen = len(val)
		}
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

func cleanupKeys(client *redis.Client) {
	keys, err := client.Keys(ctx, keyPrefix+"*").Result()
	if err != nil {
		fmt.Printf("Error getting keys: %v\n", err)
		return
	}

	if len(keys) == 0 {
		return
	}

	deleted, err := client.Del(ctx, keys...).Result()
	if err != nil {
		fmt.Printf("Error deleting keys: %v\n", err)
		return
	}

	fmt.Printf("Deleted %d keys\n", deleted)
}
