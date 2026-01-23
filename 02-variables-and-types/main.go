package main

import "fmt"

func main() {
	// ========================================
	// VARIABLE DECLARATION - 3 ways
	// ========================================

	// 1. Full declaration (rarely used)
	var name string = "Primi"

	// 2. Type inference with var (useful for zero values)
	var age = 27 // Go infers int

	// 3. Short declaration := (most common, only inside functions)
	city := "Cali" // infers string

	fmt.Println(name, age, city)

	// ========================================
	// ZERO VALUES - Go's default values (no undefined/null hell!)
	// ========================================

	var i int     // 0
	var f float64 // 0.0
	var b bool    // false
	var s string  // "" (empty string)

	fmt.Printf("Zero values: int=%d, float=%f, bool=%t, string=%q\n", i, f, b, s)

	// ========================================
	// BASIC TYPES
	// ========================================

	// Integers
	var small int8 = 127           // -128 to 127
	var medium int16 = 32767       // -32768 to 32767
	var regular int32 = 2147483647 // you know this from C
	var big int64 = 9223372036854775807
	var auto int = 42 // platform dependent (64-bit on your machine)

	// Unsigned
	var unsigned uint = 42
	var byte_ byte = 255 // alias for uint8

	// Floats
	var price float64 = 19.99 // use this by default
	var fast float32 = 3.14   // only if memory constrained

	// Strings (UTF-8 by default, immutable)
	emoji := "Go is fire 🔥"

	// Runes (Unicode code point, alias for int32)
	var r rune = '🔥' // single quotes = rune

	// Bool
	isAwesome := true

	fmt.Println(small, medium, regular, big, auto)
	fmt.Println(unsigned, byte_)
	fmt.Println(price, fast)
	fmt.Println(emoji, r, isAwesome)

	// ========================================
	// TYPE CONVERSION - explicit only! (no implicit casting)
	// ========================================

	var x int = 42
	var y float64 = float64(x) // must be explicit
	var z int = int(y)         // must be explicit

	fmt.Printf("Conversions: int=%d, float=%f, back to int=%d\n", x, y, z)

	// This WON'T compile (unlike JS/PHP):
	// var bad float64 = x  // ERROR: cannot use x (type int) as type float64

	// ========================================
	// CONSTANTS
	// ========================================

	const pi = 3.14159
	const (
		statusOK    = 200
		statusError = 500
	)

	fmt.Println("Constants:", pi, statusOK, statusError)

	// ========================================
	// iota - auto-incrementing constants (super useful for enums)
	// ========================================

	const (
		Sunday = iota // 0
		Monday        // 1
		Tuesday       // 2
		Wednesday     // 3
		Thursday      // 4
		Friday        // 5
		Saturday      // 6
	)

	fmt.Println("Days:", Sunday, Monday, Friday)

	// Common pattern: skip zero value
	const (
		_           = iota // skip 0
		KB int64 = 1 << (10 * iota) // 1024
		MB                          // 1048576
		GB                          // 1073741824
	)

	fmt.Printf("Sizes: KB=%d, MB=%d, GB=%d\n", KB, MB, GB)
}
