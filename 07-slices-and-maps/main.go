package main

import "fmt"

func main() {
	// ========================================
	// ARRAYS - fixed size, rarely used directly
	// ========================================

	var arr [3]int // fixed size of 3, zero-valued
	arr[0] = 10
	arr[1] = 20
	arr[2] = 30
	fmt.Println("Array:", arr)
	fmt.Println("Array length:", len(arr))

	// Array literal
	arr2 := [3]string{"a", "b", "c"}
	fmt.Println("Array literal:", arr2)

	// Let compiler count
	arr3 := [...]int{1, 2, 3, 4, 5}
	fmt.Println("Auto-sized array:", arr3)

	// WHY arrays are rarely used:
	// - Size is part of the type! [3]int != [4]int
	// - Passed by VALUE (copied!)
	// - Not flexible

	// ========================================
	// SLICES - dynamic, this is what you'll actually use
	// ========================================

	// Slice literal (no size specified)
	names := []string{"Primi", "Go", "Queen"}
	fmt.Println("Slice:", names)
	fmt.Println("Length:", len(names))     // current elements
	fmt.Println("Capacity:", cap(names))   // underlying array size

	// Empty slice with make (when you know approximate size)
	scores := make([]int, 0, 10) // len=0, cap=10
	fmt.Printf("Made slice: len=%d, cap=%d\n", len(scores), cap(scores))

	// ========================================
	// SLICE OPERATIONS
	// ========================================

	// Append - THE way to add elements
	names = append(names, "Fire")
	fmt.Println("After append:", names)

	// Append multiple
	names = append(names, "More", "Names")
	fmt.Println("After multi-append:", names)

	// Slicing (like Python!)
	fmt.Println("names[1:3]:", names[1:3])  // index 1 to 2
	fmt.Println("names[:2]:", names[:2])    // first 2
	fmt.Println("names[2:]:", names[2:])    // from index 2 to end

	// ========================================
	// SLICE GOTCHA - slices share underlying array!
	// ========================================

	original := []int{1, 2, 3, 4, 5}
	sliced := original[1:4] // [2, 3, 4]

	fmt.Println("Before modification:")
	fmt.Println("  original:", original)
	fmt.Println("  sliced:", sliced)

	sliced[0] = 999 // This modifies BOTH!

	fmt.Println("After modifying sliced[0]:")
	fmt.Println("  original:", original) // [1, 999, 3, 4, 5]
	fmt.Println("  sliced:", sliced)     // [999, 3, 4]

	// To avoid this, use copy
	safeCopy := make([]int, len(original))
	copy(safeCopy, original)
	safeCopy[0] = 111
	fmt.Println("Safe copy modified:", safeCopy)
	fmt.Println("Original unchanged:", original)

	// ========================================
	// ITERATING SLICES
	// ========================================

	fruits := []string{"apple", "banana", "cherry"}

	// With index and value
	for i, fruit := range fruits {
		fmt.Printf("fruits[%d] = %s\n", i, fruit)
	}

	// Just values
	for _, fruit := range fruits {
		fmt.Println("Fruit:", fruit)
	}

	// Just indices
	for i := range fruits {
		fmt.Println("Index:", i)
	}

	// ========================================
	// MAPS - key-value pairs (like objects/dicts)
	// ========================================

	// Map literal
	ages := map[string]int{
		"Primi": 27,
		"Go":    15,
		"Linux": 33,
	}
	fmt.Println("Map:", ages)

	// Access value
	fmt.Println("Primi's age:", ages["Primi"])

	// Check if key exists (THE pattern!)
	age, exists := ages["Unknown"]
	if exists {
		fmt.Println("Found:", age)
	} else {
		fmt.Println("Key not found, got zero value:", age)
	}

	// Short version
	if age, ok := ages["Go"]; ok {
		fmt.Println("Go's age:", age)
	}

	// Add/update
	ages["Rust"] = 9
	fmt.Println("After adding Rust:", ages)

	// Delete
	delete(ages, "Rust")
	fmt.Println("After deleting Rust:", ages)

	// Empty map with make
	scores2 := make(map[string]int)
	scores2["test"] = 100
	fmt.Println("Made map:", scores2)

	// ========================================
	// ITERATING MAPS
	// ========================================

	for name, age := range ages {
		fmt.Printf("%s is %d years old\n", name, age)
	}

	// Just keys
	for name := range ages {
		fmt.Println("Name:", name)
	}

	// WARNING: map iteration order is NOT guaranteed!
	// Each run might print in different order

	// ========================================
	// NIL SLICES AND MAPS
	// ========================================

	var nilSlice []int
	fmt.Println("Nil slice:", nilSlice, "is nil:", nilSlice == nil)
	// Safe to append to nil slice!
	nilSlice = append(nilSlice, 1, 2, 3)
	fmt.Println("After append:", nilSlice)

	var nilMap map[string]int
	fmt.Println("Nil map:", nilMap, "is nil:", nilMap == nil)
	// Reading from nil map returns zero value (safe)
	fmt.Println("Read from nil map:", nilMap["key"]) // 0
	// BUT writing to nil map PANICS!
	// nilMap["key"] = 1 // PANIC: assignment to entry in nil map

	// Always initialize maps before writing:
	nilMap = make(map[string]int)
	nilMap["key"] = 1 // now it's safe
	fmt.Println("Initialized map:", nilMap)
}
