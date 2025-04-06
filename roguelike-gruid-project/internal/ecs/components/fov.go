package components

import (
	"codeberg.org/anaseto/gruid" // Import base gruid package
	"codeberg.org/anaseto/gruid/rl"
)

// FOV holds data related to an entity's field of view.
type FOV struct {
	Range   int      // Vision radius
	Visible []uint64 // Bitset for tiles currently visible to this entity
	fov     *rl.FOV  // Reusable FOV calculator instance
}

// NewFOVComponent creates and initializes a new FOV component.
// It requires the map dimensions to correctly size the internal bitset and FOV calculator.
func NewFOVComponent(fovRange, mapWidth, mapHeight int) *FOV {
	bitsetSize := (mapWidth*mapHeight + 63) / 64
	mapRange := gruid.NewRange(0, 0, mapWidth, mapHeight) // Use gruid.NewRange

	return &FOV{
		Range:   fovRange,
		Visible: make([]uint64, bitsetSize),
		fov:     rl.NewFOV(mapRange), // Initialize the gruid FOV calculator
	}
}

// IsVisible checks if a point is currently visible according to this component's bitset.
// Assumes mapWidth is passed correctly or accessible.
func (f *FOV) IsVisible(p gruid.Point, mapWidth int) bool { // Use gruid.Point
	// Basic bounds check (optional, as FOV calc should handle it, but safer)
	if p.X < 0 || p.Y < 0 || p.X >= mapWidth { // Need mapHeight too ideally, but width is key for index
		return false
	}
	idx := p.Y*mapWidth + p.X
	if idx < 0 { // Should not happen with proper bounds check
		return false
	}
	sliceIdx := idx / 64
	bitIdx := uint(idx % 64)

	// Check if sliceIdx is within the bounds of the Visible slice
	if sliceIdx >= len(f.Visible) {
		return false // Out of bounds for the bitset
	}

	return (f.Visible[sliceIdx] & (1 << bitIdx)) != 0
}

// GetFOVCalculator returns the internal rl.FOV instance for calculations.
func (f *FOV) GetFOVCalculator() *rl.FOV {
	return f.fov
}

// ClearVisible resets the visible bitset.
func (f *FOV) ClearVisible() {
	for i := range f.Visible {
		f.Visible[i] = 0
	}
}

// SetVisible marks a point as visible in this component's bitset.
// Assumes mapWidth is passed correctly or accessible.
func (f *FOV) SetVisible(p gruid.Point, mapWidth int) { // Use gruid.Point
	// Basic bounds check
	if p.X < 0 || p.Y < 0 || p.X >= mapWidth { // Need mapHeight too ideally
		return
	}
	idx := p.Y*mapWidth + p.X
	if idx < 0 {
		return
	}
	sliceIdx := idx / 64
	bitIdx := uint(idx % 64)

	// Check if sliceIdx is within the bounds of the Visible slice
	if sliceIdx < len(f.Visible) {
		f.Visible[sliceIdx] |= (1 << bitIdx)
	}
}
