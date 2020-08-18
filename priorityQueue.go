package priorityqueue

import (
	"container/heap"
	"fmt"
	"sync"
)

// An QItem is something we manage in a Priority queue.
type QItem struct {
	ID       string
	ParentID string
	Value    interface{} // The value of the item; can hold any type.
	Priority int         // The Priority of the item in the queue.

	// The index is needed by update and is maintained by the heap.Interface methods.
	index int // The index of the item in the heap.
}

// A QItems implements heap.Interface and holds QItems.
type QItems []*QItem

type PriorityQueue struct {
	m         sync.Mutex
	available bool
	data      QItems
}

func NewPriorityQueue() *PriorityQueue {

	var pq PriorityQueue

	// Initialize our heap backing store
	pq.data = make(QItems, 0)
	heap.Init(&pq.data)

	return &pq
}

// Destroy clears the queue and destroys the underlying storage
func (pq *PriorityQueue) Destroy() {
	pq.Clear()
	pq.data = nil
}

func (pq *PriorityQueue) Len() int {
	pq.m.Lock()
	defer pq.m.Unlock()
	return pq.data.Len()
}

func (pq *PriorityQueue) Push(i QItem) {

	pq.m.Lock()
	defer pq.m.Unlock()
	heap.Push(&pq.data, i)

}

func (pq *PriorityQueue) Pop() (*QItem, error) {
	if pq.data.Len() > 0 {
		pq.m.Lock()
		defer pq.m.Unlock()
		r := heap.Pop(&pq.data)
		return r.(*QItem), nil
	}
	return nil, fmt.Errorf("queue is empty, nothing to Pop")
}

// UpdatePriorityById() updates the priority of an item in the queue
func (pq *PriorityQueue) UpdatePriorityByParentId(parentID string, priority int) int {
	pq.m.Lock()
	defer pq.m.Unlock()
	index := -1
	itemsUpdated := 0
	// Walk every item in the queue
	for _, element := range pq.data {
		if element.ParentID == parentID {
			index = element.index
			pq.data.update(pq.data[index], priority)
			itemsUpdated++
		}
	}
	return itemsUpdated
}

/* Clear drains all items from the queue */
func (pq *PriorityQueue) Clear() {
	pq.m.Lock()
	defer pq.m.Unlock()
	for pq.data.Len() > 0 {
		x := heap.Pop(&pq.data)
		if x != nil {
			x = nil
		}
	}
}

func (pq *PriorityQueue) locateItemByID(id string) (int, error) {
	var index = -1
	for _, element := range pq.data {
		if element.ID == id {
			index = element.index
			break
		}
	}
	if index == -1 {
		return -1, fmt.Errorf("ID Not found: [%s]", id)
	}
	return index, nil
}

// DeleteItemById() deletes an item from the queue based on the ID

func (pq *PriorityQueue) DeleteItemById(id string) error {
	pq.m.Lock()
	defer pq.m.Unlock()
	index, err := pq.locateItemByID(id)
	if err != nil {
		return err
	}
	err = pq.data.delete(index)
	if err != nil {
		return err
	}
	return nil
}

// DeleteItemByParentId() deletes all items with a matching ParentID from the queue based on the Parent ID
// Note deletes can be expensive

func (pq *PriorityQueue) DeleteItemsByParentId(parentID string) (int, error) {
	pq.m.Lock()
	defer pq.m.Unlock()

	itemsDeleted := 0

	// A place to collect the indexes for the items we want to delete
	var indexesToDelete []int

	for _, element := range pq.data {
		if element.ParentID == parentID {
			indexesToDelete = append(indexesToDelete, element.index)
		}
	}

	for index := range indexesToDelete {

		err := pq.data.delete(index)
		if err != nil {
			return itemsDeleted, err
		}
		itemsDeleted++
	}

	return itemsDeleted, nil
}

/* Implement the heap interface methods: Len, Less, Swap, Push, and Pop */

func (qData QItems) Len() int {
	return len(qData)
}

func (qData QItems) Less(i, j int) bool {
	// We want Pop to give us the highest, not lowest, Priority so we use greater than here.
	//   Priority 10 is greater than Priority 1, the larger the number the higher the Priority
	return qData[i].Priority > qData[j].Priority
}

func (qData QItems) Swap(i, j int) {
	qData[i], qData[j] = qData[j], qData[i]
	qData[i].index = i
	qData[j].index = j
}

// Push adds an item to the queue
// Note do NOT call this directly, this is called by the heap
// implementation and not your application
func (qData *QItems) Push(x interface{}) {
	n := len(*qData)
	item := x.(QItem)
	item.index = n
	*qData = append(*qData, &item)
}

// Pop removes an item to the queue
// Note do NOT call this directly, this is called by the heap
// implementation and not your application
func (qData *QItems) Pop() interface{} {
	old := *qData
	n := len(old)
	item := old[n-1]
	old[n-1] = nil  // avoid memory leak
	item.index = -1 // for safety
	*qData = old[0 : n-1]
	return item
}

// update modifies the Priority of an QItem in the queue.
func (qData *QItems) update(item *QItem, priority int) {

	item.Priority = priority
	heap.Fix(qData, item.index)
}

func (qData *QItems) delete(index int) error {
	if index < qData.Len() {
		// Remove the element at index i from a.
		heap.Remove(qData, index)
		return nil
	}
	return fmt.Errorf("Index out of range")
}
