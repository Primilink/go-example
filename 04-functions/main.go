package main

import (
	"errors"
	"fmt"
)

func main() {
	// ========================================
	// BASIC FUNCTIONS
	// ========================================

	result := add(5, 3)
	fmt.Println("5 + 3 =", result)

	greet("Primi")

	// ========================================
	// MULTIPLE RETURN VALUES - Go's killer feature!
	// ========================================

	quotient, remainder := divide(17, 5)
	fmt.Printf("17 / 5 = %d remainder %d\n", quotient, remainder)

	// The classic Go pattern: value, error
	val, err := safeDivide(10, 0)
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Result:", val)
	}

	// Don't need a return value? Use blank identifier
	_, justRemainder := divide(17, 5)
	fmt.Println("Just the remainder:", justRemainder)

	// ========================================
	// NAMED RETURN VALUES
	// ========================================

	width, height := dimensions()
	fmt.Printf("Dimensions: %d x %d\n", width, height)

	// ========================================
	// VARIADIC FUNCTIONS (like ...spread in JS)
	// ========================================

	total := sum(1, 2, 3, 4, 5)
	fmt.Println("Sum:", total)

	// Pass a slice with ...
	numbers := []int{10, 20, 30}
	fmt.Println("Sum of slice:", sum(numbers...))

	// ========================================
	// FUNCTIONS AS VALUES (first-class functions)
	// ========================================

	// Assign function to variable
	multiply := func(a, b int) int {
		return a * b
	}
	fmt.Println("3 * 4 =", multiply(3, 4))

	// Pass function as argument
	result2 := operate(10, 5, add)
	fmt.Println("operate(10, 5, add) =", result2)

	result3 := operate(10, 5, func(a, b int) int {
		return a - b
	})
	fmt.Println("operate(10, 5, subtract) =", result3)

	// ========================================
	// CLOSURES - functions that capture variables
	// ========================================

	counter := makeCounter()
	fmt.Println("Counter:", counter()) // 1
	fmt.Println("Counter:", counter()) // 2
	fmt.Println("Counter:", counter()) // 3

	// Each call to makeCounter creates a new counter
	counter2 := makeCounter()
	fmt.Println("Counter2:", counter2()) // 1 (fresh!)

	// ========================================
	// DEFER with functions
	// ========================================

	fmt.Println("\n--- Defer with cleanup pattern ---")
	doSomethingWithCleanup()

	// ========================================
	// INIT FUNCTION - runs automatically before main()
	// ========================================
	// See the init() function below - it ran before main!
	fmt.Println("Config loaded:", config)
}

// ========================================
// FUNCTION DEFINITIONS
// ========================================

// Basic function
func add(a int, b int) int {
	return a + b
}

// Same type params can be shortened
func greet(name string) {
	fmt.Println("Hello,", name)
}

// Multiple return values
func divide(a, b int) (int, int) {
	return a / b, a % b
}

// Return value + error (THE Go pattern)
func safeDivide(a, b int) (int, error) {
	if b == 0 {
		return 0, errors.New("cannot divide by zero")
	}
	return a / b, nil
}

// Named return values (auto-initialized to zero values)
func dimensions() (width int, height int) {
	width = 1920
	height = 1080
	return // naked return - returns named values
}

// Variadic function
func sum(nums ...int) int {
	total := 0
	for _, n := range nums {
		total += n
	}
	return total
}

// Function as parameter
func operate(a, b int, op func(int, int) int) int {
	return op(a, b)
}

// Closure - returns a function that captures state
func makeCounter() func() int {
	count := 0
	return func() int {
		count++
		return count
	}
}

// Defer cleanup pattern
func doSomethingWithCleanup() {
	fmt.Println("Starting work...")
	defer fmt.Println("Cleanup done!") // always runs

	fmt.Println("Doing work...")
	// even if we panic here, defer still runs!
}

// ========================================
// INIT FUNCTION - special function, runs before main
// ========================================

var config string

func init() {
	// Use for:
	// - Loading config
	// - Initializing global state
	// - Registering drivers
	config = "loaded from init()"
	fmt.Println("init() called before main()!")
}
