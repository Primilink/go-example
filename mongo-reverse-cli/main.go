package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var ctx = context.Background()

type Entry struct {
	Key       string    `bson:"_id"`
	Value     string    `bson:"value"`
	CreatedAt time.Time `bson:"created_at"`
}

func main() {
	// Connect to MongoDB
	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		mongoURI = "mongodb://localhost:27017"
	}

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		fmt.Printf("Error connecting to MongoDB: %v\n", err)
		os.Exit(1)
	}
	defer client.Disconnect(ctx)

	// Test connection
	if err := client.Ping(ctx, nil); err != nil {
		fmt.Printf("Error pinging MongoDB: %v\n", err)
		fmt.Println("Make sure MongoDB is running:")
		fmt.Println("  docker run -d -p 27017:27017 mongo")
		os.Exit(1)
	}

	collection := client.Database("reverse_db").Collection("strings")

	fmt.Println("=== MongoDB Reverse String CLI ===")
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
		cleanupEntries(collection)
		fmt.Println("Goodbye!")
		os.Exit(0)
	}()

	scanner := bufio.NewScanner(os.Stdin)

	// Show initial table
	printTable(collection)

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
		exists, err := keyExists(collection, input)
		if err != nil {
			fmt.Printf("Error checking key: %v\n", err)
			continue
		}

		if exists {
			// Key exists - delete it
			if err := deleteEntry(collection, input); err != nil {
				fmt.Printf("Error deleting entry: %v\n", err)
				continue
			}
			fmt.Printf("Deleted: %s\n", input)
		} else {
			// New key - store reversed string
			reversed := reverseString(input)
			if err := insertEntry(collection, input, reversed); err != nil {
				fmt.Printf("Error inserting entry: %v\n", err)
				continue
			}
			fmt.Printf("Stored: %s → %s\n", input, reversed)
		}

		printTable(collection)
	}
}

func reverseString(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

func keyExists(collection *mongo.Collection, key string) (bool, error) {
	count, err := collection.CountDocuments(ctx, bson.M{"_id": key})
	return count > 0, err
}

func insertEntry(collection *mongo.Collection, key, value string) error {
	entry := Entry{
		Key:       key,
		Value:     value,
		CreatedAt: time.Now(),
	}
	_, err := collection.InsertOne(ctx, entry)
	return err
}

func deleteEntry(collection *mongo.Collection, key string) error {
	_, err := collection.DeleteOne(ctx, bson.M{"_id": key})
	return err
}

func printTable(collection *mongo.Collection) {
	cursor, err := collection.Find(ctx, bson.M{}, options.Find().SetSort(bson.M{"created_at": 1}))
	if err != nil {
		fmt.Printf("Error querying entries: %v\n", err)
		return
	}
	defer cursor.Close(ctx)

	var entries []Entry
	if err := cursor.All(ctx, &entries); err != nil {
		fmt.Printf("Error decoding entries: %v\n", err)
		return
	}

	if len(entries) == 0 {
		fmt.Println("\n(no entries)")
		return
	}

	maxKeyLen := 3 // minimum "Key" header
	maxValLen := 5 // minimum "Value" header

	for _, e := range entries {
		if len(e.Key) > maxKeyLen {
			maxKeyLen = len(e.Key)
		}
		if len(e.Value) > maxValLen {
			maxValLen = len(e.Value)
		}
	}

	// Print table
	fmt.Println()
	divider := "+" + strings.Repeat("-", maxKeyLen+2) + "+" + strings.Repeat("-", maxValLen+2) + "+"

	fmt.Println(divider)
	fmt.Printf("| %-*s | %-*s |\n", maxKeyLen, "Key", maxValLen, "Value")
	fmt.Println(divider)

	for _, e := range entries {
		fmt.Printf("| %-*s | %-*s |\n", maxKeyLen, e.Key, maxValLen, e.Value)
	}

	fmt.Println(divider)
	fmt.Printf("Total: %d entries\n", len(entries))
}

func cleanupEntries(collection *mongo.Collection) {
	result, err := collection.DeleteMany(ctx, bson.M{})
	if err != nil {
		fmt.Printf("Error deleting entries: %v\n", err)
		return
	}

	fmt.Printf("Deleted %d entries\n", result.DeletedCount)
}
