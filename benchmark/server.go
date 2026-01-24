package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"sync"
	"time"
)

type Result struct {
	Number  int   `json:"number"`
	Square  int   `json:"square"`
	TimeNs  int64 `json:"time_ns"`
	Workers int   `json:"workers_used"`
}

func main() {
	http.HandleFunc("/square", squareHandler)
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	fmt.Println("🔥 Server running on http://localhost:8080")
	fmt.Println("   GET /square - returns random square calculation")
	fmt.Println("   GET /health - health check")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func squareHandler(w http.ResponseWriter, r *http.Request) {
	start := time.Now()

	// Channel to collect results
	results := make(chan Result, 101)
	var wg sync.WaitGroup

	// Spawn goroutines to calculate squares concurrently
	numWorkers := 10
	jobs := make(chan int, 101)

	// Start workers
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for n := range jobs {
				results <- Result{
					Number:  n,
					Square:  n * n,
					Workers: numWorkers,
				}
			}
		}()
	}

	// Send jobs
	go func() {
		for i := 0; i <= 100; i++ {
			jobs <- i
		}
		close(jobs)
	}()

	// Wait and close results
	go func() {
		wg.Wait()
		close(results)
	}()

	// Collect all results
	allResults := make([]Result, 0, 101)
	for res := range results {
		allResults = append(allResults, res)
	}

	// Pick random result
	randomResult := allResults[rand.Intn(len(allResults))]
	randomResult.TimeNs = time.Since(start).Nanoseconds()

	// Return JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(randomResult)
}
