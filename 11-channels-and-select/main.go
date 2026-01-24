package main

import (
	"fmt"
	"time"
)

func main() {
	// ========================================
	// CHANNELS - typed pipes for communication
	// ========================================

	// Create an unbuffered channel
	ch := make(chan string)

	// Send to channel (in goroutine because unbuffered blocks)
	go func() {
		ch <- "Hello from goroutine!" // send
	}()

	msg := <-ch // receive (blocks until message arrives)
	fmt.Println("Received:", msg)

	// ========================================
	// BUFFERED CHANNELS
	// ========================================

	// Unbuffered: sender blocks until receiver ready
	// Buffered: sender blocks only when buffer full

	buffered := make(chan int, 3) // buffer size 3

	buffered <- 1 // doesn't block
	buffered <- 2 // doesn't block
	buffered <- 3 // doesn't block
	// buffered <- 4 // would block! buffer full

	fmt.Println("Buffered:", <-buffered, <-buffered, <-buffered)

	// ========================================
	// CHANNEL DIRECTIONS (for type safety)
	// ========================================

	ping := make(chan string)
	pong := make(chan string)

	go pingFunc(ping)            // can only send
	go pongFunc(ping, pong)      // ping receive-only, pong send-only
	fmt.Println("Pong:", <-pong) // main receives

	// ========================================
	// CLOSING CHANNELS
	// ========================================

	jobs := make(chan int, 5)

	// Send some jobs
	go func() {
		for i := 1; i <= 3; i++ {
			jobs <- i
		}
		close(jobs) // signal no more jobs
	}()

	// Receive until closed
	for job := range jobs { // automatically stops when closed
		fmt.Println("Processing job:", job)
	}

	// Check if channel closed
	ch2 := make(chan int, 1)
	ch2 <- 42
	close(ch2)

	val, ok := <-ch2
	fmt.Printf("Value: %d, Channel open: %v\n", val, ok) // ok = false after close

	// ========================================
	// SELECT - multiplexing channels
	// ========================================

	fmt.Println("\n--- Select Example ---")

	ch1 := make(chan string)
	ch2Again := make(chan string)

	go func() {
		time.Sleep(50 * time.Millisecond)
		ch1 <- "message from ch1"
	}()

	go func() {
		time.Sleep(30 * time.Millisecond)
		ch2Again <- "message from ch2"
	}()

	// Select waits on multiple channels
	for i := 0; i < 2; i++ {
		select {
		case msg1 := <-ch1:
			fmt.Println("Received:", msg1)
		case msg2 := <-ch2Again:
			fmt.Println("Received:", msg2)
		}
	}

	// ========================================
	// SELECT WITH DEFAULT (non-blocking)
	// ========================================

	empty := make(chan int)

	select {
	case val := <-empty:
		fmt.Println("Got:", val)
	default:
		fmt.Println("No value ready, moving on!")
	}

	// ========================================
	// SELECT WITH TIMEOUT
	// ========================================

	fmt.Println("\n--- Timeout Example ---")

	slowCh := make(chan string)

	go func() {
		time.Sleep(200 * time.Millisecond)
		slowCh <- "slow response"
	}()

	select {
	case msg := <-slowCh:
		fmt.Println("Got:", msg)
	case <-time.After(100 * time.Millisecond):
		fmt.Println("Timeout! Didn't wait forever")
	}

	// ========================================
	// WORKER POOL PATTERN - THE Go pattern!
	// ========================================

	fmt.Println("\n--- Worker Pool ---")

	numWorkers := 3
	jobsCh := make(chan int, 10)
	results := make(chan int, 10)

	// Start workers
	for w := 1; w <= numWorkers; w++ {
		go worker(w, jobsCh, results)
	}

	// Send jobs
	for j := 1; j <= 5; j++ {
		jobsCh <- j
	}
	close(jobsCh) // no more jobs

	// Collect results
	for r := 1; r <= 5; r++ {
		result := <-results
		fmt.Println("Result:", result)
	}

	// ========================================
	// FAN-OUT, FAN-IN PATTERN
	// ========================================

	fmt.Println("\n--- Fan-Out/Fan-In ---")

	input := make(chan int)
	output := fanIn(
		square(input),
		square(input),
		square(input),
	)

	go func() {
		for i := 1; i <= 5; i++ {
			input <- i
		}
		close(input)
	}()

	for result := range output {
		fmt.Println("Squared:", result)
	}
}

// ========================================
// HELPER FUNCTIONS
// ========================================

// Send-only channel parameter
func pingFunc(ch chan<- string) {
	ch <- "ping!"
}

// Receive from one, send to another
func pongFunc(in <-chan string, out chan<- string) {
	msg := <-in
	out <- msg + " pong!"
}

// Worker for pool pattern
func worker(id int, jobs <-chan int, results chan<- int) {
	for job := range jobs {
		fmt.Printf("Worker %d processing job %d\n", id, job)
		time.Sleep(20 * time.Millisecond) // simulate work
		results <- job * 2
	}
}

// Fan-out: distribute work
func square(in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		for n := range in {
			out <- n * n
		}
		close(out)
	}()
	return out
}

// Fan-in: merge multiple channels into one
func fanIn(channels ...<-chan int) <-chan int {
	out := make(chan int)
	var count int

	for _, ch := range channels {
		count++
		go func(c <-chan int) {
			for n := range c {
				out <- n
			}
			count--
			if count == 0 {
				close(out)
			}
		}(ch)
	}

	return out
}
