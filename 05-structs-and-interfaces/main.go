package main

import (
	"fmt"
	"math"
)

func main() {
	// ========================================
	// STRUCTS - Go's "classes" (but not really)
	// ========================================

	// Create a struct instance
	user1 := User{
		Name:  "Primi",
		Email: "primi@example.com",
		Age:   27,
	}
	fmt.Println("User1:", user1)

	// Access fields
	fmt.Println("Name:", user1.Name)

	// Modify fields
	user1.Age = 28
	fmt.Println("Updated age:", user1.Age)

	// Partial initialization (other fields get zero values)
	user2 := User{Name: "Anonymous"}
	fmt.Printf("User2: %+v\n", user2) // %+v shows field names

	// Zero value struct (all fields are zero values)
	var user3 User
	fmt.Printf("User3 (zero): %+v\n", user3)

	// ========================================
	// METHODS - functions attached to types
	// ========================================

	user1.Greet()
	fmt.Println("Birth year:", user1.BirthYear())

	// ========================================
	// POINTER vs VALUE RECEIVERS
	// ========================================

	fmt.Println("\nBefore birthday:", user1.Age)
	user1.HaveBirthday() // modifies the actual struct
	fmt.Println("After birthday:", user1.Age)

	// Why this matters:
	// - Value receiver: method gets a COPY (can't modify original)
	// - Pointer receiver: method gets the ADDRESS (can modify original)

	// ========================================
	// CONSTRUCTORS - Go doesn't have them, use factory functions
	// ========================================

	user4 := NewUser("Factory", "factory@test.com", 30)
	fmt.Printf("From factory: %+v\n", user4)

	// ========================================
	// INTERFACES - implicit implementation (duck typing!)
	// ========================================

	// No "implements" keyword needed!
	// If a type has the methods, it implements the interface

	circle := Circle{Radius: 5}
	rect := Rectangle{Width: 4, Height: 3}

	// Both implement Shape interface
	printShapeInfo(circle)
	printShapeInfo(rect)

	// Slice of interfaces - polymorphism!
	shapes := []Shape{circle, rect, Circle{Radius: 10}}
	totalArea := 0.0
	for _, s := range shapes {
		totalArea += s.Area()
	}
	fmt.Printf("Total area: %.2f\n", totalArea)

	// ========================================
	// EMPTY INTERFACE - any type (like "any" in TS)
	// ========================================

	var anything interface{} // or just: var anything any
	anything = 42
	fmt.Println("anything (int):", anything)

	anything = "now I'm a string"
	fmt.Println("anything (string):", anything)

	anything = user1
	fmt.Println("anything (User):", anything)

	// Type assertion - get the concrete type back
	if str, ok := anything.(string); ok {
		fmt.Println("It's a string:", str)
	} else {
		fmt.Println("Not a string!")
	}

	// Type switch
	checkType(42)
	checkType("hello")
	checkType(3.14)
	checkType(user1)

	// ========================================
	// EMBEDDING - composition over inheritance
	// ========================================

	admin := Admin{
		User:  User{Name: "Super Primi", Email: "admin@test.com", Age: 27},
		Level: "superadmin",
	}

	// Admin "inherits" User methods (not really inheritance, it's embedding)
	admin.Greet() // calls User.Greet()

	// Admin can have its own methods too
	admin.AdminGreet()

	// Access embedded fields directly
	fmt.Println("Admin name:", admin.Name) // not admin.User.Name (though that works too)
}

// ========================================
// STRUCT DEFINITIONS
// ========================================

type User struct {
	Name  string
	Email string
	Age   int
}

// Method with VALUE receiver (gets a copy)
func (u User) Greet() {
	fmt.Printf("Hi, I'm %s!\n", u.Name)
}

func (u User) BirthYear() int {
	return 2025 - u.Age
}

// Method with POINTER receiver (can modify)
func (u *User) HaveBirthday() {
	u.Age++
}

// Factory function (constructor pattern)
func NewUser(name, email string, age int) *User {
	return &User{
		Name:  name,
		Email: email,
		Age:   age,
	}
}

// ========================================
// INTERFACES
// ========================================

type Shape interface {
	Area() float64
	Perimeter() float64
}

// Circle implements Shape (implicitly!)
type Circle struct {
	Radius float64
}

func (c Circle) Area() float64 {
	return math.Pi * c.Radius * c.Radius
}

func (c Circle) Perimeter() float64 {
	return 2 * math.Pi * c.Radius
}

// Rectangle also implements Shape
type Rectangle struct {
	Width, Height float64
}

func (r Rectangle) Area() float64 {
	return r.Width * r.Height
}

func (r Rectangle) Perimeter() float64 {
	return 2 * (r.Width + r.Height)
}

// Function accepting interface
func printShapeInfo(s Shape) {
	fmt.Printf("Area: %.2f, Perimeter: %.2f\n", s.Area(), s.Perimeter())
}

// Type switch
func checkType(val interface{}) {
	switch v := val.(type) {
	case int:
		fmt.Println("It's an int:", v)
	case string:
		fmt.Println("It's a string:", v)
	case float64:
		fmt.Println("It's a float:", v)
	default:
		fmt.Printf("Unknown type: %T\n", v)
	}
}

// ========================================
// EMBEDDING (composition)
// ========================================

type Admin struct {
	User  // embedded - Admin "has a" User, but feels like "is a"
	Level string
}

func (a Admin) AdminGreet() {
	fmt.Printf("I'm admin %s with level %s\n", a.Name, a.Level)
}
