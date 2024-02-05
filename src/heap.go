
package main


// MinHeap represents a min-heap of tuples (address, connection)
type MinHeap [][]interface{}

// Implement the heap.Interface methods for MinHeap
func (h MinHeap) Len() int           { return len(h) }

func (h MinHeap) Less(i, j int) bool {
    // Ensure to assert the correct type when accessing connection count
    return h[i][1].(int) < h[j][1].(int)
}

func (h MinHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *MinHeap) Push(x interface{}) {
    // Correctly append the new element
    *h = append(*h, x.([]interface{}))
}

func (h *MinHeap) Pop() interface{} {
    old := *h
    n := len(old)
    x := old[n-1]
    *h = old[0 : n-1]
    return x
}

