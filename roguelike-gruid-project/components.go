package main

import "codeberg.org/anaseto/gruid"

// Position component holds the entity's location on the map.
type Position struct {
	Point gruid.Point
}

// Renderable component holds information needed to draw the entity.
type Renderable struct {
	Glyph rune        // The character symbol
	Color gruid.Color // The foreground color
	// Add background color if needed later: BG gruid.Color
}

// BlocksMovement is a tag component indicating that an entity blocks movement.
// Entities with this component (like walls or other creatures) prevent others
// from moving into their tile.
type BlocksMovement struct{}
