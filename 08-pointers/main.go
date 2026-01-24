package main

import "fmt"

func main() {
	// ========================================
	// POINTER BASICS - you know this from C!
	// ========================================

	x := 42
	p := &x // p is a pointer to x

	fmt.Println("x value:", x)
	fmt.Println("x address:", &x)
	fmt.Println("p (address):", p)
	fmt.Println("*p (dereferenced):", *p)

	// Modify through pointer
	*p = 100
	fmt.Println("x after *p = 100:", x)

	// ========================================
	// DIFFERENCE FROM C
	// ========================================

	// NO pointer arithmetic! This won't compile:
	// p++        // NOPE
	// p + 1      // NOPE
	// *(p + 1)   // NOPE

	// Go is memory-safe, no segfaults from pointer math!

	// ========================================
	// ZERO VALUE OF POINTER = nil
	// ========================================

	var nilPtr *int
	fmt.Println("Nil pointer:", nilPtr)

	if nilPtr == nil {
		fmt.Println("Pointer is nil, don't dereference!")
	}

	// Dereferencing nil = panic
	// fmt.Println(*nilPtr) // PANIC: runtime error

	// ========================================
	// NEW() - allocates and returns pointer
	// ========================================

	ptr := new(int) // allocates int, returns *int
	fmt.Println("new(int):", ptr, "value:", *ptr) // zero-valued

	*ptr = 50
	fmt.Println("After assignment:", *ptr)

	// ========================================
	// POINTERS AND FUNCTIONS
	// ========================================

	val := 10

	// Pass by value - original unchanged
	doubleValue(val)
	fmt.Println("After doubleValue:", val) // still 10

	// Pass by pointer - original modified
	doublePointer(&val)
	fmt.Println("After doublePointer:", val) // now 20

	// ========================================
	// POINTERS AND STRUCTS
	// ========================================

	user := User{Name: "Primi", Age: 27}

	// Without pointer - gets a copy
	celebrateBirthdayCopy(user)
	fmt.Println("After copy birthday:", user.Age) // still 27

	// With pointer - modifies original
	celebrateBirthdayPointer(&user)
	fmt.Println("After pointer birthday:", user.Age) // now 28

	// ========================================
	// AUTOMATIC DEREFERENCING
	// ========================================

	userPtr := &user

	// In C you'd need: (*userPtr).Name
	// Go does it automatically:
	fmt.Println("Name via pointer:", userPtr.Name) // no need for (*userPtr).Name

	// ========================================
	// RETURNING POINTERS (safe in Go!)
	// ========================================

	// In C this would be a bug (returning pointer to stack variable)
	// In Go, the compiler moves it to heap automatically (escape analysis)

	newUser := createUser("Factory Queen", 30)
	fmt.Println("Created user:", newUser.Name)

	// ========================================
	// WHEN TO USE POINTERS
	// ========================================

	// Use pointers when:
	// 1. You need to modify the original value
	// 2. Struct is large (avoid copying)
	// 3. You need to represent "absence" (nil)

	// Use values when:
	// 1. Struct is small (copying is cheap)
	// 2. You want immutability
	// 3. You want simpler code

	// ========================================
	// POINTER TO POINTER (rarely needed)
	// ========================================

	a := 5
	pa := &a
	ppa := &pa

	fmt.Println("a:", a)
	fmt.Println("*pa:", *pa)
	fmt.Println("**ppa:", **ppa)
}

type User struct {
	Name string
	Age  int
}

func doubleValue(x int) {
	x = x * 2 // modifies copy only
}

func doublePointer(x *int) {
	*x = *x * 2 // modifies original
}

func celebrateBirthdayCopy(u User) {
	u.Age++ // modifies copy
}

func celebrateBirthdayPointer(u *User) {
	u.Age++ // modifies original
}

func createUser(name string, age int) *User {
	// This is SAFE in Go! Compiler handles it
	user := User{Name: name, Age: age}
	return &user // would be bug in C, fine in Go
}
