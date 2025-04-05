package game

import (
	"container/heap"
	"fmt"
	"time"

	"codeberg.org/anaseto/gruid"
	"github.com/sirupsen/logrus"
)

// GameAction is an interface for actions that can be performed in the game.
type GameAction interface {
	Execute(g *Game) (cost uint, err error) // Modified return signature
}

type MoveAction struct {
	Direction gruid.Point
}

// Execute performs the move action, returning the time cost and any error.
func (a MoveAction) Execute(g *Game) (cost uint, err error) {
	// TODO: This currently assumes the action is for the player.
	// Need to generalize for any entity ID.
	// For now, we'll pass g.PlayerID, but this needs refactoring later.
	again, err := g.EntityBump(g.PlayerID, a.Direction) // Assuming EntityBump exists or will be created
	if err != nil {
		return 0, err // No cost if error occurred
	}

	if again {
		// Bumped into something, action didn't fully succeed in moving
		return 0, nil // No time cost for a bump? Or maybe a smaller cost? Let's use 0 for now.
	}

	// Successful move
	return 100, nil // Standard move cost
}

// --- TurnEntry ---

// TurnEntry represents an item in the turn queue.
// It holds the time (priority) and the entity's ID.
// This corresponds to the (u64, Entity) tuple in Rust.
type TurnEntry struct {
	Time     uint64   // The time/turn number for this entity's action
	EntityID EntityID // The ID of the entity
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

// --- TurnQueue Struct ---

// TurnQueue manages entity turns based on time using a min-heap.
// It corresponds to the Rust TurnQueue struct.
type TurnQueue struct {
	// CurrentTime corresponds to current_time: u64
	CurrentTime uint64

	// turnQueue holds the actual heap data structure.
	// We use a pointer to the slice so methods like Push/Pop can modify it.
	// This corresponds to turn_queue: BinaryHeap<Reverse<(u64, Entity)>>
	queue *turnHeap // Note: unexported field to encourage using methods

	// OperationsSinceCleanup corresponds to operations_since_cleanup: u32
	OperationsSinceCleanup uint32

	// TotalCleanups corresponds to total_cleanups: u64
	TotalCleanups uint64

	// TotalEntitiesRemoved corresponds to total_entities_removed: u64
	TotalEntitiesRemoved uint64
}

// --- Constructor and Methods ---

// NewTurnQueue creates and initializes a new TurnQueue.
// This creates a Go equivalent of Rust's BinaryHeap<Reverse<(u64, Entity)>>,
// which is a min-heap where entries with the smallest time are processed first.
func NewTurnQueue() *TurnQueue {
	// Create an empty turnHeap
	h := &turnHeap{}

	// Initialize the heap structure
	// This establishes the heap invariant: for every element at index i,
	// the element at index 2*i+1 and 2*i+2 (if they exist) are children of the element at i.
	heap.Init(h)

	return &TurnQueue{
		CurrentTime:            0,
		queue:                  h,
		OperationsSinceCleanup: 0,
		TotalCleanups:          0,
		TotalEntitiesRemoved:   0,
	}
}

// Add adds an entity with its scheduled time to the queue.
func (tq *TurnQueue) Add(entityID EntityID, time uint64) {
	// Create a new entry with the specified time and entity ID
	entry := TurnEntry{Time: time, EntityID: entityID}

	// Push the entry onto the heap
	// The heap package will maintain the min-heap property
	heap.Push(tq.queue, entry)

	// For debugging purposes
	logrus.Debugf("Added entity %d to turn queue with time %d", entityID, time)
}

// Next removes and returns the next entity (the one with the smallest time)
// from the queue. It returns the entry and true if the queue is not empty,
// otherwise, it returns a zero TurnEntry and false.
func (tq *TurnQueue) Next() (TurnEntry, bool) {
	if tq.queue.Len() == 0 {
		return TurnEntry{}, false
	}

	// heap.Pop maintains the min-heap property by:
	// 1. Moving the last element to the root
	// 2. Shifting it down until heap property is restored
	// 3. Returning the original root
	entry := heap.Pop(tq.queue).(TurnEntry)

	logrus.Debugf("Popped entity %d from turn queue (time: %d)", entry.EntityID, entry.Time)
	return entry, true
}

// Peek returns the next entity (the one with the smallest time) without
// removing it from the queue. It returns the entry and true if the queue
// is not empty, otherwise, it returns a zero TurnEntry and false.
func (tq *TurnQueue) Peek() (TurnEntry, bool) {
	if tq.queue.Len() == 0 {
		return TurnEntry{}, false
	}

	// In a min-heap, the minimum element is always at index 0
	entry := (*tq.queue)[0]
	logrus.Debugf("Peeked entity %d from turn queue (time: %d)", entry.EntityID, entry.Time)
	return entry, true
}

// Len returns the number of entities currently in the queue.
func (tq *TurnQueue) Len() int {
	return tq.queue.Len()
}

// IsEmpty checks if the turn queue is empty.
func (tq *TurnQueue) IsEmpty() bool {
	return tq.queue.Len() == 0
}

// PrintQueue prints the current state of the turn queue.
// This is primarily for debugging purposes.
func (tq *TurnQueue) PrintQueue() {
	if tq.IsEmpty() {
		logrus.Debug("---- Turn Queue: EMPTY ----")
		return
	}

	logrus.Debug("---- Turn Queue Contents ----")
	logrus.Debugf("Current Game Time: %d\n", tq.CurrentTime)
	logrus.Debugf("Queue Size: %d\n", tq.Len())

	// Display the raw queue in heap order
	logrus.Debug("Queue (in heap order):")
	for i, entry := range *tq.queue {
		delta := int64(entry.Time) - int64(tq.CurrentTime)
		logrus.Debugf("[%d] EntityID: %d, Time: %d (Δ%d from current)\n",
			i, entry.EntityID, entry.Time, delta)
	}

	// Also display the queue in sorted order (to show actual processing order)
	logrus.Debug("\nProcessing order (sorted by time):")
	// Make a copy of the queue to sort
	sorted := make([]TurnEntry, len(*tq.queue))
	copy(sorted, *tq.queue)

	// Sort by time, then by EntityID for stable ordering
	tq.sortEntriesByTime(sorted)

	for i, entry := range sorted {
		delta := int64(entry.Time) - int64(tq.CurrentTime)
		logrus.Debugf("%d. EntityID: %d, Time: %d (Δ%d from current)\n",
			i+1, entry.EntityID, entry.Time, delta)
	}

	logrus.Debug("----------------------------")
}

// sortEntriesByTime sorts a slice of TurnEntry by time, then by EntityID
func (tq *TurnQueue) sortEntriesByTime(entries []TurnEntry) {
	// Use a simple insertion sort for small lists (typically the case for turn queues)
	for i := 1; i < len(entries); i++ {
		for j := i; j > 0; j-- {
			// Primary sort by time
			if entries[j].Time < entries[j-1].Time {
				entries[j], entries[j-1] = entries[j-1], entries[j]
			} else if entries[j].Time == entries[j-1].Time &&
				entries[j].EntityID < entries[j-1].EntityID {
				// Secondary sort by EntityID for stable ordering
				entries[j], entries[j-1] = entries[j-1], entries[j]
			} else {
				// Stop if we're in the right position
				break
			}
		}
	}
}

// --- Cleanup Logic ---

// --- CleanupMetrics Struct ---
// Corresponds to the Rust CleanupMetrics struct.
type CleanupMetrics struct {
	EntitiesRemoved int // Use int for counts, common in Go
	QueueSizeBefore int
	QueueSizeAfter  int
	ProcessingTime  time.Duration // Go's standard duration type
}

// Implement fmt.Stringer for nice printing (like Rust's Debug derive)
func (m CleanupMetrics) String() string {
	return fmt.Sprintf(
		"CleanupMetrics{Removed: %d, Before: %d, After: %d, Time: %v}",
		m.EntitiesRemoved,
		m.QueueSizeBefore,
		m.QueueSizeAfter,
		m.ProcessingTime,
	)
}

// --- Helper Functions ---

// isValIDTurnActor checks if an entity is valid to remain in the turn queue.
// Corresponds to the Rust method.
func (tq *TurnQueue) isValIDTurnActor(world *ECS, entityID EntityID) bool {
	// First, check if entity exists at all
	if !world.EntityExists(entityID) {
		return false
	}

	// Check if it has required TurnActor component
	// if !world.HasComponent(entityID, GetReflectType(TurnActor{})) {
	// 	return false
	// }

	// Check for "dead" markers (game-specific logic)
	// if world.HasComponent(entityID, GetReflectType(Dead{})) {
	// 	return false
	// }

	// Check for health <= 0 (if Health component exists)
	/*
	   // Example health check:
	   if comp, found := world.GetComponent(entityID, HealthComponentType); found {
	       if health, ok := comp.(*Health); ok { // Type assertion
	           if health.Current <= 0 {
	               return false
	           }
	       } else {
	           // Log error if component type is wrong?
	           log.Printf("Warning: Found component for HealthComponentType but it wasn't a *Health for entity %d", entityID)
	       }
	   }
	*/

	return true
}

// --- Cleanup Functions ---

func (tq *TurnQueue) getCleanupThreshold(world *ECS) uint32 {
	// Base threshold
	base_threshold := 100

	entityCount := len(world.Entities)
	queueSize := tq.Len()

	// More frequent cleanup with larger entity counts or queue sizes
	if entityCount > 1000 || queueSize > 500 {
		return uint32(base_threshold / 2)
	} else if entityCount < 100 && queueSize < 50 {
		return uint32(base_threshold * 2)
	}

	return uint32(base_threshold)
}

// CleanupDeadEntities removes invalid or dead entities from the queue.
// Corresponds to the Rust method.
func (tq *TurnQueue) CleanupDeadEntities(world *ECS) CleanupMetrics {
	// Only run periodically to amortize cost
	threshold := tq.getCleanupThreshold(world)
	if tq.OperationsSinceCleanup < threshold {
		tq.OperationsSinceCleanup++
		// Return zero-value metrics if no cleanup performed
		return CleanupMetrics{}
	}

	logrus.Debug("TurnQueue: Cleaning up dead entities...") // Use Go's log package

	queueSizeBefore := tq.Len()
	startTime := time.Now() // Use Go's time package

	// Create a temporary slice to hold valid entries.
	// Pre-allocate capacity close to the original size for efficiency.
	newQueueSlice := make(turnHeap, 0, queueSizeBefore)
	removedCount := 0

	// Process each entity by popping from the *original* heap
	// Note: We modify the original heap in place by popping.
	originalQueue := tq.queue // Keep a reference if needed, but Pop modifies tq.queue directly
	for originalQueue.Len() > 0 {
		// Pop directly from the underlying heap implementation
		entry := heap.Pop(originalQueue).(TurnEntry)

		// Check validity using the helper method
		entityValid := tq.isValIDTurnActor(world, entry.EntityID)

		if entityValid {
			// Keep valid entities by adding them to the new slice
			// W	e don't use heap.Push here yet, just build the slice.
			newQueueSlice = append(newQueueSlice, entry)
		} else {
			// Count removed entities
			removedCount++

			name, ok := world.GetName(entry.EntityID)
			if ok {
				// Log removed entity (using helper for name)
				logrus.Debugf("TurnQueue: Removed dead entity from turn queue: %s\n",
					name)
			} else {
				logrus.Debugf("TurnQueue: Removed dead entity from turn queue: %d\n",
					entry.EntityID)
			}
		}
	}

	// Replace the queue's underlying slice with the cleaned slice
	tq.queue = &newQueueSlice
	// IMPORTANT: Re-establish the heap invariant on the new slice
	heap.Init(tq.queue)

	// Update metrics
	tq.OperationsSinceCleanup = 0 // Reset counter
	tq.TotalCleanups++
	tq.TotalEntitiesRemoved += uint64(removedCount) // Convert int to uint64

	// Return cleanup metrics
	metrics := CleanupMetrics{
		EntitiesRemoved: removedCount,
		QueueSizeBefore: queueSizeBefore,
		QueueSizeAfter:  tq.Len(), // Get current length after cleanup
		ProcessingTime:  time.Since(startTime),
	}

	logrus.Debugf("TurnQueue: Cleanup finished. %s\n", metrics)
	return metrics
}
