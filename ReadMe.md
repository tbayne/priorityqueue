# priorityQueue
A golang implementation of a priority queue based upon the standard 
go library heap functionality.

Items are pushed and popped from the queue, and can have their 
priority changed as needed.

Items are always of type `QItem`, and the `QItem.value` is where
you put your values.  Additionally

```
type QItem struct {
	ID       string      // The ID of the queue item
	ParentID string      // The parent ID of the queue item
	Value    interface{} // The value of the item; can hold any type.
	Priority int         // The Priority of the item in the queue.
}
``` 

Note that it the ID value of the item is not checked for uniqueness.

## Example Usage
---
```
package main

import (
	"fmt"
	"runtime"
	"strconv"

	pq "tbayne/priorityqueue"
)

func main() {
	runtime.GOMAXPROCS(8)
	q := pq.NewPriorityQueue()

	// Push 5 items
	for i := 0; i <= 100; i++ {
		var y = pq.QItem{
			ParentID: "Parent ID String",
			ID:       strconv.Itoa(i),
			Value:    "test",
			Priority: 1 + i,
		}
		q.Push(y)
	}
	fmt.Printf("Queue Length: %d\n", q.Len())
	// Push a higher priority item
	y := pq.QItem{ID: "12345", Value: 1234555, Priority: 500}
	q.Push(y)
	fmt.Printf("Queue Length: %d\n", q.Len())
	for q.Len() > 0 {
		x, _ := q.Pop()
		fmt.Printf("QItem: %+v \n", x)
	}

	fmt.Printf("QItems Length: %d\n", q.Len())
        q.Destroy()
}
```
---

## Usage Notes
* Always access the queue through the provided member functions

* Always Clear() the queue when you are finished with it

* Calling Destory() clears the queue and deletes the underlying array

* See `priorityqueue_test.go` for more usage examples
