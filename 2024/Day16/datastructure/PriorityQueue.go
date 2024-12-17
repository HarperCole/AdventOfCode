package datastructure

import (
	"container/heap"
)

// PriorityQueueItem defines the structure of each item in the priority queue.
type PriorityQueueItem struct {
	Position  [2]int // (x, y) position
	Direction int    // Direction (N, S, E, W)
	Cost      int    // g(x) score for A* search
	Index     int    // Index in the heap for updating
}

// PriorityQueue implements a priority queue for PriorityQueueItem.
type PriorityQueue []*PriorityQueueItem

// Len returns the length of the priority queue.
func (pq PriorityQueue) Len() int { return len(pq) }

// Less compares two items to determine priority (min-heap based on GScore).
func (pq PriorityQueue) Less(i, j int) bool {
	return pq[i].Cost < pq[j].Cost
}

// Swap swaps two items in the priority queue.
func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].Index = i
	pq[j].Index = j
}

// Push adds an item to the priority queue.
func (pq *PriorityQueue) Push(x interface{}) {
	item := x.(*PriorityQueueItem)
	item.Index = len(*pq)
	*pq = append(*pq, item)
}

// Pop removes and returns the lowest-priority item (min GScore).
func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil // Avoid memory leak
	item.Index = -1
	*pq = old[0 : n-1]
	return item
}

// Update modifies the GScore of an item in the queue and adjusts its position.
func (pq *PriorityQueue) Update(item *PriorityQueueItem, cost int) {
	item.Cost = cost
	heap.Fix(pq, item.Index)
}

// NewPriorityQueue creates and initializes a new priority queue.
func NewPriorityQueue() *PriorityQueue {
	pq := &PriorityQueue{}
	heap.Init(pq)
	return pq
}
