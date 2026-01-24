package main

import (
	"fmt"
	"runtime"
	"sync"
	"time"
)

func main() {
	fmt.Println("CPUs available:", runtime.NumCPU())
	fmt.Println("Goroutines at start:", runtime.NumGoroutine())

	// ========================================
	// GOROUTINES - lightweight threads
	// ========================================

	// Just add "go" before a function call!
	go sayHello("Goroutine 1")
	go sayHello("Goroutine 2")
	go sayHello("Goroutine 3")

	// Problem: main() might exit before goroutines finish
	time.Sleep(100 * time.Millisecond) // bad solution, just for demo

	fmt.Println("Goroutines now:", runtime.NumGoroutine())

	// ========================================
	// WAITGROUP - wait for goroutines to finish
	// ========================================

	fmt.Println("\n--- WaitGroup Example ---")

	var wg sync.WaitGroup

	for i := 1; i <= 3; i++ {
		wg.Add(1) // increment counter

		go func(id int) {
			defer wg.Done() // decrement when done
			fmt.Printf("Worker %d starting\n", id)
			time.Sleep(50 * time.Millisecond)
			fmt.Printf("Worker %d done\n", id)
		}(i) // pass i as argument to avoid closure gotcha
	}

	wg.Wait() // blocks until counter is 0
	fmt.Println("All workers completed!")

	// ========================================
	// CLOSURE GOTCHA - classic bug!
	// ========================================

	fmt.Println("\n--- Closure Gotcha ---")

	// WRONG - all goroutines might see the same value
	// for i := 1; i <= 3; i++ {
	//     go func() {
	//         fmt.Println("Wrong:", i) // might print 4, 4, 4
	//     }()
	// }

	// RIGHT - pass as argument
	var wg2 sync.WaitGroup
	for i := 1; i <= 3; i++ {
		wg2.Add(1)
		go func(n int) {
			defer wg2.Done()
			fmt.Println("Right:", n) // prints 1, 2, 3 (order varies)
		}(i)
	}
	wg2.Wait()

	// ========================================
	// GOROUTINES ARE CHEAP!
	// ========================================

	fmt.Println("\n--- 1000 Goroutines ---")

	var wg3 sync.WaitGroup
	start := time.Now()

	for i := 0; i < 1000; i++ {
		wg3.Add(1)
		go func(id int) {
			defer wg3.Done()
			// simulate work
			time.Sleep(10 * time.Millisecond)
		}(i)
	}

	wg3.Wait()
	fmt.Printf("1000 goroutines completed in %v\n", time.Since(start))
	// Note: they run concurrently, so total time ≈ 10ms, not 10000ms!

	// ========================================
	// ANONYMOUS GOROUTINES
	// ========================================

	fmt.Println("\n--- Anonymous Goroutine ---")

	done := make(chan bool) // we'll learn channels next!

	go func() {
		fmt.Println("Anonymous goroutine running!")
		done <- true
	}()

	<-done // wait for signal

	// ========================================
	// GOROUTINE vs THREAD
	// ========================================

	// Threads (OS):
	// - Heavy (~1MB stack)
	// - OS scheduled
	// - Expensive to create

	// Goroutines (Go runtime):
	// - Lightweight (~2KB stack, grows as needed)
	// - Go scheduler (M:N scheduling)
	// - Can run millions of them!

	// This is why Go is AMAZING for high-concurrency:
	// 20k req/s? Just spawn a goroutine per request. No problem.

	// ========================================
	// GOMAXPROCS
	// ========================================

	// How many OS threads can run goroutines simultaneously
	fmt.Println("\nGOMAXPROCS:", runtime.GOMAXPROCS(0)) // 0 = query current
	// Defaults to number of CPUs
	// runtime.GOMAXPROCS(4) // set to 4

	fmt.Println("\nNext lesson: Channels - how goroutines communicate!")
}

func sayHello(name string) {
	fmt.Println("Hello from", name)
}
