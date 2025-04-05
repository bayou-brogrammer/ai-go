package turn

import "github.com/lecoqjacob/ai-go/roguelike-gruid-project/internal/ecs"

// --- TurnEntry ---

// TurnEntry represents an item in the turn queue.
// It holds the time (priority) and the entity's ID.
// This corresponds to the (u64, Entity) tuple in Rust.
type TurnEntry struct {
	Time     uint64       // The time/turn number for this entity's action
	EntityID ecs.EntityID // The ID of the entity
	// You could add an index field here if needed for heap.Update
	// index int
}

// --- Heap Implementation ---
// We need a type that implements heap.Interface based on TurnEntry.
// This mirrors Rust's BinaryHeap<Reverse<(u64, Entity)>> by creating a min-heap
// where the smallest time values are at the top of the heap.
type turnHeap []TurnEntry

// Len returns the number of elements in the heap.
func (h turnHeap) Len() int { return len(h) }

// Less reports whether the element with index i should sort before the element with index j.
// To create a min-heap (equivalent to Rust's BinaryHeap<Reverse<...>>), we return true
// when h[i].Time < h[j].Time. This ensures entries with the smallest time are at the top.
func (h turnHeap) Less(i, j int) bool {
	// Primary sort by time
	if h[i].Time != h[j].Time {
		return h[i].Time < h[j].Time
	}
	// Secondary sort by EntityID (for stable ordering when times are equal)
	return h[i].EntityID < h[j].EntityID
}

// Swap swaps the elements with indexes i and j.
func (h turnHeap) Swap(i, j int) { h[i], h[j] = h[j], h[i] }

// Push adds an element to the heap.
// NOTE: heap.Push calls this; you don't call it directly.
// The argument x is asserted to the concrete type TurnEntry.
func (h *turnHeap) Push(x any) {
	*h = append(*h, x.(TurnEntry))
}

// Pop removes and returns the minimum element (root) from the heap.
// NOTE: heap.Pop calls this; you don't call it directly.
func (h *turnHeap) Pop() any {
	old := *h
	n := len(old)
	item := old[n-1]  // Get the last element
	*h = old[0 : n-1] // Truncate the slice
	return item       // Return the popped item (which was the root)
}
