package main

import (
	cryptorand "crypto/rand"
	"encoding/binary"
	"fmt"
	"math"
	"math/rand/v2"
)

type Node struct {
	value int
	next  []*Node // array of pointers to next nodes at each level
}

type SkipList struct {
	head         *Node
	maxLevel     int
	currentLevel int
	probability  float64
}

func NewSkipList(maxLevel int, probability float64) *SkipList {
	// Create head node with max possible level
	head := &Node{
		value: math.MinInt, // Sentinel value
		next:  make([]*Node, maxLevel),
	}

	return &SkipList{
		head:         head,
		maxLevel:     maxLevel,
		currentLevel: 1,
		probability:  probability,
	}
}

func SecureFloat64() (float64, error) {
	// Read 8 cryptographically secure random bytes
	var buf [8]byte
	_, err := cryptorand.Read(buf[:])
	if err != nil {
		return 0, err
	}

	// Convert bytes to a uint64
	val := binary.LittleEndian.Uint64(buf[:])

	// Shift bits and divide to get a value in [0, 1.0)
	// Use 53 bits as they form the mantissa of a float64
	const maxUint53 = 1 << 53
	return float64(val>>11) / maxUint53, nil
}

func (sl *SkipList) getRandomLevel() int {
	level := 1

	for {
		random, err := SecureFloat64() // Truly generate a random value - could have been avoided
		if err != nil {
			// Fallback scenario
			random = rand.Float64()
		}
		if random <= sl.probability && level <= sl.maxLevel {
			level += 1
		} else {
			break
		}
	}

	return level
}

func (sl *SkipList) Insert(value int) {
	// Step 1: Determine random level for the new node
	level := sl.getRandomLevel()

	fmt.Printf("Level selected for %d: %d\n", value, level)

	// Step 2: Create update array to track insertion points at each level
	update := make([]*Node, sl.maxLevel)
	current := sl.head

	// Step 3: Search for insertion points from top to bottom
	for i := sl.currentLevel - 1; i >= 0; i-- {
		// Move forward while next node exists and its value < insert value
		for current.next[i] != nil && current.next[i].value < value {
			current = current.next[i]
		}
		// Record the node at this level where we should insert
		update[i] = current
	}

	// Step 4: Move to the actual insertion point at level 0
	current = current.next[0]

	// Step 5: Check if value already exists (optional - update or skip)
	if current != nil && current.value == value {
		// Value already exists - you can choose to:
		// - Skip insertion
		// - Update the node
		// For now, we'll just return
		fmt.Printf("Value %d already exists\n", value)
		return
	}

	newNode := &Node{
		value: value,
		next:  make([]*Node, level),
	}

	for i := range level {
		// If this level is higher than current max level
		if i >= sl.currentLevel {
			// Connect head directly to new node
			sl.head.next[i] = newNode
			newNode.next[i] = nil
		} else {
			// Standard insertion: newNode points to what update[i] was pointing to
			newNode.next[i] = update[i].next[i]
			// update[i] now points to newNode
			update[i].next[i] = newNode
		}
	}

	// Step 8: Update current level if necessary
	if level > sl.currentLevel {
		sl.currentLevel = level
	}
}

// Returns the node if exists, null otherwise
func (sl *SkipList) Search(number int) *Node {
	current := sl.head
	for i := sl.currentLevel - 1; i >= 0; i-- {
		for current.next[i] != nil {
			if current.next[i].value == number {
				fmt.Printf("Found number: %d\n", number)
				return current.next[i]
			}
			if current.next[i].value < number {
				current = current.next[i]
			} else {
				break
			}
		}
	}

	fmt.Printf("Didn't find the number %d\n", number)
	return nil
}

func (sl *SkipList) Delete(number int) {
	update := make([]*Node, sl.currentLevel)
	current := sl.head

	deletionMaxLevel := -1

	// Step 3: Search for insertion points from top to bottom
	for i := sl.currentLevel - 1; i >= 0; i-- {
		// Move forward while next node exists and its value < insert value
		for current.next[i] != nil && current.next[i].value < number {
			current = current.next[i]
		}

		if deletionMaxLevel == -1 && current.next[i] != nil && current.next[i].value == number {
			deletionMaxLevel = i
		}

		update[i] = current
	}

	if deletionMaxLevel == -1 {
		fmt.Println("Number doesn't exist: ", number)
		return
	}

	for i := deletionMaxLevel; i >= 0; i-- {
		temp := update[i].next[i]
		if temp != nil {
			update[i].next[i] = temp.next[i]
		}
	}

	// Decrement the current level if the tops levels are empty
	for sl.currentLevel > 1 && sl.head.next[sl.currentLevel-1] == nil {
		sl.currentLevel--
	}

	fmt.Println("Deleted number: ", number)
}

// Display prints the skip list structure (for debugging)
func (sl *SkipList) Display() {
	fmt.Println("\n=== Skip List Structure ===")
	for i := sl.currentLevel - 1; i >= 0; i-- {
		fmt.Printf("Level %d: HEAD", i+1)
		current := sl.head.next[i]
		for current != nil {
			fmt.Printf(" -> %d", current.value)
			current = current.next[i]
		}
		fmt.Println(" -> NULL")
	}
	fmt.Println("===========================")
}
