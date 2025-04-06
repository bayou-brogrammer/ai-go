package turn

import "github.com/lecoqjacob/ai-go/roguelike-gruid-project/internal/ecs"

// --- TurnEntry ---

// TurnEntry represents an item in the turn queue.
// It holds the time (priority) and the entity's ID.
// This corresponds to the (u64, Entity) tuple in Rust.
type TurnEntry struct {
	Time     uint64       // The time/turn number for this entity's action
	EntityID ecs.EntityID // The ID of the entity
}

// --- Heap Implementation ---
// This mirrors Rust's BinaryHeap<Reverse<(u64, Entity)>> by creating a min-heap
// where the smallest time values are at the top of the heap.
type turnHeap []TurnEntry

// Len returns the number of elements in the heap.
func (h turnHeap) Len() int { return len(h) }

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

func (h *turnHeap) Push(x any) {
	*h = append(*h, x.(TurnEntry))
}

func (h *turnHeap) Pop() any {
	old := *h
	n := len(old)
	item := old[n-1]
	*h = old[0 : n-1]
	return item
}
