package components

import (
	"codeberg.org/anaseto/gruid"
)

// Name component represents an entity's name
type Name struct {
	Name string
}

// Position component represents an entity's position in the game world
type Position struct {
	Point gruid.Point
}

// Renderable component represents how an entity is rendered
type Renderable struct {
	Glyph rune
	Color gruid.Color
}

// Health component represents an entity's health points
type Health struct {
	CurrentHP int
	MaxHP     int
}
