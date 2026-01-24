package main

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

func main() {
	// ========================================
	// MUTEX - mutual exclusion lock
	// ========================================

	fmt.Println("--- Mutex Example ---")

	counter := &Counter{}
	var wg sync.WaitGroup

	// Without mutex, this would have race conditions!
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			counter.Increment()
		}()
	}

	wg.Wait()
	fmt.Println("Counter value:", counter.Value())
	// Always 1000, thanks to mutex!

	// ========================================
	// RWMUTEX - multiple readers, single writer
	// ========================================

	fmt.Println("\n--- RWMutex Example ---")

	cache := &Cache{data: make(map[string]string)}
	var wg2 sync.WaitGroup

	// Multiple writers
	for i := 0; i < 5; i++ {
		wg2.Add(1)
		go func(id int) {
			defer wg2.Done()
			key := fmt.Sprintf("key%d", id)
			cache.Set(key, fmt.Sprintf("value%d", id))
		}(i)
	}

	wg2.Wait()

	// Multiple readers can read concurrently
	for i := 0; i < 5; i++ {
		wg2.Add(1)
		go func(id int) {
			defer wg2.Done()
			key := fmt.Sprintf("key%d", id)
			val := cache.Get(key)
			fmt.Printf("Read %s = %s\n", key, val)
		}(i)
	}

	wg2.Wait()

	// ========================================
	// ATOMIC OPERATIONS - lock-free
	// ========================================

	fmt.Println("\n--- Atomic Example ---")

	var atomicCounter int64
	var wg3 sync.WaitGroup

	for i := 0; i < 1000; i++ {
		wg3.Add(1)
		go func() {
			defer wg3.Done()
			atomic.AddInt64(&atomicCounter, 1)
		}()
	}

	wg3.Wait()
	fmt.Println("Atomic counter:", atomic.LoadInt64(&atomicCounter))

	// Other atomic operations:
	// atomic.StoreInt64(&val, 100)    // set
	// atomic.LoadInt64(&val)          // get
	// atomic.SwapInt64(&val, new)     // set and get old
	// atomic.CompareAndSwapInt64(&val, old, new) // CAS

	// ========================================
	// ONCE - do something exactly once
	// ========================================

	fmt.Println("\n--- Once Example ---")

	var once sync.Once
	var wg4 sync.WaitGroup

	initFunc := func() {
		fmt.Println("Initialization (only runs once!)")
	}

	// Try to init from multiple goroutines
	for i := 0; i < 5; i++ {
		wg4.Add(1)
		go func(id int) {
			defer wg4.Done()
			fmt.Printf("Goroutine %d calling init\n", id)
			once.Do(initFunc) // only first call executes
		}(i)
	}

	wg4.Wait()

	// ========================================
	// COND - condition variable
	// ========================================

	fmt.Println("\n--- Cond Example ---")

	var ready bool
	var mu sync.Mutex
	cond := sync.NewCond(&mu)

	// Waiter goroutine
	go func() {
		mu.Lock()
		for !ready { // always check condition in loop!
			cond.Wait() // releases lock, waits, reacquires lock
		}
		fmt.Println("Worker: condition met, proceeding!")
		mu.Unlock()
	}()

	// Signaler
	time.Sleep(50 * time.Millisecond)
	mu.Lock()
	ready = true
	cond.Signal() // wake one waiter
	// cond.Broadcast() // wake ALL waiters
	mu.Unlock()

	time.Sleep(50 * time.Millisecond)

	// ========================================
	// POOL - reuse expensive objects
	// ========================================

	fmt.Println("\n--- Pool Example ---")

	bufferPool := &sync.Pool{
		New: func() any {
			fmt.Println("Creating new buffer")
			return make([]byte, 1024)
		},
	}

	// First get creates new
	buf1 := bufferPool.Get().([]byte)
	fmt.Println("Got buffer 1, len:", len(buf1))

	// Return to pool
	bufferPool.Put(buf1)

	// Second get reuses
	buf2 := bufferPool.Get().([]byte)
	fmt.Println("Got buffer 2 (reused), len:", len(buf2))

	// Pool is great for:
	// - Byte buffers
	// - Temporary objects
	// - Reducing GC pressure in hot paths

	// ========================================
	// MAP - concurrent map (Go 1.9+)
	// ========================================

	fmt.Println("\n--- sync.Map Example ---")

	var m sync.Map
	var wg5 sync.WaitGroup

	// Concurrent writes - no locks needed!
	for i := 0; i < 5; i++ {
		wg5.Add(1)
		go func(id int) {
			defer wg5.Done()
			m.Store(id, fmt.Sprintf("value%d", id))
		}(i)
	}

	wg5.Wait()

	// Read
	if val, ok := m.Load(3); ok {
		fmt.Println("Key 3:", val)
	}

	// Iterate
	m.Range(func(key, value any) bool {
		fmt.Printf("  %v: %v\n", key, value)
		return true // continue iteration
	})

	// ========================================
	// WHEN TO USE WHAT
	// ========================================

	// Channels:
	// - Passing data/ownership between goroutines
	// - Signaling events
	// - When you want to "share by communicating"

	// Mutex:
	// - Protecting shared state
	// - Simple critical sections
	// - When you need to "communicate by sharing"

	// RWMutex:
	// - Many readers, few writers
	// - Read-heavy workloads

	// Atomic:
	// - Simple counters
	// - Flags
	// - When you need max performance

	// sync.Map:
	// - Concurrent map access
	// - When keys are stable (more reads than writes)
}

// ========================================
// MUTEX-PROTECTED COUNTER
// ========================================

type Counter struct {
	mu    sync.Mutex
	value int
}

func (c *Counter) Increment() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.value++
}

func (c *Counter) Value() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.value
}

// ========================================
// RWMUTEX-PROTECTED CACHE
// ========================================

type Cache struct {
	mu   sync.RWMutex
	data map[string]string
}

func (c *Cache) Set(key, value string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.data[key] = value
}

func (c *Cache) Get(key string) string {
	c.mu.RLock() // read lock - multiple readers OK
	defer c.mu.RUnlock()
	return c.data[key]
}
