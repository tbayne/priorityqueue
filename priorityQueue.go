package priorityqueue

import (
	"container/heap"
	"fmt"
	"sync"
)

// An Item is something we manage in a Priority queue.
type Item struct {
	ID       int
	Value    interface{} // The value of the item; can hold any type.
	Priority int         // The Priority of the item in the queue.
	// The index is needed by update and is maintained by the heap.Interface methods.
	index int // The index of the item in the heap.
}

// A Items implements heap.Interface and holds Items.
type Items []*Item

type PriorityQueue struct {
	m    sync.Mutex
	data Items
}

func NewPriorityQueue() *PriorityQueue {

	var pq PriorityQueue

	// Initialize our heap backing store
	pq.data = make(Items, 0)

	return &pq
}

func (pq *PriorityQueue) Len() int {
	pq.m.Lock()
	defer pq.m.Unlock()
	return pq.data.Len()
}

func (pq *PriorityQueue) Push(i Item) {

	pq.m.Lock()
	defer pq.m.Unlock()
	pq.data.Push(i)

}

// Todo: Pop() Return nil and an error if there is nothing in the queue
func (pq *PriorityQueue) Pop() (*Item, error) {
	if pq.data.Len() > 0 {
		pq.m.Lock()
		defer pq.m.Unlock()
		r := pq.data.Pop()
		return r.(*Item), nil
	}
	return nil, fmt.Errorf("Queue is empty, nothing to Pop")
}

/* Implement the heap interface methods: Len, Less, Swap, Push, and Pop */

func (pq Items) Len() int {
	return len(pq)
}

func (pq Items) Less(i, j int) bool {
	// We want Pop to give us the highest, not lowest, Priority so we use greater than here.
	//   Priority 10 is greater than Priority 1, the larger the number the higher the Priority
	return pq[i].Priority > pq[j].Priority
}

func (pq Items) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *Items) Push(x interface{}) {
	n := len(*pq)
	item := x.(Item)
	item.index = n
	*pq = append(*pq, &item)
}

func (pq *Items) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil  // avoid memory leak
	item.index = -1 // for safety
	*pq = old[0 : n-1]
	return item
}

// update modifies the Priority and value of an Item in the queue.
func (pq *Items) update(item *Item, value string, priority int) {
	item.Value = value
	item.Priority = priority
	heap.Fix(pq, item.index)
}

// This example creates a Items with some items, adds and manipulates an item,
// and then removes the items in Priority order.

/*
import (
	"container/heap"
	"fmt"
)

type Item struct {
	ID       int
	Value    interface{}
	Priority int
	index    int
}

type Items struct {
	items    []Item
	quitchan chan bool
	// pushchan chan heapPushChanMsg
	//popchan  chan heapPopChanMsg
}


// heapPopChanMsg - the message structure for a pop chan
type heapPopChanMsg struct {
	// h heap.Interface
	result chan interface{}
}

// heapPushChanMsg - the message structure for a push chan
type heapPushChanMsg struct {
	// h heap.Interface
	x interface{}
}

func NewPriorityQueue(len int) Items {

	var pq Items
	items := make([]Item, len)
	pq.items = items
	// Create our syncronizing channels:
	//pq.pushchan = make(chan heapPushChanMsg, 10)
	//pq.popchan = make(chan heapPopChanMsg, 10)
	heap.Init(pq)

	// Start the watcher
	//pq.quitchan = pq.watchHeapOps()

	// How many items do we actually contain

	return pq

}

func (pq *Items) Close() {
	// Drain queue
	// Stop the watcher
	//pq.quitchan <- true
	//close(pq.quitchan)
	//close(pq.popchan)
	//close(pq.pushchan)
	// Destroy queue
	pq.items = nil
}



// Len - get the length of the queue
func (pq Items) Len() int {
	return len(pq.items)
}

// Less - determine which is a greater Priority
func (pq Items) Less(i, j int) bool {
	return pq.items[i].Priority < pq.items[j].Priority
}

// Swap - implementation of swap for the heap interface
func (pq Items) Swap(i, j int) {
	pq.items[i], pq.items[j] = pq.items[j], pq.items[i]
	pq.items[i].index = i
	pq.items[j].index = j
}

// HeapPush - safely push item to a heap interface
func (pq Items) Push(x interface{}) {

	i := x.(Item)
	pq.items = append(pq.items, i)
	pq.pushchan <- heapPushChanMsg{x: x}

}

// HeapPop - safely pop item from a heap interface
func (pq Items) Pop() interface{} {
	var result = make(chan interface{})
	pq.popchan <- heapPopChanMsg{
		result: result,
	}
	pq.count = pq.count - 1
	return <-result
}

// watchHeapOps - watch for push/pops to our heap, and serializing the operations
// with channels
func (pq *Items) watchHeapOps() chan bool {
	var quit = make(chan bool)
	go func() {
		for {
			select {
			case <-quit:
				// TODO: update to quit gracefully
				// TODO: maybe need to dump state somewhere?
				return

			case popMsg := <-pq.popchan:
				popMsg.result <- heap.Pop(pq)
				fmt.Println("Items: Popped item from the Items")

			case pushMsg := <-pq.pushchan:
				heap.Push(pq, pushMsg.x)
				fmt.Println("Items: Pushed item to the Items")
			}
		}
	}()
	return quit
}
*/
