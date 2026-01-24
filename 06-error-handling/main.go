package main

import (
	"errors"
	"fmt"
	"os"
	"strconv"
)

func main() {
	// ========================================
	// THE GO WAY: errors are values, not exceptions
	// ========================================

	// Pattern you'll write 1000x
	result, err := divide(10, 2)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("10 / 2 =", result)

	// Handle the error case
	result2, err := divide(10, 0)
	if err != nil {
		fmt.Println("Error:", err) // prints: Error: cannot divide by zero
	}
	fmt.Println("Result:", result2)

	// ========================================
	// CREATING ERRORS
	// ========================================

	// 1. Simple error with errors.New
	err1 := errors.New("something went wrong")
	fmt.Println("Simple error:", err1)

	// 2. Formatted error with fmt.Errorf
	name := "config.json"
	err2 := fmt.Errorf("file not found: %s", name)
	fmt.Println("Formatted error:", err2)

	// ========================================
	// ERROR WRAPPING (Go 1.13+) - for context
	// ========================================

	err3 := loadConfig("missing.json")
	if err3 != nil {
		fmt.Println("Config error:", err3)
		// Output: Config error: failed to load config: open missing.json: no such file or directory
	}

	// Unwrap to check the original error
	if errors.Is(err3, os.ErrNotExist) {
		fmt.Println("The file doesn't exist!")
	}

	// ========================================
	// CUSTOM ERROR TYPES
	// ========================================

	err4 := processAge(-5)
	if err4 != nil {
		fmt.Println("Age error:", err4)

		// Type assertion to get details
		var validErr *ValidationError
		if errors.As(err4, &validErr) {
			fmt.Printf("  Field: %s, Value: %v\n", validErr.Field, validErr.Value)
		}
	}

	// ========================================
	// SENTINEL ERRORS - predefined errors to check against
	// ========================================

	_, err5 := findUser(999)
	if errors.Is(err5, ErrNotFound) {
		fmt.Println("User not found (sentinel error)")
	}

	// ========================================
	// MULTIPLE ERROR HANDLING PATTERNS
	// ========================================

	// Pattern 1: Early return (most common)
	if err := doStep1(); err != nil {
		fmt.Println("Step 1 failed:", err)
		// return err  // in real code
	}

	// Pattern 2: Inline error handling
	if val, err := strconv.Atoi("42"); err != nil {
		fmt.Println("Parse error:", err)
	} else {
		fmt.Println("Parsed value:", val)
	}

	// Pattern 3: Ignore error intentionally (use _ but be careful!)
	val, _ := strconv.Atoi("123") // only when you're SURE it won't fail
	fmt.Println("Ignored error, got:", val)

	// ========================================
	// PANIC & RECOVER - for truly exceptional cases
	// ========================================

	// panic = unrecoverable error, crashes the program
	// recover = catch a panic (use sparingly!)

	fmt.Println("\n--- Panic & Recover ---")
	safeCall(func() {
		fmt.Println("This runs fine")
	})

	safeCall(func() {
		panic("oh no everything is on fire 🔥")
	})

	fmt.Println("Program continues after recovered panic!")

	// When to use panic:
	// - Programmer error (bug that should never happen)
	// - Initialization that MUST succeed
	// - Never for expected errors like "file not found"!
}

// ========================================
// FUNCTIONS WITH ERROR RETURNS
// ========================================

func divide(a, b int) (int, error) {
	if b == 0 {
		return 0, errors.New("cannot divide by zero")
	}
	return a / b, nil
}

// Error wrapping with %w
func loadConfig(path string) error {
	_, err := os.Open(path)
	if err != nil {
		// %w wraps the original error (preserves the chain)
		return fmt.Errorf("failed to load config: %w", err)
	}
	return nil
}

// ========================================
// CUSTOM ERROR TYPE
// ========================================

type ValidationError struct {
	Field string
	Value any
	Msg   string
}

// Implement the error interface (just needs Error() string)
func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation failed on %s: %s", e.Field, e.Msg)
}

func processAge(age int) error {
	if age < 0 {
		return &ValidationError{
			Field: "age",
			Value: age,
			Msg:   "age cannot be negative",
		}
	}
	return nil
}

// ========================================
// SENTINEL ERRORS
// ========================================

var ErrNotFound = errors.New("not found")
var ErrUnauthorized = errors.New("unauthorized")

func findUser(id int) (string, error) {
	// Simulate user lookup
	if id == 999 {
		return "", ErrNotFound
	}
	return "Primi", nil
}

// ========================================
// HELPER FUNCTIONS
// ========================================

func doStep1() error {
	return nil // success
}

// Recover from panic
func safeCall(fn func()) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered from panic:", r)
		}
	}()
	fn()
}
