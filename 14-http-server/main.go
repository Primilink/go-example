package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	// ========================================
	// BASIC HTTP SERVER
	// ========================================

	// Route handlers
	// Using a custom mux (recommended for production)
	mux := http.NewServeMux()

	mux.HandleFunc("/", homeHandler)
	mux.HandleFunc("/hello", helloHandler)
	mux.HandleFunc("/json", jsonHandler)
	mux.HandleFunc("/user", userHandler) // POST/GET example

	mux.HandleFunc("/api/health", healthHandler)
	mux.HandleFunc("/api/data", dataHandler)

	// ========================================
	// MIDDLEWARE
	// ========================================

	// Wrap handlers with middleware
	mux.Handle("/api/protected", authMiddleware(http.HandlerFunc(protectedHandler)))

	// Chain multiple middlewares
	mux.Handle("/api/logged",
		loggingMiddleware(
			authMiddleware(
				http.HandlerFunc(protectedHandler),
			),
		),
	)

	// ========================================
	// CUSTOM SERVER (recommended for production)
	// ========================================

	server := &http.Server{
		Addr:         ":8080",
		Handler:      mux,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// ========================================
	// GRACEFUL SHUTDOWN
	// ========================================

	// Channel to listen for interrupt
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGTERM)

	// Start server in goroutine
	go func() {
		fmt.Println("Server starting on http://localhost:8080")
		fmt.Println("Endpoints:")
		fmt.Println("  GET  /api/health")
		fmt.Println("  GET  /api/data")
		fmt.Println("  GET  /api/protected (needs Authorization header)")
		fmt.Println("\nPress Ctrl+C to stop")

		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	// Wait for interrupt
	<-done
	fmt.Println("\nShutting down gracefully...")

	// Give active connections 30s to finish
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Shutdown error: %v", err)
	}

	fmt.Println("Server stopped")
}

// ========================================
// HANDLERS
// ========================================

func homeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to Go HTTP server!")
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	if name == "" {
		name = "World"
	}
	fmt.Fprintf(w, "Hello, %s!", name)
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status": "healthy",
		"time":   time.Now().Format(time.RFC3339),
	})
}

func jsonHandler(w http.ResponseWriter, r *http.Request) {
	response := map[string]any{
		"message": "Hello, JSON!",
		"count":   42,
		"nested": map[string]string{
			"key": "value",
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)

}

func dataHandler(w http.ResponseWriter, r *http.Request) {
	data := []map[string]any{
		{"id": 1, "name": "Item 1"},
		{"id": 2, "name": "Item 2"},
		{"id": 3, "name": "Item 3"},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

// ========================================
// HANDLING DIFFERENT METHODS
// ========================================

type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func userHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		// Return a user
		user := User{ID: 1, Name: "Primi", Email: "primi@example.com"}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(user)

	case http.MethodPost:
		// Create a user from JSON body
		var user User
		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		// In real app: save to database
		user.ID = 1 // assign ID

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(user)

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func protectedHandler(w http.ResponseWriter, r *http.Request) {
	// Get user from context (set by middleware)
	userID := r.Context().Value("userID")
	fmt.Fprintf(w, "Protected resource! User: %v", userID)
}

// ========================================
// MIDDLEWARE
// ========================================

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Call the next handler
		next.ServeHTTP(w, r)

		// Log after
		log.Printf("%s %s %v", r.Method, r.URL.Path, time.Since(start))
	})
}

func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")

		if token == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// In real app: validate token, get user
		userID := "user-123"

		// Add to context for downstream handlers
		ctx := context.WithValue(r.Context(), "userID", userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// ========================================
// TIPS FOR PRODUCTION
// ========================================

// 1. Use a router library for complex routing:
//    - chi (lightweight, stdlib compatible)
//    - gin (fast, popular)
//    - echo (similar to gin)

// 2. Always set timeouts on your server

// 3. Implement graceful shutdown

// 4. Use middleware for:
//    - Logging
//    - Authentication
//    - CORS
//    - Rate limiting
//    - Request ID tracing

// 5. Use json struct tags:
//    `json:"name,omitempty"` // omit if empty
//    `json:"-"`              // never include

// 6. Validate input! Don't trust the client
