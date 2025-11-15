package main

import (
	"fmt"
)

func main() {
	skiplist := NewSkipList(16, 0.5)

	fmt.Println("Inserting element: 10")

	skiplist.Insert(10)
	skiplist.Display()

	fmt.Println("Inserting element: 25")

	skiplist.Insert(25)
	skiplist.Display()

	fmt.Println("Inserting element: 5")

	skiplist.Insert(5)
	skiplist.Display()

	fmt.Println("Inserting element: 17")

	skiplist.Insert(17)
	skiplist.Display()

	fmt.Println("Inserting element: 30")

	skiplist.Insert(30)
	skiplist.Display()

	skiplist.Search(10)
	skiplist.Search(25)
	skiplist.Search(5)
	skiplist.Search(30)
	skiplist.Search(17)
	skiplist.Search(24)

	skiplist.Delete(10)
	skiplist.Display()
}
