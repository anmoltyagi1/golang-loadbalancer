package main

import "fmt"

// Define a type for your set
type Set map[string]bool

// Initialize a new set
func NewSet() Set {
    return make(Set)
}

// Add an element to the set
func (set Set) Add(element string) {
    set[element] = true
}

// Remove an element from the set
func (set Set) Remove(element string) {
    delete(set, element)
}

// Check if an element exists in the set
func (set Set) Contains(element string) bool {
    _, exists := set[element]
    return exists
}

// Get the size of the set
func (set Set) Size() int {
    return len(set)
}

func printSet(set Set) {
    fmt.Print("{")
    for key := range set {
        fmt.Printf("%v, ", key)
    }
    fmt.Println("}")
}
