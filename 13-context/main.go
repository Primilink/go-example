package main

import (
	"context"
	"errors"
	"fmt"
	"time"
)

func main() {
	// ========================================
	// CONTEXT - cancellation, deadlines, values
	// ========================================

	// Context is THE way to:
	// 1. Cancel operations
	// 2. Set deadlines/timeouts
	// 3. Pass request-scoped values

	// ========================================
	// CONTEXT.BACKGROUND & TODO
	// ========================================

	// Background: root context, never cancelled
	ctx := context.Background()
	fmt.Println("Background context:", ctx)

	// TODO: placeholder when you're not sure what context to use
	_ = context.TODO()

	// ========================================
	// CANCEL CONTEXT
	// ========================================

	fmt.Println("\n--- Cancel Context ---")

	// Create cancellable context
	ctx, cancel := context.WithCancel(context.Background())

	go func(ctx context.Context) {
		for {
			select {
			case <-ctx.Done():
				fmt.Println("Worker: cancelled!", ctx.Err())
				return
			default:
				fmt.Println("Worker: working...")
				time.Sleep(50 * time.Millisecond)
			}
		}
	}(ctx)

	time.Sleep(150 * time.Millisecond)
	cancel() // signal cancellation
	time.Sleep(50 * time.Millisecond)

	// ========================================
	// TIMEOUT CONTEXT
	// ========================================

	fmt.Println("\n--- Timeout Context ---")

	// Automatically cancels after duration
	ctx2, cancel2 := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel2() // always call cancel to release resources!

	result, err := slowOperation(ctx2)
	if err != nil {
		fmt.Println("Operation failed:", err)
	} else {
		fmt.Println("Result:", result)
	}

	// ========================================
	// DEADLINE CONTEXT
	// ========================================

	fmt.Println("\n--- Deadline Context ---")

	deadline := time.Now().Add(100 * time.Millisecond)
	ctx3, cancel3 := context.WithDeadline(context.Background(), deadline)
	defer cancel3()

	// Check deadline
	if d, ok := ctx3.Deadline(); ok {
		fmt.Println("Deadline set to:", d)
	}

	// ========================================
	// CONTEXT WITH VALUES
	// ========================================

	fmt.Println("\n--- Context Values ---")

	// Pass request-scoped data
	ctx4 := context.WithValue(context.Background(), "requestID", "abc-123")
	ctx4 = context.WithValue(ctx4, "userID", 42)

	processRequest(ctx4)

	// ========================================
	// REAL WORLD EXAMPLE: HTTP-like request
	// ========================================

	fmt.Println("\n--- Real World Example ---")

	// Simulate an HTTP request with timeout
	ctx5, cancel5 := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel5()

	// Add request metadata
	ctx5 = context.WithValue(ctx5, "traceID", "trace-xyz")

	// Make concurrent calls, any can cancel
	results := make(chan string, 2)
	errs := make(chan error, 2)

	go func() {
		result, err := fetchFromDB(ctx5)
		if err != nil {
			errs <- err
			return
		}
		results <- result
	}()

	go func() {
		result, err := fetchFromCache(ctx5)
		if err != nil {
			errs <- err
			return
		}
		results <- result
	}()

	// Wait for first result or context timeout
	select {
	case result := <-results:
		fmt.Println("Got result:", result)
	case err := <-errs:
		fmt.Println("Got error:", err)
	case <-ctx5.Done():
		fmt.Println("Request timed out:", ctx5.Err())
	}

	// ========================================
	// CONTEXT BEST PRACTICES
	// ========================================

	// 1. Always pass context as first parameter
	//    func DoSomething(ctx context.Context, arg1 string) error

	// 2. Don't store context in structs
	//    BAD:  type Server struct { ctx context.Context }
	//    GOOD: pass context to each method

	// 3. Use context.Background() at the top level
	//    Use context.TODO() as a placeholder

	// 4. Always call cancel() even if context times out
	//    defer cancel()

	// 5. Don't pass nil context
	//    If unsure, use context.TODO()

	// 6. Use values sparingly - only for request-scoped data
	//    NOT for passing function parameters!
	//    Good: requestID, userID, traceID
	//    Bad: database connections, config

	// ========================================
	// CONTEXT IN THE STANDARD LIBRARY
	// ========================================

	// Many stdlib functions accept context:
	// - http.NewRequestWithContext()
	// - sql.DB.QueryContext()
	// - exec.CommandContext()
	// - net.Dialer.DialContext()

	fmt.Println("\nContext is essential for production Go!")
}

// ========================================
// HELPER FUNCTIONS
// ========================================

func slowOperation(ctx context.Context) (string, error) {
	select {
	case <-time.After(200 * time.Millisecond): // simulates slow work
		return "completed!", nil
	case <-ctx.Done():
		return "", ctx.Err()
	}
}

func processRequest(ctx context.Context) {
	// Extract values (need type assertion)
	requestID := ctx.Value("requestID").(string)
	userID := ctx.Value("userID").(int)

	fmt.Printf("Processing request %s for user %d\n", requestID, userID)
}

func fetchFromDB(ctx context.Context) (string, error) {
	select {
	case <-time.After(100 * time.Millisecond):
		return "data from DB", nil
	case <-ctx.Done():
		return "", ctx.Err()
	}
}

func fetchFromCache(ctx context.Context) (string, error) {
	select {
	case <-time.After(30 * time.Millisecond):
		// Simulate cache miss sometimes
		return "", errors.New("cache miss")
	case <-ctx.Done():
		return "", ctx.Err()
	}
}
