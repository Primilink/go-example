package main

// ========================================
// IMPORTS
// ========================================

import (
	"fmt"
	"math"
	"strings"

	// Alias an import
	str "strings" // now you can use str.ToUpper()

	// Import for side effects only (runs init())
	// _ "github.com/lib/pq"

	// Internal packages (we'll create these below)
	// "go-example/09-packages-and-modules/greet"
	// "go-example/09-packages-and-modules/mathutil"
)

func main() {
	// ========================================
	// USING STANDARD LIBRARY
	// ========================================

	fmt.Println("Using math:", math.Sqrt(16))
	fmt.Println("Using strings:", strings.ToUpper("hello"))
	fmt.Println("Using alias:", str.ToLower("HELLO"))

	// ========================================
	// VISIBILITY - Uppercase = Public, lowercase = private
	// ========================================

	// This is THE rule in Go:
	// - Uppercase first letter = exported (public)
	// - Lowercase first letter = unexported (private)

	// fmt.Println  <- Println is uppercase, so it's public
	// strings.Join <- Join is uppercase, so it's public

	// You CAN'T access:
	// strings.hasPrefix <- lowercase, private to strings package

	// Same applies to:
	// - Functions
	// - Types (structs, interfaces)
	// - Struct fields
	// - Constants
	// - Variables

	// Example:
	user := User{
		Name: "Primi", // uppercase = can access from other packages
		age:  27,      // lowercase = only this package can access
	}
	fmt.Printf("User: %+v\n", user)

	// ========================================
	// PACKAGE ORGANIZATION PATTERNS
	// ========================================

	// Flat structure (small projects):
	// myproject/
	// ├── go.mod
	// ├── main.go
	// └── utils.go

	// Package-by-feature (recommended for larger projects):
	// myproject/
	// ├── go.mod
	// ├── main.go
	// ├── user/
	// │   ├── user.go
	// │   ├── repository.go
	// │   └── service.go
	// ├── order/
	// │   ├── order.go
	// │   └── service.go
	// └── internal/      <- special! can't be imported outside module
	//     └── auth/
	//         └── auth.go

	// ========================================
	// GO.MOD - Module definition
	// ========================================

	// go.mod contents:
	// module github.com/primi/myproject
	//
	// go 1.22
	//
	// require (
	//     github.com/gin-gonic/gin v1.9.1
	//     google.golang.org/grpc v1.60.0
	// )

	// ========================================
	// COMMON COMMANDS
	// ========================================

	// go mod init github.com/user/project  <- create new module
	// go mod tidy                          <- sync dependencies
	// go get github.com/pkg/errors         <- add dependency
	// go get -u ./...                      <- update all deps
	// go mod vendor                        <- copy deps locally

	// ========================================
	// INTERNAL PACKAGES
	// ========================================

	// The "internal" directory is special!
	// Packages under internal/ can only be imported by packages
	// in the same module, rooted at the parent of internal/

	// myproject/
	// └── internal/
	//     └── secrets/   <- ONLY myproject can import this

	// ========================================
	// INIT FUNCTION
	// ========================================

	// Each package can have init() functions
	// They run automatically when package is imported
	// Order: imported packages init() -> this package init() -> main()

	fmt.Println("Check the init message above!")
}

// ========================================
// EXAMPLE STRUCT WITH VISIBILITY
// ========================================

type User struct {
	Name  string // exported - other packages can access
	Email string // exported
	age   int    // unexported - only this package
}

// Exported method
func (u *User) GetAge() int {
	return u.age
}

// Unexported method
func (u *User) validateAge() bool {
	return u.age >= 0
}

// ========================================
// INIT RUNS BEFORE MAIN
// ========================================

func init() {
	fmt.Println("=== init() called! Setting things up... ===")
}

// You can have multiple init() in the same file!
func init() {
	fmt.Println("=== second init() also runs! ===")
}
