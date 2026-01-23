package main

import "fmt"

func main() {
	// ========================================
	// IF/ELSE - no parentheses needed!
	// ========================================

	age := 27

	if age >= 21 {
		fmt.Println("You can drink 🍺")
	} else if age >= 18 {
		fmt.Println("You can vote")
	} else {
		fmt.Println("You're a baby")
	}

	// IF with init statement (super useful, scopes variable to if block)
	if score := 85; score >= 70 {
		fmt.Printf("Score %d: Passed!\n", score)
	}
	// score doesn't exist here - scoped to the if block only!

	// ========================================
	// FOR - the ONLY loop in Go (no while, no do-while)
	// ========================================

	// Classic for
	for i := 0; i < 3; i++ {
		fmt.Println("Classic:", i)
	}

	// While-style (just omit init and post)
	count := 0
	for count < 3 {
		fmt.Println("While-style:", count)
		count++
	}

	// Infinite loop (use break to exit)
	attempts := 0
	for {
		attempts++
		if attempts >= 3 {
			fmt.Println("Breaking out after", attempts, "attempts")
			break
		}
	}

	// Range - iterating over collections (you'll use this A LOT)
	names := []string{"Primi", "Go", "Queen"}

	for index, value := range names {
		fmt.Printf("Index %d: %s\n", index, value)
	}

	// Don't need index? Use blank identifier _
	for _, name := range names {
		fmt.Println("Name:", name)
	}

	// Only need index?
	for i := range names {
		fmt.Println("Just index:", i)
	}

	// ========================================
	// SWITCH - no break needed! (implicit break)
	// ========================================

	day := "Friday"

	switch day {
	case "Monday":
		fmt.Println("Ugh, Monday")
	case "Friday":
		fmt.Println("TGIF! 🎉")
	case "Saturday", "Sunday": // multiple values
		fmt.Println("Weekend vibes")
	default:
		fmt.Println("Just another day")
	}

	// Switch with no condition (cleaner than if-else chains)
	score := 85

	switch {
	case score >= 90:
		fmt.Println("A - Excellent!")
	case score >= 80:
		fmt.Println("B - Good job!")
	case score >= 70:
		fmt.Println("C - Passing")
	default:
		fmt.Println("F - Try again")
	}

	// fallthrough - explicit (opposite of C where break is explicit)
	switch num := 1; num {
	case 1:
		fmt.Println("One")
		fallthrough // continues to next case
	case 2:
		fmt.Println("This also prints because of fallthrough")
	case 3:
		fmt.Println("This won't print")
	}

	// ========================================
	// DEFER - runs when function exits (LIFO order)
	// ========================================

	fmt.Println("\n--- Defer examples ---")
	fmt.Println("Start")

	defer fmt.Println("Deferred 1 - runs last")
	defer fmt.Println("Deferred 2 - runs second to last")

	fmt.Println("Middle")
	fmt.Println("End")

	// Defers run here, in reverse order (LIFO - stack)

	// Common use: cleanup resources
	// file := openFile("data.txt")
	// defer file.Close()  // guaranteed to close, even if panic occurs
	// ... work with file ...
}
