package priorityqueue

import (
	"container/heap"
	"reflect"
	"strconv"
	"testing"
)

func assertEqual(t *testing.T, a interface{}, b interface{}) {
	if a == b {
		return
	}
	t.Errorf("Received %v (type %v), expected %v (type %v)", a, reflect.TypeOf(a), b, reflect.TypeOf(b))
}

func populateQueue(pq *PriorityQueue, itemsToPopulate int) {
	for i := 0; i < itemsToPopulate; i++ {
		var y = QItem{
			ParentID: "12345",
			ID:       strconv.Itoa(i),
			Value:    "test",
			Priority: i + 1,
		}
		pq.Push(y)
	}
	//heap.Init(&pq.data)
}

func Test_NewPriorityQueue(t *testing.T) {
	pq := NewPriorityQueue()
	assertEqual(t, pq.Len(), 0)
}

func Test_Len(t *testing.T) {
	expectedLen := 2
	pq := NewPriorityQueue()
	var y = QItem{
		ID:       "Test",
		Value:    "test",
		Priority: 1,
	}
	pq.Push(y)
	pq.Push(y)

	assertEqual(t, pq.Len(), expectedLen)

}

func Test_Push(t *testing.T) {
	expectedItems := 10
	pq := NewPriorityQueue()
	populateQueue(pq, expectedItems)

	if pq.Len() != expectedItems {
		t.Errorf("Underlying slice not populated")
	}

	expectedPriority := 10

	y, _ := pq.Pop()
	actualPriority := y.Priority
	if actualPriority != expectedPriority {
		t.Errorf("Actual priority: %d should be: %d", actualPriority, expectedPriority)
	}

}

func Test_Pop(t *testing.T) {
	expectedItems := 10
	pq := NewPriorityQueue()
	populateQueue(pq, expectedItems)
	for pq.Len() > 0 {
		_, err := pq.Pop()
		if err != nil {
			t.Errorf("Error popping item from queue")
		}
	}
}

func Test_ClearQueue(t *testing.T) {
	expectedItems := 10
	pq := NewPriorityQueue()
	populateQueue(pq, expectedItems)
	pq.Clear()
	if pq.Len() != 0 {
		t.Errorf("Queue is not empty after Clear()")
	}
}

func Test_PopPastEndOfQueue(t *testing.T) {
	expectedItems := 10
	pq := NewPriorityQueue()
	populateQueue(pq, expectedItems)
	for pq.Len() > 0 {
		_, err := pq.Pop()
		if err != nil {
			t.Errorf("Error popping item from queue")
		}
	}
	_, err := pq.Pop()
	if err == nil {
		t.Errorf("Pop past end of queue should return an error")
	}
}

func Test_PopWithHigherPriorityPush(t *testing.T) {
	expectedItems := 10
	pq := NewPriorityQueue()
	populateQueue(pq, expectedItems)
	// At this point we have 10 items in the queue
	// The highest priority item should be equal
	// to expectedItems

	// Inject a high priority item into the queue
	expectedPriority := 500
	var y = QItem{
		ID:       strconv.Itoa(1234),
		Value:    "test",
		Priority: expectedPriority,
	}
	pq.Push(y)

	// Now attempt to pop the highest priority item
	x, err := pq.Pop()
	if err != nil {
		t.Error(err)
	} else {
		if x.Priority != expectedPriority {
			t.Errorf("Highest priority item failed to pop out of the queue\nItem that popped: %-v", x)
		}
	}
}

func Test_PopWithLowerPriorityPush(t *testing.T) {
	expectedItems := 10
	pq := NewPriorityQueue()
	populateQueue(pq, expectedItems)
	// At this point we have 10 items in the queue
	// The highest priority item should be equal
	// to expectedItems

	// Inject a high priority item into the queue
	HighPriority := 500
	var y = QItem{
		ID:       strconv.Itoa(500),
		Value:    "test",
		Priority: HighPriority,
	}
	pq.Push(y)

	// Inject a medium priority item into the queue
	// Inject a high priority item into the queue
	MediumPriority := 100
	var yy = QItem{
		ID:       strconv.Itoa(100),
		Value:    "test",
		Priority: MediumPriority,
	}
	pq.Push(yy)
	heap.Init(&pq.data)
	// Now attempt to pop the highest priority item
	x, err := pq.Pop()
	if err != nil {
		t.Error(err)
	} else {
		if x.Priority != HighPriority {
			t.Errorf("Highest priority item failed to pop out of the queue: %v", x)
		}
	}
}

func Test_UpdatePriorityByParentId(t *testing.T) {
	expectedItems := 10
	pq := NewPriorityQueue()
	populateQueue(pq, expectedItems)
	// At this point we have 10 items in the queue
	// The highest priority item should be equal
	// to expectedItems, the parentID is "12345"
	parentID := "12345"
	NewPriority := 500
	pq.UpdatePriorityByParentId(parentID, NewPriority)
	for pq.Len() > 0 {
		x, err := pq.Pop()
		if err != nil {
			t.Errorf("Error retrieving item: %e", err)
		}
		if x.Priority != NewPriority {
			t.Errorf("Priority not set for item: %v", x)
		}
	}

}

func Test_DeleteItemById(t *testing.T) {
	expectedItems := 10
	pq := NewPriorityQueue()
	populateQueue(pq, expectedItems)
	// At this point we have 10 items in the queue
	// The highest priority item should be equal
	// to expectedItems, the parentID is "12345"
	// with ID values of "0" to "9"
	queueLength := pq.Len()
	IdToDelete := "5"
	err := pq.DeleteItemById(IdToDelete)
	if err != nil {
		t.Errorf("Error deleting item by id: %e", err)
	}

	if queueLength <= pq.Len() {
		t.Errorf("Queue length did not change")
	}
	for pq.Len() > 0 {
		x, err := pq.Pop()
		if err != nil {
			t.Errorf("Error accessing queue: %e", err)
		}
		if x.ID == IdToDelete {
			t.Errorf("Item with ID: %s was not deleted from the queue: %v", IdToDelete, x)
		}
	}
}

func Test_DeleteItemsByParentId(t *testing.T) {
	expectedItems := 10
	pq := NewPriorityQueue()
	populateQueue(pq, expectedItems)
	// At this point we have 10 items in the queue
	// The highest priority item should be equal
	// to expectedItems

	// Add two items with a unique ParentID to the queue
	parentID := "parentid"
	var y = QItem{
		ParentID: parentID,
		ID:       "testitemID",
		Value:    "test",
		Priority: 7,
	}
	var yy = QItem{
		ParentID: parentID,
		ID:       "testitemID",
		Value:    "test2",
		Priority: 7,
	}
	baseQueueLength := pq.Len()
	pq.Push(y)
	pq.Push(yy)

	queueLength := pq.Len()
	itemsToDelete := queueLength - baseQueueLength

	_, err := pq.DeleteItemsByParentId(parentID)
	if err != nil {
		t.Errorf("Error deleting item by parentID: %e", err)
	}

	if queueLength <= pq.Len() {
		t.Errorf("Queue length did not change")
	}

	if queueLength-pq.Len() != itemsToDelete {
		t.Errorf("failed to delete both items")
	}

	for pq.Len() > 0 {
		x, err := pq.Pop()
		if err != nil {
			t.Errorf("Error accessing queue: %e", err)
		}
		if x.ID == parentID {
			t.Errorf("Item with parent ID: %s still present in the queue: %v", parentID, x)
		}
	}
}
